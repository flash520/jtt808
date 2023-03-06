package protocol

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

// 数据压缩上报
type T808_0x0901 struct {
	// 未压缩消息
	Uncompressed []byte
}

func (entity *T808_0x0901) MsgID() MsgID {
	return MsgT808_0x0901
}

func (entity *T808_0x0901) Encode() ([]byte, error) {
	// 压缩消息
	buffer := bytes.NewBuffer(nil)
	gzipWriter := gzip.NewWriter(buffer)
	gzipWriter.Write(entity.Uncompressed)
	gzipWriter.Close()

	// 写入消息长度
	writer := NewWriter()
	writer.WriteUint32(uint32(buffer.Len()))

	// 写入压缩消息
	fmt.Println(len(entity.Uncompressed))
	writer.Write(buffer.Bytes())
	return writer.Bytes(), nil
}

func (entity *T808_0x0901) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取消息长度
	size, err := reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	// 读取压缩消息
	compressed, err := reader.Read(int(size))
	if err != nil {
		return 0, err
	}

	// 解压缩消息
	buffer := bytes.NewBuffer(nil)
	gzipReader, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return 0, err
	}
	defer gzipReader.Close()
	var temp [256]byte
	for {
		n, err := gzipReader.Read(temp[:])
		if err != nil && err != io.EOF {
			break
		}

		if n > 0 {
			buffer.Write(temp[:n])
		}

		if err == io.EOF {
			break
		}
	}
	entity.Uncompressed = buffer.Bytes()
	return len(data) - reader.Len(), nil
}
