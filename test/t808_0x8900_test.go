package test

import (
	"reflect"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

func TestT808_0x8900_EncodeDecode(t *testing.T) {
	message := protocol.T808_0x8900{
		Type: 0x01,
		Data: []byte{0x34, 0x45, 0x23, 0x45},
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x8900
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
	assert.True(t, reflect.DeepEqual(message, message2))
}
