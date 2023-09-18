package types

type Session interface {
	GetSessionNumber() uint64
	GetSessionHeight() uint64
	GetBlockHash() []byte
}
