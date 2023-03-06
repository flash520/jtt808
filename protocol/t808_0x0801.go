package protocol

import (
	"bytes"
	"io"
	"io/ioutil"
)

// 多媒体数据上传
type T808_0x0801 struct {
	// 多媒体 ID
	MediaID uint32
	// 多媒体类型
	Type T808_0x0800_MediaType
	// 多媒体格式编码
	Coding T808_0x0800_MediaCoding
	// 事件项编码
	Event byte
	// 通道 ID
	ChannelID byte
	// 位置信息汇报
	Location T808_0x0200
	// 多媒体数据包
	Packet io.Reader
}

func (entity *T808_0x0801) MsgID() MsgID {
	return MsgT808_0x0801
}

func (entity *T808_0x0801) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入媒体ID
	writer.WriteUint32(entity.MediaID)

	// 写入媒体类型
	writer.WriteByte(byte(entity.Type))

	// 写入媒体编码
	writer.WriteByte(byte(entity.Coding))

	// 写入事件类型
	writer.WriteByte(entity.Event)

	// 写入通道ID
	writer.WriteByte(entity.ChannelID)

	// 写入定位信息
	entity.Location.Extras = nil
	data, err := entity.Location.Encode()
	if err != nil {
		return nil, err
	}
	writer.Write(data)

	// 写入数据包
	if entity.Packet != nil {
		data, err = ioutil.ReadAll(entity.Packet)
		if err != nil {
			return nil, err
		}
		writer.Write(data)
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x0801) Decode(data []byte) (int, error) {
	if len(data) < 36 {
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

	// 读取事件类型
	entity.Event, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取通道ID
	entity.ChannelID, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取定位信息
	buf, err := reader.Read(28)
	if err != nil {
		return 0, err
	}
	if _, err = entity.Location.Decode(buf); err != nil {
		return 0, err
	}
	return len(data) - reader.Len(), nil
}

func (entity *T808_0x0801) GetTag() uint32 {
	return entity.MediaID
}

func (entity *T808_0x0801) GetReader() io.Reader {
	return entity.Packet
}

func (entity *T808_0x0801) SetReader(reader io.Reader) {
	entity.Packet = reader
}

func (entity *T808_0x0801) DecodePacket(data []byte) error {
	n, err := entity.Decode(data)
	if err != nil {
		return err
	}
	entity.Packet = bytes.NewReader(data[n:])
	return nil
}
