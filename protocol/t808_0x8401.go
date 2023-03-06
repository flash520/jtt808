package protocol

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

// 设置电话本
type T808_0x8401 struct {
	// 设置类型
	Type T808_0x8401_Type
	// 联系人项
	Contacts []T808_0x8401_Contact
}

// 设置类型
type T808_0x8401_Type byte

var (
	// 删除所有
	T808_0x8401_TypeRemoveAll T808_0x8401_Type = 0
	// 更新
	T808_0x8401_TypeUpdate T808_0x8401_Type = 1
	// 追加
	T808_0x8401_TypeAppend T808_0x8401_Type = 2
	// 修改
	T808_0x8401_TypeModify T808_0x8401_Type = 3
)

// 电话本联系人
type T808_0x8401_Contact struct {
	// 标志
	Flag T808_0x8401_ContactFlag
	// 电话号码
	Number string
	// 联系人
	Contact string
}

// 联系人标志
type T808_0x8401_ContactFlag byte

var (
	// 呼入
	T808_0x8401_ContactFlagIncomingCall T808_0x8401_ContactFlag = 1
	// 呼出
	T808_0x8401_ContactFlagOutgoingCall T808_0x8401_ContactFlag = 2
	// 呼入和呼出
	T808_0x8401_ContactFlagBoth T808_0x8401_ContactFlag = 3
)

func (entity *T808_0x8401) MsgID() MsgID {
	return MsgT808_0x8401
}

func (entity *T808_0x8401) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入设置类型
	writer.WriteByte(byte(entity.Type))

	// 写入联系人总数
	writer.WriteByte(byte(len(entity.Contacts)))

	// 写入联系人列表
	for _, contact := range entity.Contacts {
		// 写入标志
		writer.WriteByte(byte(contact.Flag))

		// 写入号码长度
		reader := bytes.NewReader([]byte(contact.Number))
		number, err := ioutil.ReadAll(
			transform.NewReader(reader, simplifiedchinese.GB18030.NewEncoder()))
		if err != nil {
			return nil, err
		}
		writer.WriteByte(byte(len(number)))

		// 写入电话号码
		writer.Write(number)

		// 写入联系人长度
		reader = bytes.NewReader([]byte(contact.Contact))
		contact, err := ioutil.ReadAll(
			transform.NewReader(reader, simplifiedchinese.GB18030.NewEncoder()))
		if err != nil {
			return nil, err
		}
		writer.WriteByte(byte(len(contact)))

		// 写入联系人
		writer.Write(contact)
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x8401) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取设置类型
	typ, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	entity.Type = T808_0x8401_Type(typ)

	// 读取联系人总数
	count, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取联系人列表
	for i := 0; i < int(count); i++ {
		var contact T808_0x8401_Contact

		// 读取标志
		flag, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}
		contact.Flag = T808_0x8401_ContactFlag(flag)

		// 读取号码长度
		size, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}

		// 读取电话号码
		contact.Number, err = reader.ReadString(int(size))
		if err != nil {
			return 0, err
		}

		// 读取联系人长度
		size, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}

		// 读取联系人
		contact.Contact, err = reader.ReadString(int(size))
		if err != nil {
			return 0, err
		}
		entity.Contacts = append(entity.Contacts, contact)
	}
	return len(data) - reader.Len(), nil
}
