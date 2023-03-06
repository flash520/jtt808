package protocol

// 数据下行透传
type T808_0x8900 struct {
	// 透传消息类型
	Type byte
	// 透传消息内容
	Data []byte
}

func (entity *T808_0x8900) MsgID() MsgID {
	return MsgT808_0x8900
}

func (entity *T808_0x8900) Encode() ([]byte, error) {
	writer := NewWriter()
	writer.WriteByte(entity.Type)
	writer.Write(entity.Data)
	return writer.Bytes(), nil
}

func (entity *T808_0x8900) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, ErrInvalidBody
	}
	entity.Type, entity.Data = data[0], data[1:]
	return len(data), nil
}
