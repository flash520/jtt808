package test

import (
	"reflect"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

func TestT808_0x0100_EncodeDecode(t *testing.T) {
	message := protocol.T808_0x0100{
		ProvinceID:    023,
		CityID:        023,
		ManufactureID: "test",
		Model:         "t900",
		TerminalID:    "abc123",
		PlateColor:    45,
		LicenseNo:     "asdfg阿伯吃得",
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x0100
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
	assert.True(t, reflect.DeepEqual(message, message2))
}
