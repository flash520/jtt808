package protocol

// 查询车辆位置
type T808_0x8201 struct {
}

func (entity *T808_0x8201) MsgID() MsgID {
	return MsgT808_0x8201
}

func (entity *T808_0x8201) Encode() ([]byte, error) {
	return nil, nil
}

func (entity *T808_0x8201) Decode(data []byte) (int, error) {
	return 0, nil
}
