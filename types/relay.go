package types

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Relay struct {
	RelayRequest
	RelayResponse
}

func (r *Relay) Serialize() []byte {
	return []byte(fmt.Sprintf(
		"%d:%s:%s:%d:%s",
		r.Height,
		r.Req.URL.String(),
		string(r.Payload),
		r.StatusCode,
		string(r.Signature),
	))
}

// RelayRequest is validated to trigger a relays and record the response
type RelayRequest struct {
	Height uint64
	Req    *http.Request
}

// RelayResponse returns a payload and status code from a request
type RelayResponse struct {
	Payload    []byte
	StatusCode int
	Err        error
	Signature  []byte
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
		res, err := http.Get(r.Req.URL.String())
		if err != nil {
			return &RelayResponse{
				Err: err,
			}
		}
		defer res.Body.Close()
		// read off payload body
		payload, err := io.ReadAll(res.Body)
		relayRes := &RelayResponse{
			StatusCode: res.StatusCode,
			Payload:    payload,
			Err:        err,
		}
		return relayRes
	case http.MethodPost:
		log.Fatalf("todo: support request method type %s", r.Req.Method)
	}
	return &RelayResponse{
		Err: fmt.Errorf("todo: support request method type %s", r.Req.Method),
	}
}
