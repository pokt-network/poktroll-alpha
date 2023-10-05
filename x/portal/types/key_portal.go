package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// PortalKeyPrefix is the prefix to retrieve all Portals
	PortalKeyPrefix = "Portal/value/"
	// PortalDelegationsKeyPrefix is the prefix to retrieve all Portal PubKeys an Application is delegated to
	PortalDelegationsKeyPrefix = "Portal/delegatedPortals/"
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

// PortalDelegationsKey returns the store key to retrieve all Portals and app is delegated to from the index fields
func PortalDelegationsKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
