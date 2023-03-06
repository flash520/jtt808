package protocol

// 录音开始命令
type T808_0x8804 struct {
	// 录音命令
	// 0：停止录音； 0x01：开始录音
	Cmd byte
	// 录音时间
	// 单位为秒（ s），0表示一直录音
	Duration uint16
	// 保存标志
	Save T808_0x8801_SaveFlag
	// 音频采样率
	AudioSampleRate T808_0x8804_AudioSampleRate
}

// 音频采样率
type T808_0x8804_AudioSampleRate byte

var (
	T808_0x8804_AudioSampleRate8k  T808_0x8804_AudioSampleRate = 0
	T808_0x8804_AudioSampleRate11k T808_0x8804_AudioSampleRate = 1
	T808_0x8804_AudioSampleRate23k T808_0x8804_AudioSampleRate = 2
	T808_0x8804_AudioSampleRate32k T808_0x8804_AudioSampleRate = 3
)

func (entity *T808_0x8804) MsgID() MsgID {
	return MsgT808_0x8804
}

func (entity *T808_0x8804) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入录音命令
	writer.WriteByte(entity.Cmd)

	// 写入录音时间
	writer.WriteUint16(entity.Duration)

	// 写入保存标志
	writer.WriteByte(byte(entity.Save))

	// 写入音频采样率
	writer.WriteByte(byte(entity.AudioSampleRate))
	return writer.Bytes(), nil
}

func (entity *T808_0x8804) Decode(data []byte) (int, error) {
	if len(data) < 5 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取录音命令
	var err error
	entity.Cmd, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取录音时间
	entity.Duration, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取保存标志
	saveFlag, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	entity.Save = T808_0x8801_SaveFlag(saveFlag)

	// 读取音频采样率
	audioSampleRate, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	entity.AudioSampleRate = T808_0x8804_AudioSampleRate(audioSampleRate)
	return len(data) - reader.Len(), nil
}
