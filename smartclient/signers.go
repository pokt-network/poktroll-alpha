package smartclient

import (
	"fmt"

	ring_types "github.com/athanorlabs/go-dleq/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	ring "github.com/noot/ring-go"
)

type Signer interface {
	Sign(relayRequest [32]byte) (signature []byte, err error)
}

// SimpleSigner implements the Signer interface using a keyring and a key name
type SimpleSigner struct {
	keyring keyring.Keyring
	keyName string
	pubKey  types.PubKey
}

// NewSimpleSigner creates a new SimpleSigner and stores the public key to be used for potential verifications
func NewSimpleSigner(keyring keyring.Keyring, keyName string) *SimpleSigner {
	keyRecord, err := keyring.Key(keyName)
	if err != nil {
		panic(fmt.Errorf("key not found: %v", err))
	}

	pubKey, err := keyRecord.GetPubKey()
	if err != nil {
		panic(fmt.Errorf("failed to get pubKey: %v", err))
	}

	return &SimpleSigner{keyring: keyring, keyName: keyName, pubKey: pubKey}
}

// Sign implements the Signer interface
func (signer *SimpleSigner) Sign(data [32]byte) (signature []byte, err error) {
	sig, _, err := signer.keyring.Sign(signer.keyName, data[:])
	return sig, err
}

// RingSigner implements the Signer interface using a ring and a scalar point on the ring's curve
type RingSigner struct {
	ring    *ring.Ring
	privKey ring_types.Scalar
}

// NewRingSigner creates a new RingSigner instance with the ring and private key provided
func NewRingSigner(ring *ring.Ring, privKey ring_types.Scalar) *RingSigner {
	return &RingSigner{ring: ring, privKey: privKey}
}

// Sign uses the ring and private key to sign the message provided and returns the
// serialised ring signature that can be deserialised and verified by the verifier
func (r *RingSigner) Sign(message [32]byte) ([]byte, error) {
	ringSig, err := r.ring.Sign(message, r.privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}
	return ringSig.Serialize()
}
