package protocol

// 人工确认报警消息
type T808_0x8203 struct {
	// 报警消息流水号
	MsgSerialNo uint16
	// 人工确认报警类型
	Type uint32
}

func (entity *T808_0x8203) MsgID() MsgID {
	return MsgT808_0x8203
}

func (entity *T808_0x8203) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入消息序列号
	writer.WriteUint16(entity.MsgSerialNo)

	// 写入报警类型
	writer.WriteUint32(entity.Type)
	return writer.Bytes(), nil
}

func (entity *T808_0x8203) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取消息序列号
	var err error
	entity.MsgSerialNo, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取报警类型
	entity.Type, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}
	return len(data) - reader.Len(), nil
}
