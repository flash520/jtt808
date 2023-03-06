package protocol

// 定位数据批量上传
type T808_0x0704 struct {
	// 位置数据类型
	// 0：正常位置批量汇报， 1：盲区补报
	Type byte
	// 位置汇报数据项
	Items []T808_0x0200
}

func (entity *T808_0x0704) MsgID() MsgID {
	return MsgT808_0x0704
}

func (entity *T808_0x0704) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入数据项个数
	writer.WriteUint16(uint16(len(entity.Items)))

	// 写入位置数据类型
	writer.WriteByte(entity.Type)

	// 写入位置汇报数据项
	for _, position := range entity.Items {
		data, err := position.Encode()
		if err != nil {
			return nil, err
		}
		writer.WriteUint16(uint16(len(data)))
		writer.Write(data)
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x0704) Decode(data []byte) (int, error) {
	if len(data) < 3 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取数据项个数
	count, err := reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 写入位置数据类型
	entity.Type, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 写入位置汇报数据项
	entity.Items = make([]T808_0x0200, 0, count)
	for i := 0; i < int(count); i++ {
		size, err := reader.ReadUint16()
		if err != nil {
			return 0, err
		}

		buf, err := reader.Read(int(size))
		if err != nil {
			return 0, err
		}

		var position T808_0x0200
		_, err = position.Decode(buf)
		if err != nil {
			return 0, err
		}
		entity.Items = append(entity.Items, position)
	}
	return len(data) - reader.Len(), nil
}
