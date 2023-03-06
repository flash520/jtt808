package protocol

// 终端升级结果通知
type T808_0x0108 struct {
	// 升级类型
	// 0：终端， 12：道路运输证 IC 卡读卡器， 52：北斗卫星定位模块
	Type byte
	// 升级结果
	// 0：成功， 1：失败， 2：取消
	Result byte
}

func (entity *T808_0x0108) MsgID() MsgID {
	return MsgT808_0x0108
}

func (entity *T808_0x0108) Encode() ([]byte, error) {
	return []byte{entity.Type, entity.Result}, nil
}

func (entity *T808_0x0108) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, ErrInvalidBody
	}

	entity.Type, entity.Result = data[0], data[1]
	return 2, nil
}
