package test

import (
	"reflect"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

func TestT808_0x8700_EncodeDecode(t *testing.T) {
	message := protocol.T808_0x8700{
		Cmd:  26,
		Data: []byte{1, 3, 5, 7, 9},
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x8700
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
	assert.True(t, reflect.DeepEqual(message, message2))
}
