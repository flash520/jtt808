package protocol

// 提问答案
type T808_0x0302 struct {
	// 应答流水号
	ReplyMsgSerialNo uint16
	// 答案 ID
	AnswerID byte
}

func (entity *T808_0x0302) MsgID() MsgID {
	return MsgT808_0x0302
}

func (entity *T808_0x0302) Encode() ([]byte, error) {
	writer := NewWriter()
	writer.WriteUint16(entity.ReplyMsgSerialNo)
	writer.WriteByte(entity.AnswerID)
	return writer.Bytes(), nil
}

func (entity *T808_0x0302) Decode(data []byte) (int, error) {
	if len(data) < 3 {
		return 0, ErrInvalidBody
	}

	reader := NewReader(data)
	entity.ReplyMsgSerialNo, _ = reader.ReadUint16()
	entity.AnswerID, _ = reader.ReadByte()
	return len(data) - reader.Len(), nil
}
