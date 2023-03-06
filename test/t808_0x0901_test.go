package test

import (
	"reflect"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

func TestT808_0x0901_EncodeDecode(t *testing.T) {
	message := protocol.T808_0x0901{
		Uncompressed: make([]byte, 768),
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x0901
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
	assert.True(t, reflect.DeepEqual(message, message2))
}
