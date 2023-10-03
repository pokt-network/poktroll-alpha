package smartclient

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/types"
)

type SimpleSigner struct {
	keyring keyring.Keyring
	keyName string
	pubKey  types.PubKey
}

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

func (signer *SimpleSigner) Sign(data [32]byte) (signature []byte, err error) {
	sig, _, err := signer.keyring.Sign(signer.keyName, data[:])
	return sig, err
}

func (signer *SimpleSigner) Verify(data [32]byte, signature []byte) bool {
	return signer.pubKey.VerifySignature(data[:], signature)
}
