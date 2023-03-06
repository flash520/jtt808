package test

import (
	"reflect"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

func TestT808_0x0800_EncodeDecode(t *testing.T) {
	message := protocol.T808_0x0800{
		MediaID:   3456,
		Type:      protocol.T808_0x0800_MediaTypeVideo,
		Coding:    protocol.T808_0x0800_MediaCodingMP3,
		Event:     23,
		ChannelID: 17,
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x0800
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
	assert.True(t, reflect.DeepEqual(message, message2))
}
