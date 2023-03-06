package test

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

func TestT808_0x8801_EncodeDecode(t *testing.T) {
	message := protocol.T808_0x8801{
		ChannelID:    28,
		Cmd:          456,
		Duration:     360,
		Save:         protocol.T808_0x8801_SaveFlagLocal,
		Resolution:   protocol.T808_0x8801_Resolution176x144,
		Quality:      50,
		Lighting:     50,
		Contrast:     50,
		Saturability: 50,
		Chroma:       50,
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x8801
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
}
