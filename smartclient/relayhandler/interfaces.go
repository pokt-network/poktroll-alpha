package relayhandler

import svcTypes "poktroll/x/service/types"

type ServicesEndpoints *map[string][]svcTypes.Endpoint

type EndpointSelectionStrategy interface {
	SelectEndpoint(endpoints []svcTypes.Endpoint) *svcTypes.Endpoint
}

type Signer interface {
	Sign(relayRequest [32]byte) (signature []byte, err error)
}
