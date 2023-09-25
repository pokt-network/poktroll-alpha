package keeper

// HACK/IMPROVE: using "legacy" errors to save time; replace with custom error
// protobuf types. See: https://docs.cosmos.network/v0.47/core/events.
const (
	EventTypeClaim          = "claim"
	AttributeKeySmtRootHash = "smt_root_hash"
)
