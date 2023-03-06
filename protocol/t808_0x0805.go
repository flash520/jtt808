package protocol

// 摄像头立即拍摄命令应答
type T808_0x0805 struct {
	// 应答流水号
	ReplyMsgSerialNo uint16
	// 结果
	// 0：成功； 1：失败； 2：通道不支持。
	Result T808_0x0805_Result
	// 多媒体 ID 个数
	MediaIDs []uint32
}

// 处理结果
type T808_0x0805_Result byte

var (
	// 成功
	T808_0x0805_ResultSuccess T808_0x0805_Result = 0
	// 失败
	T808_0x0805_ResultFail T808_0x0805_Result = 1
	// 不支持
	T808_0x0805_ResultNotSupported T808_0x0805_Result = 2
)

func (entity *T808_0x0805) MsgID() MsgID {
	return MsgT808_0x0805
}

func (entity *T808_0x0805) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入应答流水号
	writer.WriteUint16(entity.ReplyMsgSerialNo)

	// 写入处理结果
	writer.WriteByte(byte(entity.Result))

	// 写入媒体ID数量
	writer.WriteUint16(uint16(len(entity.MediaIDs)))

	// 写入媒体ID列表
	for _, id := range entity.MediaIDs {
		writer.WriteUint32(id)
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x0805) Decode(data []byte) (int, error) {
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

	// 读取处理结果
	result, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	entity.Result = T808_0x0805_Result(result)

	// 读取媒体ID数量
	count, err := reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取媒体ID列表
	entity.MediaIDs = make([]uint32, 0, count)
	for i := 0; i < int(count); i++ {
		id, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		entity.MediaIDs = append(entity.MediaIDs, id)
	}
	return len(data) - reader.Len(), nil
}
