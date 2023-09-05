package types

type Session struct {
	Id                 string
	BlockHeight        uint64
	BlockHash          string
	SessionNumber      uint64
	SessionBlockHeight uint64
}
