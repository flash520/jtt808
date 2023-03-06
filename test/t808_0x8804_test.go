package test

import (
	"reflect"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

func TestT808_0x8804_EncodeDecode(t *testing.T) {
	message := protocol.T808_0x8804{
		Cmd:             67,
		Duration:        3500,
		Save:            protocol.T808_0x8801_SaveFlagRemote,
		AudioSampleRate: protocol.T808_0x8804_AudioSampleRate8k,
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x8804
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
	assert.True(t, reflect.DeepEqual(message, message2))
}
