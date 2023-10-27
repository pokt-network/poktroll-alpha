package types

func (session *Session) GetSessionEndHeight() uint64 {
	return session.SessionBlockStartHeight + session.NumBlocksPerSession
}
