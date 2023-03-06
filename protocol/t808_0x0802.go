package protocol

// 存储多媒体数据检索应答
type T808_0x0802 struct {
	// 应答流水号
	ReplyMsgSerialNo uint16
	// 检索项
	Items []T808_0x0802_Item
}

// 多媒体检索项
type T808_0x0802_Item struct {
	// 多媒体 ID
	MediaID uint32
	// 多媒体类型
	Type T808_0x0800_MediaType
	// 通道 ID
	ChannelID byte
	// 事件项编码
	// 0：平台下发指令； 1：定时动作； 2：抢劫报警触 发；3：碰撞侧翻报警触发；其他保留
	Event byte
	// 位置信息汇报
	Location T808_0x0200
}

func (entity *T808_0x0802) MsgID() MsgID {
	return MsgT808_0x0802
}

func (entity *T808_0x0802) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入应答流水号
	writer.WriteUint16(entity.ReplyMsgSerialNo)

	// 写入多媒体检索项数
	writer.WriteUint16(uint16(len(entity.Items)))

	// 写入多媒体检索项
	for _, item := range entity.Items {
		// 写入多媒体ID
		writer.WriteUint32(item.MediaID)

		// 写入多媒体类型
		writer.WriteByte(byte(item.Type))

		// 写入通道ID
		writer.WriteByte(item.ChannelID)

		// 写入事件项编码
		writer.WriteByte(item.Event)

		// 写入定位信息
		item.Location.Extras = nil
		data, err := item.Location.Encode()
		if err != nil {
			return nil, err
		}
		writer.Write(data)
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x0802) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取应答流水号
	var err error
	entity.ReplyMsgSerialNo, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取多媒体检索项数
	count, err := reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取多媒体检索项
	entity.Items = make([]T808_0x0802_Item, 0, count)
	for i := 0; i < int(count); i++ {
		var item T808_0x0802_Item

		// 读取多媒体ID
		item.MediaID, err = reader.ReadUint32()
		if err != nil {
			return 0, err
		}

		// 读取多媒体类型
		mediaType, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}
		item.Type = T808_0x0800_MediaType(mediaType)

		// 读取通道ID
		item.ChannelID, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}

		// 读取事件项编码
		item.Event, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}

		// 读取定位信息
		buf, err := reader.Read(28)
		if err != nil {
			return 0, err
		}
		_, err = item.Location.Decode(buf)
		if err != nil {
			return 0, err
		}
		entity.Items = append(entity.Items, item)
	}
	return len(data) - reader.Len(), nil
}
