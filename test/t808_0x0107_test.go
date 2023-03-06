package test

import (
	"reflect"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

func TestT808_0x0107_EncodeDecode(t *testing.T) {
	message := protocol.T808_0x0107{
		TerminalType:    123,
		ManufactureID:   "wer",
		Model:           "1231231",
		TerminalID:      "t800",
		Sim:             "8618202374929",
		HardwareVersion: "v0.1.1",
		SoftwareVersion: "v0.1.2",
		GNSSProperty:    23,
		COMMProperty:    34,
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x0107
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
	assert.True(t, reflect.DeepEqual(message, message2))
}
