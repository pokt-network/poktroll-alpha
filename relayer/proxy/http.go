package proxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	serviceTypes "poktroll/x/service/types"
	"poktroll/x/servicer/types"
	sessionTypes "poktroll/x/session/types"
)

type httpProxy struct {
	serviceId             *serviceTypes.ServiceId
	serviceForwardingAddr string
	sessionQueryClient    sessionTypes.QueryClient
	client                types.ServicerClient
	relayNotifier         chan *RelayWithSession
	signResponse          responseSigner
}

func NewHttpProxy(
	serviceId *serviceTypes.ServiceId,
	serviceForwardingAddr string,
	sessionQueryClient sessionTypes.QueryClient,
	client types.ServicerClient,
	relayNotifier chan *RelayWithSession,
	signResponse responseSigner,
) *httpProxy {
	return &httpProxy{
		serviceId:             serviceId,
		serviceForwardingAddr: serviceForwardingAddr,
		sessionQueryClient:    sessionQueryClient,
		client:                client,
		relayNotifier:         relayNotifier,
		signResponse:          signResponse,
	}
}

func (httpProxy *httpProxy) Start(advertisedEndpointUrl string) error {
	return http.ListenAndServe(mustGetHostAddress(advertisedEndpointUrl), httpProxy)
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
		ServiceId:   httpProxy.serviceId.Id,
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

	url, err := parseURLWithScheme(httpProxy.serviceForwardingAddr)
	if err != nil {
		replyWithHTTPError(400, err, httpResponseWriter)
	}

	serviceRequest := &http.Request{
		Method: req.Method,
		Header: req.Header,
		URL:    url,
		Host:   url.Host,
		Body:   io.NopCloser(bytes.NewBuffer(relayRequest.Payload)),
	}
	relayResponse, err := httpProxy.executeRelay(serviceRequest, relayRequest.Payload)
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
		// HACK: the application address should be populated by the requesting client
		relayRequest.ApplicationAddress = "pokt1mrqt5f7qh8uxs27cjm9t7v9e74a9vvdnq5jva4"
	}
	return relayRequest, nil
}

func proxyHTTPServiceRequest(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
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
	clientError := err
	if statusCode == 500 {
		clientError = fmt.Errorf("internal server error")
		log.Printf("internal server error: %v", err)
	}

	if _, replyError := wr.Write([]byte(clientError.Error())); replyError != nil {
		log.Printf("failed sending error response: %v", replyError)
	}
}