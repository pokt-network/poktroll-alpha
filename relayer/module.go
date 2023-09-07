package relayer

import (
	"fmt"

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
	// channel signaling a relay has been completed.
	relay chan *types.Relay
}

func NewRelayerModule() modules.RelayerModule {
	return &relayer{
		in:    make(chan *types.Relay),
		relay: make(chan *types.Relay),
	}
}

func (r *relayer) Resolve(injector *di.Injector, path *[]string) {}

func (r *relayer) CascadeStart() error {
	return r.Start()
}

// Start begins listening for requests on its given pipeline and handles them appropriately
func (r *relayer) Start() error {
	go r.start()
	return nil
}

func (r *relayer) start() {
	// receive incoming relay requests
	for relay := range r.in {
		if err := r.validate(relay); err != nil {
			// update the error response and output the relay.
			relay.Res = &types.RelayResponse{
				Err:        fmt.Errorf("ErrRelayFailedValidation: %w", err).Error(),
				Payload:    nil,
				StatusCode: 400,
			}
			r.relay <- relay
		}

		// Execute the relay, return the response and error
		completed := r.execute(relay)
		r.relay <- completed

		// store the relay
		go func(relay *types.Relay) {
			if err := r.store(relay); err != nil {
				fmt.Printf("TODO ADD A LOGGER LOL")
			}
		}(relay)
	}
}

// Returns a single-listener channel that emits completed relays.
func (r *relayer) Relays() <-chan *types.Relay {
	return r.relay
}

func (r *relayer) Stop() error {
	// TODO: what needs to happen here?
	return nil
}

// validate checks if a relay request is formatted correctly and meets requirements
func (r *relayer) validate(relay *types.Relay) error {
	return fmt.Errorf("not impl")
}

// execute runs the relay and returns the error and response, no matter what
func (r *relayer) execute(relay *types.Relay) *types.Relay {
	return &types.Relay{
		Res: &types.RelayResponse{
			Err: "not impl",
		},
	}
}

// store persists the relay to disk for later reference
func (r *relayer) store(*types.Relay) error {
	return fmt.Errorf("not impl")
}
