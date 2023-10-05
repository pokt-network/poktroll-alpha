package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// ServicersKeyPrefix is the prefix to retrieve Servicers
	ServicersKeyPrefix = "Servicers/value/"
	// TODO_CONSIDERATION: should this be moved to x/servicer/types/key_claims.go?
	// ClaimsKeyPrefix is the prefix to retrieve Claims
	ClaimsKeyPrefix = "Claims/value/"
	// ProofsKeyProfix is the prefix to retrieve Proofs for claims
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

// TODO_CONSIDERATION: should this be moved to x/servicer/types/key_claims.go?
func ClaimsKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
