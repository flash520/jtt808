package protocol

// 事件报告
type T808_0x0301 struct {
	// 事件 ID
	EventID byte
}

func (entity *T808_0x0301) MsgID() MsgID {
	return MsgT808_0x0301
}

func (entity *T808_0x0301) Encode() ([]byte, error) {
	return []byte{entity.EventID}, nil
}

func (entity *T808_0x0301) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, ErrInvalidBody
	}
	entity.EventID = data[0]
	return 1, nil
}
