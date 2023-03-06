package protocol

import "time"

// 存储多媒体数据检索
type T808_0x8802 struct {
	// 多媒体类型
	Type T808_0x0800_MediaType
	// 通道 ID
	ChannelID byte
	// 事件项编码
	// 0：平台下发指令； 1：定时动作； 2：抢劫报警触发；3：碰撞侧翻报警触发；其他保留
	Event byte
	// 起始时间
	StartTime time.Time
	// 结束时间
	EndTime time.Time
}

func (entity *T808_0x8802) MsgID() MsgID {
	return MsgT808_0x8802
}

func (entity *T808_0x8802) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入媒体类型
	writer.WriteByte(byte(entity.Type))

	// 写入通道ID
	writer.WriteByte(entity.ChannelID)

	// 写入事件项编码
	writer.WriteByte(entity.Event)

	// 写入开始时间
	writer.WriteBcdTime(entity.StartTime)

	// 写入结束时间
	writer.WriteBcdTime(entity.EndTime)
	return writer.Bytes(), nil
}

func (entity *T808_0x8802) Decode(data []byte) (int, error) {
	if len(data) < 15 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取媒体类型
	mediaType, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	entity.Type = T808_0x0800_MediaType(mediaType)

	// 读取通道ID
	entity.ChannelID, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取事件项编码
	entity.Event, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取开始时间
	entity.StartTime, err = reader.ReadBcdTime()
	if err != nil {
		return 0, err
	}

	// 读取结束时间
	entity.EndTime, err = reader.ReadBcdTime()
	if err != nil {
		return 0, err
	}
	return len(data) - reader.Len(), nil
}
