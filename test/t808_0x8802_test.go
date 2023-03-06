package test

import (
	"reflect"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

func TestT808_0x8802_EncodeDecode(t *testing.T) {
	message := protocol.T808_0x8802{
		Type:      protocol.T808_0x0800_MediaTypeAudio,
		ChannelID: 56,
		Event:     87,
		StartTime: time.Unix(time.Now().Unix(), 0),
		EndTime:   time.Unix(time.Now().Unix(), 0),
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x8802
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
	assert.True(t, reflect.DeepEqual(message, message2))
}
