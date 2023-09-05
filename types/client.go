package types

// TODO: move these to more appripriate places and/or define as protobufs

type Block struct {
	Height uint64
	Hash   string
}

type TxResult struct {
	Hash   string
	Height uint64
}
