package relayhandler

// TODO_REFACTOR: Avoid an `interfaces` file. Each major types/actor/module/submodule should have its own
// file and the interface should be there.

import svcTypes "poktroll/x/service/types"

type ServicesEndpoints *map[string][]svcTypes.Endpoint

type EndpointSelectionStrategy interface {
	SelectEndpoint(endpoints []svcTypes.Endpoint) *svcTypes.Endpoint
}
