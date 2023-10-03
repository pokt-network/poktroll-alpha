package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// PortalKeyPrefix is the prefix to retrieve all Portals
	PortalKeyPrefix = "Portal/value/"
)

// PortalKey returns the store key to retrieve a Portal from the index fields
func PortalKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
