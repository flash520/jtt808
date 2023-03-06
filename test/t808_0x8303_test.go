package test

import (
	"reflect"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

func TestT808_0x8303_EncodeDecode(t *testing.T) {
	message := protocol.T808_0x8303{
		Type: protocol.T808_0x8303_TypeModify,
		Items: []protocol.T808_0x8303_Item{
			{Name: "点播项1"},
			{Name: "点播项2"},
			{Name: "点播项3"},
		},
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x8303
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
	assert.True(t, reflect.DeepEqual(message, message2))
}
