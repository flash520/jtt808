package test

import (
	"reflect"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

func TestT808_0x0001_EncodeDecode(t *testing.T) {
	message := protocol.T808_0x0001{
		ReplyMsgSerialNo: 1024,
		ReplyMsgID:       123,
		Result:           protocol.T808_0x8001ResultSuccess,
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x0001
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
	assert.True(t, reflect.DeepEqual(message, message2))
}
