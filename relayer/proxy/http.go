package proxy

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log"
	"net"
	"net/http"

	"poktroll/x/servicer/types"
	sessionTypes "poktroll/x/session/types"
)

type httpProxy struct {
	serviceAddr        string
	sessionQueryClient sessionTypes.QueryClient
	client             types.ServicerClient
	relayNotifier      chan *RelayWithSession
	signResponse       responseSigner
	serviceId          string
}

func NewHttpProxy(
	serviceAddr string,
	sessionQueryClient sessionTypes.QueryClient,
	client types.ServicerClient,
	relayNotifier chan *RelayWithSession,
	signResponse responseSigner,
	serviceId string,
) *httpProxy {
	return &httpProxy{
		serviceAddr:        serviceAddr,
		sessionQueryClient: sessionQueryClient,
		client:             client,
		relayNotifier:      relayNotifier,
		signResponse:       signResponse,
		serviceId:          serviceId,
	}
}

// ServeHTTP implements the http.Handler interface; called by http.ListenAndServe().
// It re-uses the incoming request, updating the host and URL to match the service,
// the body to a new io.ReadCloser containing the relay request payload, and then
// sending it to the service.
func (httpProxy *httpProxy) ServeHTTP(httpResponseWriter http.ResponseWriter, req *http.Request) {
	relayRequest, err := newHTTPRelayRequest(req)
	if err != nil {
		log.Printf("failed creating relay request: %v", err)
		replyWithHTTPError(500, err, httpResponseWriter)
		return
	}

	query := &sessionTypes.QueryGetSessionRequest{
		AppAddress:  relayRequest.ApplicationAddress,
		ServiceId:   httpProxy.serviceId,
		BlockHeight: httpProxy.client.LatestBlock().Height(),
	}

	// INVESTIGATE: get the context instead of creating a new one?
	sessionResult, err := httpProxy.sessionQueryClient.GetSession(context.TODO(), query)
	if err != nil {
		log.Printf("failed getting session: %v", err)
		replyWithHTTPError(500, err, httpResponseWriter)
		return
	}

	if err := validateSessionRequest(&sessionResult.Session, relayRequest); err != nil {
		replyWithHTTPError(400, err, httpResponseWriter)
		return
	}

	relayResponse, err := httpProxy.executeRelay(req, relayRequest.Payload)
	if err != nil {
		log.Printf("failed executing relay: %v", err)
		replyWithHTTPError(500, err, httpResponseWriter)
		return
	}

	if err := sendRelayResponse(relayResponse, httpResponseWriter); err != nil {
		log.Printf("failed sending relay response: %v", err)
		return
	}

	relayWithSession := &RelayWithSession{
		Relay: &types.Relay{
			Req: relayRequest,
			Res: relayResponse,
		},
		Session: &sessionResult.Session,
	}

	httpProxy.relayNotifier <- relayWithSession
}

func (httpProxy *httpProxy) executeRelay(req *http.Request, requestPayload []byte) (*types.RelayResponse, error) {
	// Change the request host to the service address
	// DISCUSS: create a new request instead of mutating the existing one?
	req.Host = httpProxy.serviceAddr
	req.URL.Host = httpProxy.serviceAddr
	req.Body = io.NopCloser(bytes.NewBuffer(requestPayload))

	serviceResponse, err := proxyHTTPServiceRequest(req)
	if err != nil {
		return nil, err
	}

	relayResponse, err := newRelayResponse(serviceResponse)
	if err != nil {
		return nil, err
	}

	if err := httpProxy.signResponse(relayResponse); err != nil {
		return nil, err
	}
	return relayResponse, nil
}

func newHTTPRelayRequest(req *http.Request) (*types.RelayRequest, error) {
	requestHeaders := make(map[string]string)
	for k, v := range req.Header {
		requestHeaders[k] = v[0]
	}

	relayRequest := &types.RelayRequest{
		Method:  req.Method,
		Url:     req.URL.String(),
		Headers: requestHeaders,
	}

	if req.Body != nil {
		// Read the request body
		requestBody, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		relayRequest.Payload = requestBody
	}
	return relayRequest, nil
}

func proxyHTTPServiceRequest(req *http.Request) (*http.Response, error) {
	// Connect to the service
	remoteConnection, err := net.Dial("tcp", req.Host)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = remoteConnection.Close()
	}()

	// Send the request to the service
	err = req.Write(remoteConnection)
	if err != nil {
		return nil, err
	}

	// Read the response from the service
	return http.ReadResponse(bufio.NewReader(remoteConnection), req)
}

func newRelayResponse(serviceResponse *http.Response) (_ *types.RelayResponse, err error) {
	relayResponse := &types.RelayResponse{
		Headers:    make(map[string]string),
		StatusCode: int32(serviceResponse.StatusCode),
	}

	if serviceResponse.Body != nil {
		// Read the response from the service
		relayResponse.Payload, err = io.ReadAll(serviceResponse.Body)
		if err != nil {
			return nil, err
		}
	}

	for key, value := range serviceResponse.Header {
		// TECHDEBT: this drops all but the first value for headers with
		// multiple values
		relayResponse.Headers[key] = value[0]
	}
	return relayResponse, nil
}

func sendRelayResponse(relayResponse *types.RelayResponse, wr http.ResponseWriter) error {
	// Set HTTP statuscode to match the service response's
	wr.WriteHeader(int(relayResponse.StatusCode))

	// Set relay response headers to match the service response's
	for k, v := range relayResponse.Headers {
		wr.Header().Add(k, v)
	}

	// Send the response to the client
	if _, err := wr.Write(relayResponse.Payload); err != nil {
		return err
	}
	return nil
}

// TODO: send appropriate error instead of the original error
func replyWithHTTPError(statusCode int, err error, wr http.ResponseWriter) {
	wr.WriteHeader(statusCode)
	if _, replyError := wr.Write([]byte(err.Error())); replyError != nil {
		log.Printf("failed sending error response: %v", replyError)
	}
}
