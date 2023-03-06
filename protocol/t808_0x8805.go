package protocol

// 单条存储多媒体数据检索上传命令
type T808_0x8805 struct {
	// 多媒体 ID
	MediaID uint32
	// 删除标志
	RemoveFlag byte
}

func (entity *T808_0x8805) MsgID() MsgID {
	return MsgT808_0x8805
}

func (entity *T808_0x8805) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入媒体ID
	writer.WriteUint32(entity.MediaID)

	// 写入删除标志
	writer.WriteByte(entity.RemoveFlag)
	return writer.Bytes(), nil
}

func (entity *T808_0x8805) Decode(data []byte) (int, error) {
	if len(data) < 3 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取媒体ID
	var err error
	entity.MediaID, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	// 读取删除标志
	entity.RemoveFlag, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}
	return len(data) - reader.Len(), nil
}
