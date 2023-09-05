package relayer

import (
	"net/http"
	"net/url"
	"time"

	"poktroll/modules"
	"poktroll/runtime/di"
	"poktroll/types"
)

/////////////
// RELAYER //
/////////////

// relayer receives Relays from in and executes the relay
// if passes validation. The request and response are then
// stored by height.
type relayer struct {
	// in receives relays
	in chan *types.Relay
	// relays are sent on out after they're done being processed
	out chan *types.Relay
	// channel signaling a relay has been completed
	relay chan *types.Relay
}

func NewRelayerModule() modules.RelayerModule {
	return &relayer{
		in:    make(chan *types.Relay),
		out:   make(chan *types.Relay),
		relay: make(chan *types.Relay),
	}
}

func (r *relayer) Resolve(injector *di.Injector, path *[]string) {}

func (r *relayer) CascadeStart() error {
	return r.Start()
}

// Start begins listening for requests on its given pipeline and handles them appropriately
func (r *relayer) Start() error {
	go func() {
		// receive incoming relay requests
		//for req := range r.in {
		for range time.NewTicker(1000 * time.Millisecond).C {
			//if err := req.Validate(); err != nil {
			//	// respond with error
			//	r.out <- &types.Relay{
			//		RelayRequest: req.RelayRequest,
			//		RelayResponse: types.RelayResponse{
			//			Err: err,
			//		},
			//	}
			//}

			// Execute the relay, return the response and error
			//res := req.Execute()
			//done := &types.Relay{
			//	RelayRequest:  req.RelayRequest,
			//	RelayResponse: *res,
			//}

			// Emit it on the output channel
			// TECHDEBT(dylan): the output channel could be attached to the relay
			// request instead of the relayer to allow for completely concurrent
			// handling of all relays
			//r/.out <- done

			go func() {
				req := http.Request{}
				req.URL, _ = url.Parse("http://localhost:8081")
				r.relay <- &types.Relay{
					types.RelayRequest{
						Height: 0,
						Req:    &req,
					},
					types.RelayResponse{
						Payload:    []byte("123"),
						StatusCode: 200,
						Err:        nil,
						Signature:  []byte("sig"),
					},
				}

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
			}()
		}
	}()
	return nil
}

func (r *relayer) Relays() <-chan *types.Relay {
	return r.relay
}

func (r *relayer) Stop() error {
	// TODO: what needs to happen here?
	return nil
}
