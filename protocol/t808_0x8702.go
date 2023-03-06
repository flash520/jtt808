package protocol

// 上报驾驶员身份信息请求
type T808_0x8702 struct {
}

func (entity *T808_0x8702) MsgID() MsgID {
	return MsgT808_0x8702
}

func (entity *T808_0x8702) Encode() ([]byte, error) {
	return nil, nil
}

func (entity *T808_0x8702) Decode(data []byte) (int, error) {
	return 0, nil
}
