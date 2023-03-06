package protocol

// 临时位置跟踪控制
type T808_0x8202 struct {
	// 时间间隔
	// 单位为秒 （s），0 则停止跟踪。 停止跟踪无需带后继字段
	Interval uint16
	// 位置跟踪有效期
	// 单位为秒 （s），终端在接收到位置跟踪控制消息后，
	// 在有效期截止时间之前，依据消息中的时间间隔发
	// 送位置汇报
	Expire uint32
}

func (entity *T808_0x8202) MsgID() MsgID {
	return MsgT808_0x8202
}

func (entity *T808_0x8202) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入时间间隔
	writer.WriteUint16(entity.Interval)

	// 写入跟踪有效期
	writer.WriteUint32(entity.Expire)
	return writer.Bytes(), nil
}

func (entity *T808_0x8202) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取时间间隔
	var err error
	entity.Interval, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取跟踪有效期
	entity.Expire, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}
	return len(data) - reader.Len(), nil
}
