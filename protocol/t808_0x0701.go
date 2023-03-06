package protocol

// 电子运单上报
type T808_0x0701 struct {
	// 电子运单长度
	Size uint32
	// 电子运单内容
	Content []byte
}

func (entity *T808_0x0701) MsgID() MsgID {
	return MsgT808_0x0701
}

func (entity *T808_0x0701) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入运单内容长度
	writer.WriteUint32(entity.Size)

	// 写入运单内容
	writer.Write(entity.Content)
	return writer.Bytes(), nil
}

func (entity *T808_0x0701) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取运单内容长度
	var err error
	entity.Size, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	//  读取运单内容
	entity.Content, err = reader.Read()
	if err != nil {
		return 0, err
	}
	return len(data) - reader.Len(), nil
}
