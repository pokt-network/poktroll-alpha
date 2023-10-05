package smartclient

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/types"
)

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
