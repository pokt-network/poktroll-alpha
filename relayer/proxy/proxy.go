package proxy

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"net/http"

	"poktroll/utils"
	"poktroll/x/servicer/types"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

type Proxy struct {
	localAddr        string
	serviceAddr      string
	keyring          keyring.Keyring
	keyName          string
	logger           *log.Logger
	output           chan *types.Relay
	outputObservable utils.Observable[*types.Relay]
}

func NewProxy(logger *log.Logger, keyring keyring.Keyring, keyName string) *Proxy {
	proxy := &Proxy{
		output:  make(chan *types.Relay),
		logger:  logger,
		keyring: keyring,
		keyName: keyName,
	}

	proxy.outputObservable, _ = utils.NewControlledObservable[*types.Relay](proxy.output)

	proxy.localAddr = "localhost:8545"
	proxy.serviceAddr = "localhost:8546"

	go proxy.listen()

	return proxy
}

func (proxy *Proxy) Relays() utils.Observable[*types.Relay] {
	return proxy.outputObservable
}

func (proxy *Proxy) listen() {
	if err := http.ListenAndServe(proxy.localAddr, proxy); err != nil {
		proxy.logger.Fatal(err)
	}
}

// ServeHTTP implements the http.Handler interface; called by http.ListenAndServe().
// It re-uses the incoming request, updating the host and URL to match the service,
// the body to a new io.ReadCloser containing the relay request payload, and then
// sending it to the service.
func (proxy *Proxy) ServeHTTP(httpResponseWriter http.ResponseWriter, req *http.Request) {
	relayRequest, err := newRelayRequest(req)
	if err != nil {
		if err := proxy.replyWithError(500, err, httpResponseWriter); err != nil {
			// TECHDEBT: log error
		}
		return
	}

	// Change the request host to the service address
	req.Host = proxy.serviceAddr
	req.URL.Host = proxy.serviceAddr
	req.Body = io.NopCloser(bytes.NewBuffer(relayRequest.Payload))

	relayResponse, err := proxy.executeRelay(req)
	if err != nil {
		if err := proxy.replyWithError(500, err, httpResponseWriter); err != nil {
			// TECHDEBT: log error
		}
		return
	}

	if err := sendRelayResponse(relayResponse, httpResponseWriter); err != nil {
		// TODO: log error
		return
	}

	relay := &types.Relay{
		Req: relayRequest,
		Res: relayResponse,
	}

	proxy.output <- relay
}

func (proxy *Proxy) signResponse(relayResponse *types.RelayResponse) error {
	relayResBz, err := relayResponse.Marshal()
	if err != nil {
		return err
	}

	relayResponse.Signature, _, err = proxy.keyring.Sign(proxy.keyName, relayResBz)
	return nil
}

func (proxy *Proxy) replyWithError(statusCode int, err error, wr http.ResponseWriter) error {
	wr.WriteHeader(statusCode)
	if _, err := wr.Write([]byte(err.Error())); err != nil {
		return err
	}
	return nil
}

func (proxy *Proxy) executeRelay(req *http.Request) (*types.RelayResponse, error) {
	serviceResponse, err := proxyServiceRequest(req)
	//http.ReadResponse(bufio.NewReader(remoteConnection), req)
	if err != nil {
		return nil, err
	}

	relayResponse, err := newRelayResponse(serviceResponse)
	if err != nil {
		return nil, err
	}

	if err := proxy.signResponse(relayResponse); err != nil {
		return nil, err
	}
	return relayResponse, nil
}

func newRelayRequest(req *http.Request) (*types.RelayRequest, error) {
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

func proxyServiceRequest(req *http.Request) (*http.Response, error) {
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
