package protocol

// 行驶记录数据上传
type T808_0x0700 struct {
	// 应答流水号
	ReplyMsgSerialNo uint16
	// 命令字
	Cmd byte
	// 数据块
	Data []byte
}

func (entity *T808_0x0700) MsgID() MsgID {
	return MsgT808_0x0700
}

func (entity *T808_0x0700) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入应答流水号
	writer.WriteUint16(entity.ReplyMsgSerialNo)

	// 写入命令字
	writer.WriteByte(entity.Cmd)

	// 写入数据块
	writer.Write(entity.Data)
	return writer.Bytes(), nil
}

func (entity *T808_0x0700) Decode(data []byte) (int, error) {
	if len(data) < 3 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取应答流水号
	var err error
	entity.ReplyMsgSerialNo, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取命令字
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
