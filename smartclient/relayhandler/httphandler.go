package relayhandler

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"poktroll/utils"
	svcTypes "poktroll/x/service/types"
	"poktroll/x/servicer/types"
)

// handleHTTPRelays handles HTTP relay requests (currently JSON_RPC)
// it takes the http ServeHTTP arguments, the service id, and the rpc type
// to construct a RelayRequest and send it to the relayer
func (relayHandler *RelayHandler) handleHTTPRelays(
	w http.ResponseWriter,
	req *http.Request,
	serviceId string,
	rpcType svcTypes.RPCType,
) {
	// get the current session for the service
	session := relayHandler.getServiceCurrentSession(serviceId)
	headers := cloneHeaders(req.Header)

	// if the request has a body, read it into a byte slice and use it as the payload
	var payload []byte = nil
	if req.Body != nil {
		var err error
		payload, err = io.ReadAll(req.Body)
		if err != nil {
			// reply to the client with an error if the body could not be read
			// TODO: derive a client facing error from the error
			utils.ReplyWithHTTPError(500, err, w)
		}
	}

	relayRequest := &types.RelayRequest{
		Headers:            headers,
		Method:             req.Method,
		Url:                req.URL.String(),
		Payload:            payload,
		SessionId:          session.SessionId,
		ApplicationAddress: relayHandler.applicationAddress,
	}

	// update signer if not already present
	// NOTE: this can only be nil when the signer is a ring signer
	if relayHandler.signer == nil && relayHandler.signingKey != nil {
		if err := relayHandler.updateSinger(); err != nil {
			utils.ReplyWithHTTPError(500, err, w)
			return
		}
	}
	signature, err := signRelayRequest(relayRequest, relayHandler.signer)
	if err != nil {
		utils.ReplyWithHTTPError(500, err, w)
		return
	}

	relayRequest.ApplicationSignature = signature
	relayRequestBz, err := relayRequest.Marshal()
	if err != nil {
		utils.ReplyWithHTTPError(500, err, w)
		return
	}

	// get a relayer endpoint to send the request to
	relayerUrl := relayHandler.getSessionRelayerUrl(session, rpcType)
	if relayerUrl == "" {
		utils.ReplyWithHTTPError(500, errNoRelayerUrl, w)
	}

	// create a reader from the relay request bytes and send it to the relayer
	relayRequestReader := bytes.NewReader(relayRequestBz)

	// perform the relay request and get the http.Response
	relayHTTPResponse, err := http.DefaultClient.Post(relayerUrl, "application/json", relayRequestReader)
	if err != nil {
		utils.ReplyWithHTTPError(500, err, w)
		return
	}

	// read the response body into a byte slice
	relayResponseBz, err := io.ReadAll(relayHTTPResponse.Body)
	if err != nil {
		utils.ReplyWithHTTPError(500, err, w)
		return
	}

	// verify that the relay response is properly formed and signed by the relayer then get the RelayResponse
	relayResponse, err := getVerifiedRelayResponse(relayResponseBz)
	if err != nil {
		utils.ReplyWithHTTPError(500, err, w)
		return
	}

	// send back the relay original response to the client with the same status code and headers
	// that were wrapped by the relayer into the RelayResponse
	w.WriteHeader(int(relayResponse.StatusCode))
	for key, values := range relayResponse.Headers {
		headerValues := strings.Split(values, ", ")
		for _, v := range headerValues {
			w.Header().Add(key, v)
		}
	}

	// write the relay response payload (original response) to the client
	if _, err := w.Write(relayResponse.Payload); err != nil {
		utils.ReplyWithHTTPError(500, err, w)
		return
	}
}

// TODO_REFACTOR: Move this into a shared utils (with the relayer) and reuse
// cloneHeaders clones the headers map and joins the values of each header with a comma
// to comply with the RelayRequest protobuf definition
func cloneHeaders(headers map[string][]string) map[string]string {
	clonedHeaders := make(map[string]string)
	for key, value := range headers {
		clonedHeaders[key] = strings.Join(value, ", ")
	}

	return clonedHeaders
}
