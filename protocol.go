package jtt808

import (
	"bytes"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"io"
	
	"github.com/flash520/link"
	log "github.com/sirupsen/logrus"
	
	"github.com/flash520/jtt808/errors"
	"github.com/flash520/jtt808/protocol"
)

type Protocol struct {
	privateKey *rsa.PrivateKey
}

// NewCodec 创建编解码器
func (p Protocol) NewCodec(rw io.ReadWriter) (link.Codec, error) {
	codec := &ProtocolCodec{
		w:               rw,
		r:               rw,
		privateKey:      p.privateKey,
		bufferReceiving: bytes.NewBuffer(nil),
	}
	codec.closer, _ = rw.(io.Closer)
	return codec, nil
}

// ProtocolCodec 编解码器
type ProtocolCodec struct {
	w               io.Writer
	r               io.Reader
	closer          io.Closer
	publicKey       *rsa.PublicKey
	privateKey      *rsa.PrivateKey
	bufferReceiving *bytes.Buffer
}

// GetPublicKey 获取RSA公钥
func (codec *ProtocolCodec) GetPublicKey() *rsa.PublicKey {
	return codec.publicKey
}

// SetPublicKey 设置RSA公钥
func (codec *ProtocolCodec) SetPublicKey(publicKey *rsa.PublicKey) {
	codec.publicKey = publicKey
}

// Close 关闭读写
func (codec *ProtocolCodec) Close() error {
	if codec.closer != nil {
		return codec.closer.Close()
	}
	return nil
}

// Send 发送消息
func (codec *ProtocolCodec) Send(msg interface{}) error {
	message, ok := msg.(protocol.Message)
	if !ok {
		log.WithFields(log.Fields{
			"reason": errors.ErrInvalidMessage,
		}).Error("[JT/T 808] failed to write message")
		return errors.ErrInvalidMessage
	}
	
	var err error
	var data []byte
	if codec.publicKey == nil || !message.Header.Property.IsEnableEncrypt() {
		data, err = message.Encode()
	} else {
		data, err = message.Encode(codec.publicKey)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"id":     fmt.Sprintf("0x%x", message.Header.MsgID),
			"reason": err,
		}).Error("[JT/T 808] failed to write message")
		return err
	}
	
	count, err := codec.w.Write(data)
	if err != nil {
		log.WithFields(log.Fields{
			"id":     fmt.Sprintf("0x%x", message.Header.MsgID),
			"reason": err,
		}).Error("[JT/T 808] failed to write message")
		return err
	}
	
	log.WithFields(log.Fields{
		"id":    fmt.Sprintf("0x%x", message.Header.MsgID),
		"bytes": count,
	}).Debug("[JT/T 808] write message success")
	return nil
}

// Receive 接收消息
func (codec *ProtocolCodec) Receive() (interface{}, error) {
	message, ok, err := codec.readFromBuffer()
	if ok {
		return message, nil
	}
	if err != nil {
		return nil, err
	}
	
	var buffer [128]byte
	for {
		count, err := io.ReadAtLeast(codec.r, buffer[:], 1)
		if err != nil {
			return nil, err
		}
		codec.bufferReceiving.Write(buffer[:count])
		
		if codec.bufferReceiving.Len() == 0 {
			continue
		}
		if codec.bufferReceiving.Len() > 0xffff {
			return nil, errors.ErrBodyTooLong
		}
		
		message, ok, err := codec.readFromBuffer()
		if ok {
			return message, nil
		}
		if err != nil {
			return nil, err
		}
	}
}

// readFromBuffer 从缓冲区读取
func (codec *ProtocolCodec) readFromBuffer() (protocol.Message, bool, error) {
	if codec.bufferReceiving.Len() == 0 {
		return protocol.Message{}, false, nil
	}
	
	data := codec.bufferReceiving.Bytes()
	if data[0] != protocol.PrefixID {
		i := 0
		for ; i < len(data); i++ {
			if data[i] == protocol.PrefixID {
				break
			}
		}
		codec.bufferReceiving.Next(i)
		log.WithFields(log.Fields{
			"data":   hex.EncodeToString(data),
			"reason": errors.ErrNotFoundPrefixID,
		}).Error("[JT/T 808] failed to receive message")
		return protocol.Message{}, false, errors.ErrNotFoundPrefixID
	}
	
	end := 1
	for ; end < len(data); end++ {
		if data[end] == protocol.PrefixID {
			break
		}
	}
	if end == len(data) {
		return protocol.Message{}, false, nil
	}
	
	var message protocol.Message
	if err := message.Decode(data[:end+1], codec.privateKey); err != nil {
		codec.bufferReceiving.Next(end + 1)
		log.WithFields(log.Fields{
			"data":   fmt.Sprintf("0x%x", hex.EncodeToString(data[:end+1])),
			"reason": err,
		}).Error("[JT/T 808] failed to receive message")
		return protocol.Message{}, false, err
	}
	codec.bufferReceiving.Next(end + 1)
	
	log.WithFields(log.Fields{
		"id": fmt.Sprintf("0x%x", message.Header.MsgID),
	}).Debug("[JT/T 808] new message received")
	log.WithFields(log.Fields{
		"data": fmt.Sprintf("0x%x", hex.EncodeToString(data[:end+1])),
	}).Trace("[JT/T 808] message hex string")
	return message, true, nil
}
