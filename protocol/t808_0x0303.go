package protocol

// 信息点播/取消
type T808_0x0303 struct {
	// 信息类型
	Type byte
	// 点播/ 取消标志
	Flag T808_0x0303_Flag
}

// 信息点播标志
type T808_0x0303_Flag byte

var (
	// 点播
	T808_0x0303_FlagRequest T808_0x0303_Flag = 0
	// 取消
	T808_0x0303_FlagCancel T808_0x0303_Flag = 1
)

func (entity *T808_0x0303) MsgID() MsgID {
	return MsgT808_0x0303
}

func (entity *T808_0x0303) Encode() ([]byte, error) {
	return []byte{entity.Type, byte(entity.Flag)}, nil
}

func (entity *T808_0x0303) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, ErrInvalidBody
	}
	entity.Type, entity.Flag = data[0], T808_0x0303_Flag(data[1])
	return 1, nil
}
