package protocol

// 终端控制
type T808_0x8105 struct {
	// 命令字
	Cmd T808_0x8105_Command
	// 命令参数
	Data string
}

// 终端控制命令
type T808_0x8105_Command byte

var (
	// 升级
	T808_0x8105_CommandUpgrade T808_0x8105_Command = 1
	// 设置服务器
	T808_0x8105_CommandHost T808_0x8105_Command = 2
	// 关机
	T808_0x8105_CommandShutdown T808_0x8105_Command = 3
	// 重置
	T808_0x8105_CommandReset T808_0x8105_Command = 4
	// 恢复出厂设置
	T808_0x8105_CommandFactoryReset T808_0x8105_Command = 5
	// 关闭无线网络
	T808_0x8105_CommandCloseNetwork T808_0x8105_Command = 6
	// 关闭所有网络
	T808_0x8105_CommandCloseAllNetwork T808_0x8105_Command = 7
)

func (entity *T808_0x8105) MsgID() MsgID {
	return MsgT808_0x8105
}

func (entity *T808_0x8105) Encode() ([]byte, error) {
	writer := NewWriter()
	writer.WriteByte(byte(entity.Cmd))
	if len(entity.Data) > 0 {
		if err := writer.WritString(entity.Data); err != nil {
			return nil, err
		}
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x8105) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, ErrInvalidBody
	}

	reader := NewReader(data[1:])
	if reader.Len() > 0 {
		data, err := reader.ReadString()
		if err != nil {
			return 0, err
		}
		entity.Data = data
	}

	entity.Cmd = T808_0x8105_Command(data[0])
	return len(data) - reader.Len() - 1, nil
}
