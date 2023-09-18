package types

type Block interface {
	Height() uint64
	Hash() []byte
}
