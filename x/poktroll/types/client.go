package types

// TODO: move these to more appripriate places and/or define as protobufs

type Block struct {
	Height int64
	Hash   []byte
}

type TxResult struct {
	Hash   string
	Height int64
}
