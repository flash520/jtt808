package test

import (
	"reflect"
	"testing"
	"time"
	
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

func TestT808_0x8606_EncodeDecode(t *testing.T) {
	timeTooLong := uint16(300)
	timeTooShort := uint16(60)
	maxSpeed := uint16(200)
	overSpeedDuration := byte(25)
	
	var attr protocol.T808_0x8606_SectionAttribute
	attr.SetTimeLimit(true)
	attr.SetSpeedLimit(true)
	attr.SetWestLongitude(true)
	
	message := protocol.T808_0x8606{
		ID:        1,
		Attribute: 123,
		StartTime: time.Unix(time.Now().Unix(), 0),
		EndTime:   time.Unix(time.Now().Unix(), 0),
		Points: []protocol.T808_0x8606_Point{
			{
				ID:                1,
				SectionID:         1,
				Lat:               decimal.NewFromFloat(23.562345),
				Lng:               decimal.NewFromFloat(-128.323123),
				Width:             15,
				Attribute:         attr,
				TimeTooLong:       &timeTooLong,
				TimeTooShort:      &timeTooShort,
				MaxSpeed:          &maxSpeed,
				OverSpeedDuration: &overSpeedDuration,
			},
		},
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x8606
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
	assert.True(t, reflect.DeepEqual(message, message2))
}
