package protocol

// 查询终端参数
type T808_0x8107 struct {
}

func (entity *T808_0x8107) MsgID() MsgID {
	return MsgT808_0x8107
}

func (entity *T808_0x8107) Encode() ([]byte, error) {
	return nil, nil
}

func (entity *T808_0x8107) Decode(data []byte) (int, error) {
	return 0, nil
}
