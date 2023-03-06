package protocol

// 查询终端参数
type T808_0x8106 struct {
	// 参数 ID 列表
	Params []uint32
}

func (entity *T808_0x8106) MsgID() MsgID {
	return MsgT808_0x8106
}

func (entity *T808_0x8106) Encode() ([]byte, error) {
	writer := NewWriter()
	writer.WriteByte(byte(len(entity.Params)))
	for _, param := range entity.Params {
		writer.WriteUint32(param)
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x8106) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, ErrInvalidBody
	}

	count := int(data[0])
	reader := NewReader(data[1:])
	entity.Params = make([]uint32, 0, count)
	for i := 0; i < count; i++ {
		id, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		entity.Params = append(entity.Params, id)
	}
	return len(data) - reader.Len() - 1, nil
}
