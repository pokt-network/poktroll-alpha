package relayer

import (
	"log"
	"net/http"
	"net/url"
	"testing"

	"poktroll/types"
)

// localAnvil points to the default URL of a locally run anvil node
var localAnvil = "http://127.0.0.1:8545/"

// TestRelayWorker relies on a local anvil node running at the defualt port
// tl;dr- `make local_anvil` from the command line to start up a default
// local anvil node for simulating ethereum relays
func TestRelayWorker(t *testing.T) {
	t.Run("should handle a single GET relay", func(t *testing.T) {
		worker := setupRelayer(t)
		u, err := url.Parse(localAnvil)
		if err != nil {
			log.Fatal(err)
		}

		// feed a relay in
		worker.in <- &types.Relay{
			Req: &types.RelayRequest{
				Req: &types.HttpRequest{
					Url:    u.String(),
					Method: http.MethodGet,
				},
			},
		}

		// wait for output
		got := <-worker.relay
		if got.GetRes().Err != "" {
			t.Errorf("failed to process relay: %+v", got.GetRes().Err)
		}
	})
	t.Run("should fail with invalid height error", func(t *testing.T) {
		worker := setupRelayer(t)
		u, err := url.Parse(localAnvil)
		if err != nil {
			log.Fatal(err)
		}

		// feed a relay in
		worker.in <- &types.Relay{
			Req: &types.RelayRequest{
				Req: &types.HttpRequest{
					Url: u.String(),
				},
			},
		}

		// wait for output
		got := <-worker.relay
		if got.Res.Err != "ErrInvalidHeight" {
			t.Errorf("got: %+v - wanted: %v", got, "ErrInvalidHeight")
		}
	})
}

func setupRelayer(t *testing.T) *relayer {
	worker := &relayer{
		in:    make(chan *types.Relay),
		relay: make(chan *types.Relay),
	}
	if err := worker.Start(); err != nil {
		log.Fatalf("relayer failed to start: %+v", err)
	}

	return worker
}
