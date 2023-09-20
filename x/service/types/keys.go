package types

const (
	// ModuleName defines the module name
	ModuleName = "services" // TODO_REFACTOR: Use `services` when scaffolding this module BECAUSE otherwise `assertNoCommonPrefix`  throws a panic

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_service"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
