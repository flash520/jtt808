package protocol

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

// 信息服务
type T808_0x8304 struct {
	// 信息类型
	Type byte
	// 信息内容
	Content string
}

func (entity *T808_0x8304) MsgID() MsgID {
	return MsgT808_0x8304
}

func (entity *T808_0x8304) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入信息类型
	writer.WriteByte(entity.Type)

	// 写入信息长度
	reader := bytes.NewReader([]byte(entity.Content))
	data, err := ioutil.ReadAll(
		transform.NewReader(reader, simplifiedchinese.GB18030.NewEncoder()))
	if err != nil {
		return nil, err
	}
	writer.WriteUint16(uint16(len(data)))

	// 写入信息内容
	writer.Write(data)
	return writer.Bytes(), nil
}

func (entity *T808_0x8304) Decode(data []byte) (int, error) {
	if len(data) < 3 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取信息类型
	var err error
	entity.Type, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取信息长度
	size, err := reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取信息内容
	entity.Content, err = reader.ReadString(int(size))
	if err != nil {
		return 0, err
	}
	return len(data) - reader.Len(), nil
}
