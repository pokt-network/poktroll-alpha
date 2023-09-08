package types

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func (r *Relay) Serialize() []byte {
	return []byte(fmt.Sprintf(
		"%d:%s:%s:%d:%s",
		r.Req.Height,
		r.Req.Req.Url,
		string(r.Res.Payload),
		r.Res.StatusCode,
		string(r.Res.Signature),
	))
}

// Validate returns nil if every check has passed
func (r *RelayRequest) Validate() error {
	log.Printf("validating request: %+v", r)
	if r.Height == 0 {
		return fmt.Errorf("ErrInvalidHeight")
	}
	if r.Req == nil {
		return fmt.Errorf("ErrInvalidRequest")
	}
	return nil
}

// Execute runs the relay request after the relay has been validated and
// returns both the payload and the error, if any
func (r *RelayRequest) Execute() *RelayResponse {
	switch r.Req.Method {
	case http.MethodGet:
		res, err := http.Get(r.Req.Url)
		if err != nil {
			return &RelayResponse{
				Err: err.Error(),
			}
		}
		defer res.Body.Close()
		// read off payload body
		payload, err := io.ReadAll(res.Body)
		relayRes := &RelayResponse{
			StatusCode: int32(res.StatusCode),
			Payload:    payload,
			Err:        err.Error(),
		}
		return relayRes
	case http.MethodPost:
		log.Fatalf("todo: support request method type %s", r.Req.Method)
	}
	return &RelayResponse{
		Err: fmt.Sprintf("todo: support request method type %s", r.Req.Method),
	}
}
