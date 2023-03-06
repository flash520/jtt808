package protocol

// 行驶记录参数下传命令
type T808_0x8701 struct {
	// 命令字
	Cmd byte
	// 数据块
	Data []byte
}

func (entity *T808_0x8701) MsgID() MsgID {
	return MsgT808_0x8701
}

func (entity *T808_0x8701) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入命令字
	writer.WriteByte(entity.Cmd)

	// 写入数据块
	writer.Write(entity.Data)
	return writer.Bytes(), nil
}

func (entity *T808_0x8701) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取命令字
	var err error
	entity.Cmd, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取数据块
	entity.Data, err = reader.Read()
	if err != nil {
		return 0, err
	}
	return len(data) - reader.Len(), nil
}
