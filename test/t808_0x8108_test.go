package test

import (
	"reflect"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

func TestT808_0x8108_EncodeDecode(t *testing.T) {
	message := protocol.T808_0x8108{
		Type:          1,
		ManufactureID: "abc12",
		Version:       "10.1.1",
		Size:          102400,
		Data:          []byte{1, 2, 3, 4},
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x8108
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
	assert.True(t, reflect.DeepEqual(message, message2))
}
