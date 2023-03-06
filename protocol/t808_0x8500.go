package protocol

// 车门控制
type T808_0x8500 struct {
	// 控制标志
	Flag T808_0x8500_Flag
}

// 车门控制标志
type T808_0x8500_Flag byte

// 设置上锁
func (flag *T808_0x8500_Flag) SetLock(b bool) {
	SetBitByte((*byte)(flag), 0, b)
}

func (entity *T808_0x8500) MsgID() MsgID {
	return MsgT808_0x8500
}

func (entity *T808_0x8500) Encode() ([]byte, error) {
	return []byte{byte(entity.Flag)}, nil
}

func (entity *T808_0x8500) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, ErrInvalidBody
	}
	entity.Flag = T808_0x8500_Flag(data[0])
	return 1, nil
}
