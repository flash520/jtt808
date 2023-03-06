package test

import (
	"reflect"
	"testing"
	"time"
	
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

func TestT808_0x8602_EncodeDecode(t *testing.T) {
	message := protocol.T808_0x8602{
		Action: protocol.AreaActionAdd,
		Items: []protocol.RectArea{
			{
				ID:             1,
				Attribute:      1,
				LeftTopLat:     decimal.NewFromFloat(123.456523),
				LeftTopLon:     decimal.NewFromFloat(-23.234212),
				RightBottomLat: decimal.NewFromFloat(23.456534),
				RightBottomLon: decimal.NewFromFloat(-13.456521),
				StartTime:      time.Unix(time.Now().Unix(), 0),
				EndTime:        time.Unix(time.Now().Unix(), 0),
				MaxSpeed:       1024,
				Duration:       60,
			},
			{
				ID:             2,
				Attribute:      1,
				LeftTopLat:     decimal.NewFromFloat(123.456523),
				LeftTopLon:     decimal.NewFromFloat(-23.234212),
				RightBottomLat: decimal.NewFromFloat(23.456534),
				RightBottomLon: decimal.NewFromFloat(-13.456521),
				StartTime:      time.Unix(time.Now().Unix(), 0),
				EndTime:        time.Unix(time.Now().Unix(), 0),
				MaxSpeed:       1024,
				Duration:       60,
			},
		},
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x8602
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
	assert.True(t, reflect.DeepEqual(message, message2))
}
