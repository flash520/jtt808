package protocol

// 多媒体事件信息上传
type T808_0x0800 struct {
	// 多媒体数据 ID
	MediaID uint32
	// 多媒体类型
	Type T808_0x0800_MediaType
	// 多媒体格式编码
	Coding T808_0x0800_MediaCoding
	// 事件项编码
	Event byte
	// 通道 ID
	ChannelID byte
}

// 多媒体类型
type T808_0x0800_MediaType byte

var (
	T808_0x0800_MediaTypeImage T808_0x0800_MediaType = 0
	T808_0x0800_MediaTypeAudio T808_0x0800_MediaType = 1
	T808_0x0800_MediaTypeVideo T808_0x0800_MediaType = 2
)

// 多媒体编码
type T808_0x0800_MediaCoding byte

var (
	T808_0x0800_MediaCodingJPEG T808_0x0800_MediaCoding = 0
	T808_0x0800_MediaCodingTIF  T808_0x0800_MediaCoding = 1
	T808_0x0800_MediaCodingMP3  T808_0x0800_MediaCoding = 2
	T808_0x0800_MediaCodingWAV  T808_0x0800_MediaCoding = 3
	T808_0x0800_MediaCodingWMV  T808_0x0800_MediaCoding = 4
)

func (entity *T808_0x0800) MsgID() MsgID {
	return MsgT808_0x0800
}

func (entity *T808_0x0800) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入媒体ID
	writer.WriteUint32(entity.MediaID)

	// 写入媒体类型
	writer.WriteByte(byte(entity.Type))

	// 写入媒体编码
	writer.WriteByte(byte(entity.Coding))

	// 写入媒体事件
	writer.WriteByte(entity.Event)

	// 写入通道ID
	writer.WriteByte(entity.ChannelID)
	return writer.Bytes(), nil
}

func (entity *T808_0x0800) Decode(data []byte) (int, error) {
	if len(data) < 8 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取媒体ID
	var err error
	entity.MediaID, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	// 读取媒体类型
	mediaType, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	entity.Type = T808_0x0800_MediaType(mediaType)

	// 读取媒体编码
	coding, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	entity.Coding = T808_0x0800_MediaCoding(coding)

	// 读取媒体事件
	entity.Event, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取通道ID
	entity.ChannelID, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}
	return len(data) - reader.Len(), nil
}
