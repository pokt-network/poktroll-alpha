package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrivateKey_Identity(t *testing.T) {
	privKey, err := GeneratePrivateKey()
	require.NoError(t, err)

	secpPrivKey := privKey.(Secp256k1PrivateKey)

	require.Equal(t, privKey.Bytes(), secpPrivKey.key.Bytes())
	require.Equal(t, privKey.Bytes(), secpPrivKey.key.Key)

	privKey2, err := NewPrivateKeyFromBytes(privKey.Bytes())
	require.NoError(t, err)
	require.Equal(t, privKey2.Bytes(), privKey.Bytes())

	privKey2, err = NewPrivateKeyFromBytes(secpPrivKey.key.Key)
	require.NoError(t, err)
	require.Equal(t, privKey2.Bytes(), privKey.Bytes())

	// NB: the marshalled key is not the same as the raw bytes (RFC 8032)
	keyBz, err := secpPrivKey.key.Marshal()
	require.NoError(t, err)
	require.NotEqual(t, privKey.Bytes(), keyBz)
}

func TestPublicKey_Identity(t *testing.T) {
	privKey, err := GeneratePrivateKey()
	require.NoError(t, err)

	pubKey := privKey.PublicKey()

	secpPubKey := pubKey.(Secp256k1PublicKey)

	require.Equal(t, pubKey.Bytes(), secpPubKey.key.Bytes())
	require.Equal(t, pubKey.Bytes(), secpPubKey.key.Key)

	pubKey2, err := NewPublicKeyFromBytes(pubKey.Bytes())
	require.NoError(t, err)

	pubKey2, err = NewPublicKeyFromBytes(secpPubKey.key.Key)
	require.NoError(t, err)
	require.Equal(t, pubKey2.Bytes(), pubKey.Bytes())

	// NB: the marshalled key is not the same as the raw bytes (RFC 8032)
	keyBz, err := secpPubKey.key.Marshal()
	require.NoError(t, err)
	require.NotEqual(t, pubKey.Bytes(), keyBz)
}
