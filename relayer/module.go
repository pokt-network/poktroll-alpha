package relayer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

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
	logger *modules.Logger
	// input receives relays
	input chan *types.Relay
	// channel signaling a output has been completed.
	output chan *types.Relay
}

// NewRelayerModule creates a new input and output channel
// for handling a FIFO queue of relay requests and binds them
// to a new RelayerModule.
func NewRelayerModule() modules.RelayerModule {
	return &relayer{
		input:  make(chan *types.Relay),
		output: make(chan *types.Relay),
	}

	return &relayer
}

func (r *relayer) Hydrate(injector *di.Injector, path *[]string) {
	globalLogger := di.Hydrate(modules.LoggerModuleToken, injector, path)
	r.logger = globalLogger.CreateLoggerForModule(modules.RelayerToken.Id())
}

func (r *relayer) CascadeStart() error {
	return r.Start()
}

// Start begins listening for requests on its given pipeline and handles them appropriately
func (r *relayer) Start() error {
	go r.start(context.TODO())
	return nil
}

// start listens for relays on the input channel, validates, then executes them if they're
// valid and outputs them on the relays channel when they're completed.
// it is a block channel that is meant to be called in a goroutine.
func (r *relayer) start(ctx context.Context) {
	// receive incoming relay requests
	for relay := range r.input {
		if err := r.validate(relay); err != nil {
			// update the error response and output the relay.
			relay.Res = &types.RelayResponse{
				Err:        fmt.Errorf("ErrRelayFailedValidation: %w", err).Error(),
				Payload:    nil,
				StatusCode: 400,
			}
			r.output <- relay
		}

		// Execute the relay, return the response and error
		completed := r.execute(relay)
		r.output <- completed

		// store the relay
		go func(relay *types.Relay) {
			if err := r.store(relay); err != nil {
				r.logger.Err(err).Msg("failed to store relay")
			}
		}(relay)
	}
}

// Returns a single-listener channel that emits completed relays.
func (r *relayer) Relays() <-chan *types.Relay {
	return r.output
}

func (r *relayer) Stop() error {
	// TODO: what needs to happen here?
	return nil
}

// validate checks if a relay request is formatted correctly and meets requirements
func (r *relayer) validate(relay *types.Relay) error {
	if relay.Req.Height < 1 {
		return fmt.Errorf("ErrInvalidRelayHeight")
	}
	// DISCUSS is this too heavy handed to parse it through Go's std lib?
	_, err := url.Parse(relay.Req.Req.Url)
	if err != nil {
		return fmt.Errorf("ErrInvalidURL: %s", relay.Req.Req.Url)
	}
	return nil
}

// execute routes the relay to the appropriate handler basedon
// its protocol and runs the given relay after it has been validated.
// it returns the finished relay, with any errors that occurred during
// execution.
func (r *relayer) execute(relay *types.Relay) *types.Relay {
	switch relay.Req.Req.Method {
	case http.MethodGet:
		return r.handleHTTPGet(relay)
	case http.MethodPost:
		// TODO: handleJSONRPC()
		// TODO: handleHTTPPost()
	}
	return &types.Relay{
		Res: &types.RelayResponse{Err: "method type not supported"},
	}
}

// store persists the relay to disk for later reference
func (r *relayer) store(*types.Relay) error {
	return fmt.Errorf("not impl")
}

// handleHTTPGet executes an HTTP GET call to the specified relay's URL
// and returns the relay response and payload
func (r *relayer) handleHTTPGet(relay *types.Relay) *types.Relay {
	res, err := http.Get(relay.Req.Req.Url)
	if err != nil {
		return &types.Relay{
			Req: relay.Req,
			Res: &types.RelayResponse{
				Err:        err.Error(),
				StatusCode: 500,
			},
		}
	}
	defer res.Body.Close()

	payload, err := io.ReadAll(res.Body)
	if err != nil {
		return &types.Relay{
			Req: relay.Req,
			Res: &types.RelayResponse{
				Payload:    payload,
				Err:        err.Error(),
				StatusCode: 500,
				// TODO handle signature correctly
			},
		}
	}

	return &types.Relay{
		Req: relay.Req,
		Res: &types.RelayResponse{
			Payload:    payload,
			StatusCode: 200,
			// TODO handle signature correctly
		},
	}
}
