package protocol

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

// 位置点播菜单设置
type T808_0x8303 struct {
	// 设置类型
	Type T808_0x8303_Type
	// 信息项列表
	Items []T808_0x8303_Item
}

// 设置类型
type T808_0x8303_Type byte

var (
	// 删除所有
	T808_0x8303_TypeRemoveAll T808_0x8303_Type = 0
	// 更新
	T808_0x8303_TypeUpdate T808_0x8303_Type = 1
	// 追加
	T808_0x8303_TypeAppend T808_0x8303_Type = 2
	// 修改
	T808_0x8303_TypeModify T808_0x8303_Type = 3
)

// 信息点播信息项
type T808_0x8303_Item struct {
	// 信息名称
	Name string
}

func (entity *T808_0x8303) MsgID() MsgID {
	return MsgT808_0x8303
}

func (entity *T808_0x8303) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入设置类型
	writer.WriteByte(byte(entity.Type))

	// 写入信息项总数
	writer.WriteByte(byte(len(entity.Items)))

	// 写入信息项列表
	for _, item := range entity.Items {
		reader := bytes.NewReader([]byte(item.Name))
		data, err := ioutil.ReadAll(
			transform.NewReader(reader, simplifiedchinese.GB18030.NewEncoder()))
		if err != nil {
			return nil, err
		}

		// 写入名称长度
		writer.WriteUint16(uint16(len(data)))

		// 写入信息名称
		writer.Write(data)
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x8303) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取设置类型
	typ, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	entity.Type = T808_0x8303_Type(typ)

	// 读取信息项总数
	count, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取信息项列表
	for i := 0; i < int(count); i++ {
		// 读取名称长度
		size, err := reader.ReadUint16()
		if err != nil {
			return 0, err
		}

		// 读取信息名称
		name, err := reader.ReadString(int(size))
		if err != nil {
			return 0, err
		}
		entity.Items = append(entity.Items, T808_0x8303_Item{Name: name})
	}
	return len(data) - reader.Len(), nil
}
