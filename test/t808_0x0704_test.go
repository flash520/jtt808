package test

import (
	"reflect"
	"testing"
	"time"
	
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

func TestT808_0x0704_EncodeDecode(t *testing.T) {
	message := protocol.T808_0x0704{
		Type: 1,
		Items: []protocol.T808_0x0200{
			{
				Alarm:     2342,
				Status:    8,
				Lat:       decimal.NewFromFloat(23.562345),
				Lng:       decimal.NewFromFloat(-128.323123),
				Altitude:  2345,
				Speed:     160,
				Direction: 72,
				Time:      time.Unix(time.Now().Unix(), 0),
			},
		},
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x0704
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
	assert.True(t, reflect.DeepEqual(message, message2))
}
