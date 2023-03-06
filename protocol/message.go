package protocol

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	"reflect"
	
	log "github.com/sirupsen/logrus"
	
	"github.com/flash520/jtt808/errors"
)

// 消息包
type Message struct {
	Header Header
	Body   Entity
}

// 协议编码
func (message *Message) Encode(key ...*rsa.PublicKey) ([]byte, error) {
	// 编码消息体
	count := 0
	var err error
	var body []byte
	checkSum := byte(0x00)
	if message.Body != nil && !reflect.ValueOf(message.Body).IsNil() {
		body, err = message.Body.Encode()
		if err != nil {
			return nil, err
		}
		
		if len(key) > 0 && key[0] != nil {
			message.Header.Property.enableEncrypt()
			body, err = EncryptOAEP(sha1.New(), key[0], body, nil)
			if err != nil {
				log.WithFields(log.Fields{
					"id":     fmt.Sprintf("0x%x", message.Header.MsgID),
					"reason": err,
				}).Warn("[JT/T 808] encrypt body failed")
				return nil, err
			}
		}
	}
	checkSum, count = message.computeChecksum(body, checkSum, count)
	
	// 编码消息头
	message.Header.MsgID = message.Body.MsgID()
	err = message.Header.Property.SetBodySize(uint16(len(body)))
	if err != nil {
		return nil, err
	}
	header, err := message.Header.Encode()
	if err != nil {
		return nil, err
	}
	checkSum, count = message.computeChecksum(header, checkSum, count)
	
	// 二进制转义
	buffer := bytes.NewBuffer(nil)
	buffer.Grow(count + 2)
	buffer.WriteByte(PrefixID)
	message.write(buffer, header).write(buffer, body).write(buffer, []byte{checkSum})
	buffer.WriteByte(PrefixID)
	return buffer.Bytes(), nil
}

// 协议解码
func (message *Message) Decode(data []byte, key ...*rsa.PrivateKey) error {
	// 检验标志位
	if len(data) < 2 || data[0] != PrefixID || data[len(data)-1] != PrefixID {
		return errors.ErrInvalidMessage
	}
	data = data[1 : len(data)-1]
	if len(data) == 0 {
		return errors.ErrInvalidMessage
	}
	
	// 获取校验和
	sum := data[len(data)-1]
	if data[len(data)-2] != EscapeByte {
		data = data[:len(data)-1]
	} else {
		if (data[len(data)-1]) == EscapeByteSufix1 {
			sum = EscapeByte
		} else if data[len(data)-1] == EscapeByteSufix2 {
			sum = PrefixID
		} else {
			return errors.ErrInvalidMessage
		}
		data = data[:len(data)-2]
	}
	
	// 二进制转义
	checkSum := byte(0x00)
	buffer := make([]byte, 0, len(data))
	for i := 0; i < len(data); {
		b := data[i]
		if b != EscapeByte {
			checkSum = checkSum ^ b
			buffer = append(buffer, b)
			i++
			continue
		}
		
		if i+1 >= len(data) {
			return errors.ErrInvalidMessage
		}
		
		b = data[i+1]
		if b == EscapeByteSufix1 {
			checkSum = checkSum ^ EscapeByte
			buffer = append(buffer, EscapeByte)
		} else if b == EscapeByteSufix2 {
			checkSum = checkSum ^ PrefixID
			buffer = append(buffer, PrefixID)
		} else {
			return errors.ErrInvalidMessage
		}
		i += 2
	}
	
	// 检查校验和
	if len(buffer) == 0 || checkSum != sum {
		return errors.ErrInvalidCheckSum
	}
	
	// 解码消息头
	if len(buffer) < MessageHeaderSize {
		return errors.ErrInvalidHeader
	}
	var header Header
	err := header.Decode(buffer)
	if err != nil {
		return err
	}
	if !header.Property.IsEnablePacket() {
		buffer = buffer[MessageHeaderSize:]
	} else {
		buffer = buffer[MessageHeaderSize+4:]
	}
	
	// 解码消息体
	if uint16(len(buffer)) != header.Property.GetBodySize() {
		log.WithFields(log.Fields{
			"id":     fmt.Sprintf("0x%x", header.MsgID),
			"expect": header.Property.GetBodySize(),
			"actual": len(buffer),
		}).Warn("[JT/T 808] body length mismatch")
	} else {
		if header.Property.IsEnableEncrypt() {
			if len(key) == 0 || key[0] == nil {
				log.WithFields(log.Fields{
					"id":     fmt.Sprintf("0x%x", header.MsgID),
					"reason": "private key not found",
				}).Warn("[JT/T 808] decrypt body failed")
				return errors.ErrDecryptMessageFailed
			}
			
			buffer, err = DecryptOAEP(sha1.New(), key[0], buffer, nil)
			if err != nil {
				log.WithFields(log.Fields{
					"id":     fmt.Sprintf("0x%x", header.MsgID),
					"reason": err,
				}).Warn("[JT/T 808] decrypt body failed")
				return errors.ErrDecryptMessageFailed
			}
		}
		
		entity, _, err := message.decode(uint16(header.MsgID), buffer)
		if err == nil {
			message.Body = entity
		} else {
			log.WithFields(log.Fields{
				"id":     fmt.Sprintf("0x%x", header.MsgID),
				"reason": err,
			}).Warn("[JT/T 808] failed to decode message")
		}
	}
	message.Header = header
	return nil
}
func (message *Message) decode(typ uint16, data []byte) (Entity, int, error) {
	creator, ok := entityMapper[typ]
	if !ok {
		return nil, 0, errors.ErrTypeNotRegistered
	}
	
	entity := creator()
	entityPacket, ok := interface{}(entity).(EntityPacket)
	if !ok {
		count, err := entity.Decode(data)
		if err != nil {
			return nil, 0, err
		}
		return entity, count, nil
	}
	
	err := entityPacket.DecodePacket(data)
	if err != nil {
		return nil, 0, err
	}
	return entityPacket, len(data), nil
}

// 写入二进制数据
func (message *Message) write(buffer *bytes.Buffer, data []byte) *Message {
	for _, b := range data {
		if b == PrefixID {
			buffer.WriteByte(EscapeByte)
			buffer.WriteByte(EscapeByteSufix2)
		} else if b == EscapeByte {
			buffer.WriteByte(EscapeByte)
			buffer.WriteByte(EscapeByteSufix1)
		} else {
			buffer.WriteByte(b)
		}
	}
	return message
}

// 校验和累加计算
func (message *Message) computeChecksum(data []byte, checkSum byte, count int) (byte, int) {
	for _, b := range data {
		checkSum = checkSum ^ b
		if b != PrefixID && b != EscapeByte {
			count++
		} else {
			count += 2
		}
	}
	return checkSum, count
}
