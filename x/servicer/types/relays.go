package types

func (req RelayRequest) GetSignableBytes() ([]byte, error) {
	req.Signature = nil
	return req.Marshal()
}

func (res RelayResponse) GetSignableBytes() ([]byte, error) {
	res.Signature = nil
	return res.Marshal()
}
