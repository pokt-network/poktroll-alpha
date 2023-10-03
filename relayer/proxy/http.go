package proxy

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"poktroll/relayer/client"
	"poktroll/utils"
	serviceTypes "poktroll/x/service/types"
	servicerTypes "poktroll/x/servicer/types"
	sessionTypes "poktroll/x/session/types"
)

// TODO_COMMENT: For this and other important foundational structs, we really need to comment the type and
// the fields inside of it.
type httpProxy struct {
	// TODO: Replace with sessionHeader
	serviceId *serviceTypes.ServiceId
	// TODO: replace with servicerAddress?
	serviceForwardingAddr string
	sessionQueryClient    sessionTypes.QueryClient
	client                client.ServicerClient
	relayNotifier         chan *RelayWithSession
	signResponseFn        responseSigner
	servicerAddress       string
}

func NewHttpProxy(
	serviceId *serviceTypes.ServiceId,
	serviceForwardingAddr string,
	sessionQueryClient sessionTypes.QueryClient,
	client client.ServicerClient,
	relayNotifier chan *RelayWithSession,
	signResponse responseSigner,
	servicerAddress string,
) *httpProxy {
	return &httpProxy{
		serviceId:             serviceId,
		serviceForwardingAddr: serviceForwardingAddr,
		sessionQueryClient:    sessionQueryClient,
		client:                client,
		relayNotifier:         relayNotifier,
		signResponseFn:        signResponse,
		servicerAddress:       servicerAddress,
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
	ctx := req.Context()
	relayRequest, err := newHTTPRelayRequest(req)
	if err != nil {
		utils.ReplyWithHTTPError(500, err, httpResponseWriter)
		return
	}

	query := &sessionTypes.QueryGetSessionRequest{
		AppAddress:  relayRequest.ApplicationAddress,
		ServiceId:   httpProxy.serviceId.Id,
		BlockHeight: httpProxy.client.LatestBlock(ctx).Height(),
	}

	sessionResult, err := httpProxy.sessionQueryClient.GetSession(ctx, query)
	if err != nil {
		utils.ReplyWithHTTPError(500, err, httpResponseWriter)
		return
	}

	if err := validateSessionRequest(&sessionResult.Session, relayRequest); err != nil {
		utils.ReplyWithHTTPError(400, err, httpResponseWriter)
		return
	}

	url, err := parseURLWithScheme(httpProxy.serviceForwardingAddr)
	if err != nil {
		utils.ReplyWithHTTPError(400, err, httpResponseWriter)
	}

	headers := make(http.Header)
	for k, v := range relayRequest.Headers {
		headers.Add(k, v)
	}

	serviceRequest := &http.Request{
		Method: relayRequest.Method,
		Header: headers,
		URL:    url,
		Host:   url.Host,
		Body:   io.NopCloser(bytes.NewBuffer(relayRequest.Payload)),
	}
	relayResponse, err := httpProxy.executeRelay(serviceRequest, sessionResult.Session.SessionId)
	if err != nil {
		utils.ReplyWithHTTPError(500, err, httpResponseWriter)
		return
	}

	if err := sendRelayResponse(relayResponse, httpResponseWriter); err != nil {
		log.Printf("failed sending relay response: %v", err)
		return
	}

	relayWithSession := &RelayWithSession{
		Relay: &servicerTypes.Relay{
			Req: relayRequest,
			Res: relayResponse,
		},
		Session: &sessionResult.Session,
	}

	httpProxy.relayNotifier <- relayWithSession
}

func (httpProxy *httpProxy) executeRelay(req *http.Request, sessionId string) (*servicerTypes.RelayResponse, error) {
	// Change the request host to the service address
	serviceResponse, err := proxyHTTPServiceRequest(req)
	if err != nil {
		return nil, err
	}

	relayResponse, err := newRelayResponse(serviceResponse, httpProxy.servicerAddress, sessionId)
	if err != nil {
		return nil, err
	}

	signature, err := httpProxy.signResponseFn(relayResponse)
	if err != nil {
		return nil, err
	}

	relayResponse.ServicerSignature = signature
	return relayResponse, nil
}

func newHTTPRelayRequest(req *http.Request) (*servicerTypes.RelayRequest, error) {
	requestBz, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	var relayRequest servicerTypes.RelayRequest
	if err := relayRequest.Unmarshal(requestBz); err != nil {
		return nil, err
	}

	return &relayRequest, nil
}

func proxyHTTPServiceRequest(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

func newRelayResponse(
	serviceResponse *http.Response,
	servicerAddress string,
	sessionId string,
) (_ *servicerTypes.RelayResponse, err error) {
	relayResponse := &servicerTypes.RelayResponse{
		Headers:         make(map[string]string),
		StatusCode:      int32(serviceResponse.StatusCode),
		ServicerAddress: servicerAddress,
		SessionId:       sessionId,
	}

	if serviceResponse.Body != nil {
		// Read the response from the service
		relayResponse.Payload, err = io.ReadAll(serviceResponse.Body)
		if err != nil {
			return nil, err
		}
	}

	for key, value := range serviceResponse.Header {
		relayResponse.Headers[key] = strings.Join(value, ", ")
	}
	return relayResponse, nil
}

func sendRelayResponse(relayResponse *servicerTypes.RelayResponse, wr http.ResponseWriter) error {
	relayResponseBz, err := relayResponse.Marshal()
	if err != nil {
		return err
	}

	if _, err := wr.Write(relayResponseBz); err != nil {
		return err
	}
	return nil
}
