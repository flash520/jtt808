package protocol

// 设置终端参数
type T808_0x8103 struct {
	// 参数项列表
	Params []*Param
}

func (entity *T808_0x8103) MsgID() MsgID {
	return MsgT808_0x8103
}

func (entity *T808_0x8103) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入参数数量
	writer.WriteByte(byte(len(entity.Params)))

	// 写入参数列表
	for _, param := range entity.Params {
		// 写入参数ID
		writer.WriteUint32(param.id)

		// 写入参数长度
		writer.WriteByte(byte(len(param.serialized)))

		// 写入参数数据
		writer.Write(param.serialized)
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x8103) Decode(data []byte) (int, error) {
	if len(data) <= 3 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取参数个数
	paramNums, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取参数信息
	params := make([]*Param, 0, paramNums)
	for i := 0; i < int(paramNums); i++ {
		// 读取参数ID
		id, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}

		// 读取数据长度
		size, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}

		// 读取数据内容
		value, err := reader.Read(int(size))
		if err != nil {
			return 0, err
		}
		params = append(params, &Param{
			id:         id,
			serialized: value,
		})
	}

	entity.Params = params
	return len(data) - reader.Len(), nil
}
