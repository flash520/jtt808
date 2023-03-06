package protocol

// 文本信息下发
type T808_0x8300 struct {
	// 标志
	Flag byte
	// 文本信息
	Text string
}

func (entity *T808_0x8300) MsgID() MsgID {
	return MsgT808_0x8300
}

func (entity *T808_0x8300) Encode() ([]byte, error) {
	writer := NewWriter()
	writer.WriteByte(entity.Flag)
	if len(entity.Text) > 0 {
		if err := writer.WritString(entity.Text); err != nil {
			return nil, err
		}
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x8300) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, ErrInvalidBody
	}

	entity.Flag = data[0]
	reader := NewReader(data[1:])
	if reader.Len() > 0 {
		data, err := reader.ReadString()
		if err != nil {
			return 0, err
		}
		entity.Text = data
	}
	return len(data) - reader.Len() - 1, nil
}
