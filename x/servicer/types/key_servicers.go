package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// ServicersKeyPrefix is the prefix to retrieve Servicers
	ServicersKeyPrefix = "Servicers/value/"
	// ClaimsKeyPrefix is the prefix to retrieve Claims
	ClaimsKeyPrefix = "Claims/value/"
	ProofsKeyPrefix = "Proofs/value/"
)

// ServicersKey returns the store key to retrieve a Servicers from the index fields
func ServicersKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
