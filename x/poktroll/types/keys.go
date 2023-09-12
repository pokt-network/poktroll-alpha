package types

const (
	// ModuleName defines the module name
	ModuleName = "poktroll"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_poktroll"

	// Key to store servicer-related info
	ServicerPrefix = "servicer"

	// Key to store watcher-related info
	WatcherPrefix = "watcher"

	// Key to store portal-related info
	PortalPrefix = "portal"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
