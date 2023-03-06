package protocol

// 终端注销
type T808_0x0003 struct {
}

func (entity *T808_0x0003) MsgID() MsgID {
	return MsgT808_0x0003
}

func (entity *T808_0x0003) Encode() ([]byte, error) {
	return nil, nil
}

func (entity *T808_0x0003) Decode(data []byte) (int, error) {
	return 0, nil
}
