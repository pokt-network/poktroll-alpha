package relayer

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

type Relayer struct {
	localAddr        string
	serviceAddr      string
	logger           *log.Logger
	output           chan *types.Relay
	outputObservable utils.Observable[*types.Relay]
}

func NewRelayer(logger *log.Logger) *Relayer {
	r := &Relayer{output: make(chan *types.Relay), logger: logger}
	r.outputObservable, _ = utils.NewControlledObservable[*types.Relay](r.output)

	r.localAddr = "localhost:8545"
	r.serviceAddr = "localhost:8546"

	go r.listen()

	return r
}

func (r *Relayer) Relays() utils.Observable[*types.Relay] {
	return r.outputObservable
}

func (r *Relayer) listen() {
	if err := http.ListenAndServe(r.localAddr, r); err != nil {
		r.logger.Fatal(err)
	}
}

func (r *Relayer) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
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
			r.replyWithError(500, err, wr)
			return
		}
		relayRequest.Payload = requestBody
	}

	// Change the request host to the service address
	req.Host = r.serviceAddr
	req.URL.Host = r.serviceAddr
	req.Body = io.NopCloser(bytes.NewBuffer(relayRequest.Payload))

	// Connect to the service
	remoteConnection, err := net.Dial("tcp", r.serviceAddr)
	if err != nil {
		r.replyWithError(500, err, wr)
		return
	}
	defer remoteConnection.Close()

	// Send the request to the service
	err = req.Write(remoteConnection)
	if err != nil {
		r.replyWithError(500, err, wr)
		return
	}

	// Read the response from the service
	response, err := http.ReadResponse(bufio.NewReader(remoteConnection), req)
	if err != nil {
		r.replyWithError(500, err, wr)
		return
	}

	var responseBody []byte
	if response.Body != nil {
		// Read the request body
		responseBody, err = io.ReadAll(response.Body)
		if err != nil {
			r.replyWithError(500, err, wr)
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

	relay.Res.Signature = r.signResponse(relay)

	r.output <- relay
}

func (r *Relayer) signResponse(relay *types.Relay) []byte {
	return nil
}

func (r *Relayer) replyWithError(statusCode int, err error, wr http.ResponseWriter) {
	wr.WriteHeader(statusCode)
	wr.Write([]byte(err.Error()))
}
