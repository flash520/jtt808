package protocol

// 终端应答结果
type T808_0x8100_Result byte

const (
	// 成功
	T808_0x8100_ResultSuccess = 0
	// 车辆已被注册
	T808_0x8100_ResultCarRegistered = 1
	// 数据库中无该车辆
	T808_0x8100_ResultCarNotFound = 2
	// 终端已被注册
	T808_0x8100_ResultTerminalRegistered = 3
	// 数据库中无该终端
	T808_0x8100_ResultTerminalNotFound = 4
)

// 终端应答
type T808_0x8100 struct {
	// 应答流水号
	MsgSerialNo uint16
	// 结果
	Result T808_0x8100_Result
	// 鉴权码
	AuthKey string
}

func (entity *T808_0x8100) MsgID() MsgID {
	return MsgT808_0x8100
}

func (entity *T808_0x8100) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入消息流水号
	writer.WriteUint16(entity.MsgSerialNo)

	// 写入响应结果
	writer.WriteByte(byte(entity.Result))

	// 写入鉴权码
	if len(entity.AuthKey) > 0 {
		if err := writer.WritString(entity.AuthKey); err != nil {
			return nil, err
		}
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x8100) Decode(data []byte) (int, error) {
	if len(data) < 3 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取流水号
	msgSerialNo, err := reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取响应结果
	temp, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取鉴权码
	if reader.Len() > 0 {
		entity.AuthKey, err = reader.ReadString()
		if err != nil {
			return 0, err
		}
	}

	entity.Result = T808_0x8100_Result(temp)
	entity.MsgSerialNo = msgSerialNo
	return len(data) - reader.Len(), nil
}
