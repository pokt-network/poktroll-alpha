package types

type Session interface {
	SessionNumber() uint64
	SessionHeight() uint64
	BlockHash() []byte
}
