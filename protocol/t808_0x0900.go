package protocol

// 数据上行透传
type T808_0x0900 struct {
	// 透传消息类型
	// GNSS模块详细定位数据 0x00 GNSS模块详细定位数据
	// 道路运输证 IC 卡信息 0x0B
	// 道路运输证 IC 卡信息上传消息为 64Byte，下传消息为24Byte。道路运输证 IC 卡认证透传超时时间为 30s。 超时后，不重发
	// 串口 1 透传 0x41 串口 1 透传消息
	// 串口 2 透传 0x42 串口 2 透传消息
	// 用户自定义透传 0xF0-0xFF 用户自定义透传消息
	Type byte
	// 透传消息内容
	Data []byte
}

func (entity *T808_0x0900) MsgID() MsgID {
	return MsgT808_0x0900
}

func (entity *T808_0x0900) Encode() ([]byte, error) {
	writer := NewWriter()
	writer.WriteByte(entity.Type)
	writer.Write(entity.Data)
	return writer.Bytes(), nil
}
func (entity *T808_0x0900) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, ErrInvalidBody
	}
	entity.Type, entity.Data = data[0], data[1:]
	return len(data), nil
}
