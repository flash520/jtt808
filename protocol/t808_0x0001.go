package protocol

// 终端应答
type T808_0x0001 struct {
	// 应答流水号
	ReplyMsgSerialNo uint16
	// 应答 ID
	ReplyMsgID uint16
	// 结果
	Result Result
}

func (entity *T808_0x0001) MsgID() MsgID {
	return MsgT808_0x0001
}

func (entity *T808_0x0001) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入消息序列号
	writer.WriteUint16(entity.ReplyMsgSerialNo)

	// 写入响应消息ID
	writer.WriteUint16(entity.ReplyMsgID)

	// 写入响应结果
	writer.WriteByte(byte(entity.Result))
	return writer.Bytes(), nil
}

func (entity *T808_0x0001) Decode(data []byte) (int, error) {
	if len(data) < 5 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取消息序列号
	var err error
	entity.ReplyMsgSerialNo, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取响应消息ID
	entity.ReplyMsgID, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取响应结果
	result, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	entity.Result = Result(result)
	return len(data) - reader.Len(), nil
}
