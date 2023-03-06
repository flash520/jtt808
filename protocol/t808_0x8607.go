package protocol

// 删除路线
type T808_0x8607 struct {
	// 路线ID列表
	IDs []uint32
}

func (entity *T808_0x8607) MsgID() MsgID {
	return MsgT808_0x8607
}

func (entity *T808_0x8607) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入路线数
	writer.WriteByte(byte(len(entity.IDs)))

	// 写入路线ID列表
	for _, id := range entity.IDs {
		writer.WriteUint32(id)
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x8607) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取路线数
	count, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取路线ID列表
	for i := 0; i < int(count); i++ {
		id, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		entity.IDs = append(entity.IDs, id)
	}
	return len(data) - reader.Len(), nil
}
