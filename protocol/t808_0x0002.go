package protocol

// 终端心跳
type T808_0x0002 struct {
}

func (entity *T808_0x0002) MsgID() MsgID {
	return MsgT808_0x0002
}

func (entity *T808_0x0002) Encode() ([]byte, error) {
	return nil, nil
}

func (entity *T808_0x0002) Decode(data []byte) (int, error) {
	return 0, nil
}
