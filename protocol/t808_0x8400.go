package protocol

// 电话回拨
type T808_0x8400 struct {
	// 标志
	Flag T808_0x8400_Flag
	// 电话号码
	Number string
}

// 回拨标志
type T808_0x8400_Flag byte

var (
	// 正常
	T808_0x8400_FlagNormal T808_0x8400_Flag = 0
	// 监听
	T808_0x8400_FlagMonitor T808_0x8400_Flag = 1
)

func (entity *T808_0x8400) MsgID() MsgID {
	return MsgT808_0x8400
}

func (entity *T808_0x8400) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入回拨标志
	writer.WriteByte(byte(entity.Flag))

	// 写入电话号码
	writer.WritString(entity.Number)
	return writer.Bytes(), nil
}

func (entity *T808_0x8400) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取回拨标志
	flag, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	entity.Flag = T808_0x8400_Flag(flag)

	// 读取电话号码
	entity.Number, err = reader.ReadString()
	if err != nil {
		return 0, err
	}
	return len(data) - reader.Len(), nil
}
