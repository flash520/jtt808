package test

import (
	"reflect"
	"testing"
	"time"
	
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

func TestT808_0x8600_EncodeDecode(t *testing.T) {
	message := protocol.T808_0x8600{
		Action: protocol.AreaActionAdd,
		Items: []protocol.CircleArea{
			{
				ID:        1,
				Attribute: 1,
				Lat:       decimal.NewFromFloat(123.456789),
				Lng:       decimal.NewFromFloat(34.567891),
				Radius:    100,
				StartTime: time.Unix(time.Now().Unix(), 0),
				EndTime:   time.Unix(time.Now().Unix(), 0),
				MaxSpeed:  1024,
				Duration:  60,
			},
			{
				ID:        2,
				Attribute: 1,
				Lat:       decimal.NewFromFloat(123.456789),
				Lng:       decimal.NewFromFloat(34.567891),
				Radius:    100,
				StartTime: time.Unix(time.Now().Unix(), 0),
				EndTime:   time.Unix(time.Now().Unix(), 0),
				MaxSpeed:  1024,
				Duration:  60,
			},
		},
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x8600
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
	assert.True(t, reflect.DeepEqual(message, message2))
}
