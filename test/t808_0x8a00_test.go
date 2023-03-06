package test

import (
	"reflect"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

func TestT808_0x8A00_EncodeDecode(t *testing.T) {
	privateKey, err := GetTestPrivateKey()
	if err != nil {
		t.Error(err)
	}
	
	message := protocol.T808_0x8A00{
		PublicKey: &privateKey.PublicKey,
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x8A00
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
	assert.True(t, reflect.DeepEqual(message, message2))
}
