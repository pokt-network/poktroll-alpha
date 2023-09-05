package client

import (
	"context"

	"poktroll/modules"
	"poktroll/runtime/di"
	"poktroll/types"
)

var _ modules.ServicerClient = &servicerClient{}

type servicerClient struct {
	rpcClient modules.RelayerModule
}

func NewServicerClient() modules.ServicerClient {
	return &servicerClient{}
}

func (sc *servicerClient) Resolve(injector *di.Injector, path *[]string) {

}

func (sc *servicerClient) CascadeStart() error {
	panic("implement me")
}

func (sc *servicerClient) Start() error {
	panic("implement me")
}

func (sc *servicerClient) Stop() {
	panic("implement me")
}

func (sc *servicerClient) HandleRelay(
	ctx context.Context,
	relay *types.Relay,
) <-chan types.Maybe[*types.Relay] {
	//sc.relayer
	return nil
}
