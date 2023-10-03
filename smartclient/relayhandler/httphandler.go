package relayhandler

import (
	"bytes"
	"crypto/sha256"
	"io"
	"net/http"
	"strings"

	"poktroll/utils"
	"poktroll/x/servicer/types"
)

func (relayHandler *RelayHandler) handleHTTPRelays(w http.ResponseWriter, req *http.Request, serviceId string) {
	relayHandler.sessionMutex.RLock()
	defer relayHandler.sessionMutex.RUnlock()

	headers := make(map[string]string)
	for key, value := range req.Header {
		headers[key] = strings.Join(value, ", ")
	}

	var payload []byte = nil
	if req.Body != nil {
		var err error
		payload, err = io.ReadAll(req.Body)
		if err != nil {
			utils.ReplyWithHTTPError(500, err, w)
		}
	}

	relayRequest := &types.RelayRequest{
		Headers:            headers,
		Method:             req.Method,
		Url:                req.URL.String(),
		Payload:            payload,
		SessionId:          relayHandler.currentSession.SessionId,
		ApplicationAddress: relayHandler.currentSession.Application.Address,
	}

	relayRequestBz, err := relayRequest.Marshal()
	if err != nil {
		utils.ReplyWithHTTPError(500, err, w)
		return
	}

	relayRequestHash := sha256.Sum256(relayRequestBz)
	relaySig, err := relayHandler.signer.Sign(relayRequestHash)
	if err != nil {
		utils.ReplyWithHTTPError(500, err, w)
		return
	}

	relayRequest.ApplicationSignature = relaySig
	relayRequestBz, err = relayRequest.Marshal()
	if err != nil {
		utils.ReplyWithHTTPError(500, err, w)
		return
	}

	servicerUrl := relayHandler.getSessionConnectionInfo(serviceId).Url
	if err != nil {
		utils.ReplyWithHTTPError(500, err, w)
		return
	}

	relayRequestReader := bytes.NewReader(relayRequestBz)
	relayHTTPResponse, err := http.DefaultClient.Post(servicerUrl, "application/json", relayRequestReader)
	if err != nil {
		utils.ReplyWithHTTPError(500, err, w)
		return
	}

	relayResponseBz, err := io.ReadAll(relayHTTPResponse.Body)
	if err != nil {
		utils.ReplyWithHTTPError(500, err, w)
		return
	}

	var relayResponse types.RelayResponse
	err = relayResponse.Unmarshal(relayResponseBz)
	if err != nil {
		utils.ReplyWithHTTPError(500, err, w)
		return
	}

	sig := relayResponse.ServicerSignature
	relayResponse.ServicerSignature = nil
	relayResponseBz, err = relayResponse.Marshal()
	relayResponseHash := sha256.Sum256(relayResponseBz)
	if relayHandler.signer.Verify(relayResponseHash, sig) {
		utils.ReplyWithHTTPError(500, err, w)
		return
	}

	w.WriteHeader(int(relayResponse.StatusCode))

	for k, v := range relayResponse.Headers {
		w.Header().Add(k, v)
	}

	if _, err := w.Write(relayResponse.Payload); err != nil {
		utils.ReplyWithHTTPError(500, err, w)
		return
	}
}
