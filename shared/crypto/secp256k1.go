package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"

	cosmosSecp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
)

const (
	AddressLen = 20
)

type (
	Secp256k1PublicKey struct {
		key *cosmosSecp256k1.PubKey
	}
	Secp256k1PrivateKey struct {
		key *cosmosSecp256k1.PrivKey
	}
)

var (
	PublicKeyLen  = cosmosSecp256k1.PubKeySize
	PrivateKeyLen = cosmosSecp256k1.PrivKeySize
)

func NewAddress(hexString string) (Address, error) {
	bz, err := hex.DecodeString(hexString)
	if err != nil {
		return bz, ErrCreateAddress(err)
	}
	return NewAddressFromBytes(bz)
}

func NewAddressFromBytes(bz []byte) (Address, error) {
	bzLen := len(bz)
	if bzLen != AddressLen {
		return bz, ErrInvalidAddressLen(bzLen)
	}
	return bz, nil
}

func (a Address) String() string {
	return hex.EncodeToString(a)
}

func NewPrivateKey(hexString string) (PrivateKey, error) {
	bz, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, ErrCreatePrivateKey(err)
	}
	return NewPrivateKeyFromBytes(bz)
}

func GeneratePrivateKey() (PrivateKey, error) {
	pk := cosmosSecp256k1.GenPrivKey()
	return Secp256k1PrivateKey{key: pk}, nil
}

func GeneratePrivateKeyWithReader(rand io.Reader) (PrivateKey, error) {
	// Secret is hashed by `#GenPrivKeyFromSecret`, using a buffer size which is
	// equivalent to the private key size to provide at least as much entropy.
	secret := make([]byte, cosmosSecp256k1.PrivKeySize)
	if _, err := io.ReadFull(rand, secret); err != nil {
		return nil, err
	}

	pk := cosmosSecp256k1.GenPrivKeyFromSecret(secret)
	return Secp256k1PrivateKey{key: pk}, nil
}

func NewPrivateKeyFromBytes(bz []byte) (PrivateKey, error) {
	bzLen := len(bz)
	if bzLen != cosmosSecp256k1.PrivKeySize {
		return nil, ErrInvalidPrivateKeyLen(bzLen)
	}
	pk := &cosmosSecp256k1.PrivKey{Key: bz}
	return Secp256k1PrivateKey{key: pk}, nil
}

func NewPrivateKeyFromSeed(seed []byte) (PrivateKey, error) {
	if len(seed) < cosmosSecp256k1.PrivKeySize {
		return nil, ErrInvalidPrivateKeySeedLenError(len(seed))
	}
	pk := cosmosSecp256k1.GenPrivKeyFromSecret(seed)
	return Secp256k1PrivateKey{key: pk}, nil
}

var _ PrivateKey = Secp256k1PrivateKey{}

func (priv Secp256k1PrivateKey) Bytes() []byte {
	return priv.key.Bytes()
}

func (priv Secp256k1PrivateKey) String() string {
	return hex.EncodeToString(priv.Bytes())
}

func (priv Secp256k1PrivateKey) Equals(other PrivateKey) bool {
	return priv.key.Equals(other.(Secp256k1PrivateKey).key)
}

func (priv Secp256k1PrivateKey) PublicKey() PublicKey {
	pubKey := priv.key.PubKey().(*cosmosSecp256k1.PubKey)
	return Secp256k1PublicKey{key: pubKey}
}

func (priv Secp256k1PrivateKey) Address() Address {
	publicKey := priv.PublicKey()
	return publicKey.Address()
}

func (priv Secp256k1PrivateKey) Sign(msg []byte) ([]byte, error) {
	return priv.key.Sign(msg)
}

func (priv Secp256k1PrivateKey) Size() int {
	return cosmosSecp256k1.PrivKeySize
}

func (priv Secp256k1PrivateKey) Seed() []byte {
	return priv.key.Key
}

func (priv *Secp256k1PrivateKey) UnmarshalJSON(data []byte) error {
	var privateKey string
	if err := json.Unmarshal(data, &privateKey); err != nil {
		return err
	}
	return priv.UnmarshalText([]byte(privateKey))
}

func (priv *Secp256k1PrivateKey) UnmarshalText(data []byte) error {
	privateKey := string(data)
	keyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return err
	}
	privKey, err := NewPrivateKeyFromBytes(keyBytes)
	if err != nil {
		return err
	}
	priv.key = privKey.(*Secp256k1PrivateKey).key
	return nil
}

var _ PublicKey = Secp256k1PublicKey{}

func NewPublicKey(hexString string) (PublicKey, error) {
	bz, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, ErrCreatePublicKey(err)
	}
	return NewPublicKeyFromBytes(bz)
}

func NewPublicKeyFromBytes(bz []byte) (PublicKey, error) {
	bzLen := len(bz)
	if bzLen != cosmosSecp256k1.PubKeySize {
		return nil, ErrInvalidPublicKeyLen(bzLen)
	}

	pubKey := &cosmosSecp256k1.PubKey{Key: bz}
	return Secp256k1PublicKey{key: pubKey}, nil
}

func (pub Secp256k1PublicKey) Bytes() []byte {
	return pub.key.Bytes()
}

func (pub Secp256k1PublicKey) String() string {
	return hex.EncodeToString(pub.Bytes())
}

func (pub Secp256k1PublicKey) Address() Address {
	hash := sha256.Sum256(pub.key.Key[:])
	return hash[:AddressLen]
}

func (pub Secp256k1PublicKey) Equals(other PublicKey) bool {
	return pub.key.Equals(other.(Secp256k1PublicKey).key)
}

func (pub Secp256k1PublicKey) Verify(msg, sig []byte) bool {
	return pub.key.VerifySignature(msg, sig)
}

func (pub Secp256k1PublicKey) Size() int {
	return pub.key.Size()
}

func GeneratePublicKey() (PublicKey, error) {
	pk, err := GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	return pk.PublicKey(), nil
}

func GenerateAddress() (Address, error) {
	pk, err := GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	return pk.Address(), nil
}

func (pub Secp256k1PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(pub.String())
}

func (pub *Secp256k1PublicKey) UnmarshalJSON(data []byte) error {
	var publicKey string
	if err := json.Unmarshal(data, &publicKey); err != nil {
		return err
	}
	keyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return err
	}
	pubKey, err := NewPublicKeyFromBytes(keyBytes)
	if err != nil {
		return err
	}
	*pub = pubKey.(Secp256k1PublicKey)
	return nil
}
