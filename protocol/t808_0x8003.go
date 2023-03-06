package protocol

// 补传分包请求
type T808_0x8003 struct {
	// 原始消息流水号
	MsgSerialNo uint16
	// 重传包 ID 列表
	PacketIDs []uint16
}

func (entity *T808_0x8003) MsgID() MsgID {
	return MsgT808_0x8003
}

func (entity *T808_0x8003) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入原始消息流水号
	writer.WriteUint16(entity.MsgSerialNo)

	// 写入重传包总数
	writer.WriteByte(byte(len(entity.PacketIDs)))

	// 写入重传包 ID 列表
	for _, id := range entity.PacketIDs {
		writer.WriteUint16(id)
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x8003) Decode(data []byte) (int, error) {
	if len(data) < 3 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取原始消息流水号
	var err error
	entity.MsgSerialNo, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取重传包总数
	count, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取重传包 ID 列表
	for i := 0; i < int(count); i++ {
		id, err := reader.ReadUint16()
		if err != nil {
			return 0, err
		}
		entity.PacketIDs = append(entity.PacketIDs, id)
	}
	return len(data) - reader.Len(), nil
}
