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
)

type Proxy struct {
	localAddr        string
	serviceAddr      string
	logger           *log.Logger
	output           chan *types.Relay
	outputObservable utils.Observable[*types.Relay]
}

func NewProxy(logger *log.Logger) *Proxy {
	proxy := &Proxy{output: make(chan *types.Relay), logger: logger}
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

func (proxy *Proxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
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
			proxy.replyWithError(500, err, wr)
			return
		}
		relayRequest.Payload = requestBody
	}

	// Change the request host to the service address
	req.Host = proxy.serviceAddr
	req.URL.Host = proxy.serviceAddr
	req.Body = io.NopCloser(bytes.NewBuffer(relayRequest.Payload))

	// Connect to the service
	remoteConnection, err := net.Dial("tcp", proxy.serviceAddr)
	if err != nil {
		proxy.replyWithError(500, err, wr)
		return
	}
	defer remoteConnection.Close()

	// Send the request to the service
	err = req.Write(remoteConnection)
	if err != nil {
		proxy.replyWithError(500, err, wr)
		return
	}

	// Read the response from the service
	response, err := http.ReadResponse(bufio.NewReader(remoteConnection), req)
	if err != nil {
		proxy.replyWithError(500, err, wr)
		return
	}

	var responseBody []byte
	if response.Body != nil {
		// Read the request body
		responseBody, err = io.ReadAll(response.Body)
		if err != nil {
			proxy.replyWithError(500, err, wr)
			return
		}
	}

	wr.WriteHeader(response.StatusCode)

	responseHeaders := make(map[string]string)
	for k, v := range response.Header {
		wr.Header().Add(k, v[0])
	}

	// Send the response to the client
	_, err = wr.Write(responseBody)
	if err != nil {
		// TODO: handle error
		return
	}

	relay := &types.Relay{
		Req: relayRequest,
		Res: &types.RelayResponse{
			StatusCode: int32(response.StatusCode),
			Headers:    responseHeaders,
			Payload:    responseBody,
		},
	}

	relay.Res.Signature = proxy.signResponse(relay)

	proxy.output <- relay
}

func (r *Proxy) signResponse(relay *types.Relay) []byte {
	return nil
}

func (proxy *Proxy) replyWithError(statusCode int, err error, wr http.ResponseWriter) {
	wr.WriteHeader(statusCode)
	wr.Write([]byte(err.Error()))
}
