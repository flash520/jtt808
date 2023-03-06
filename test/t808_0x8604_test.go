package test

import (
	"reflect"
	"testing"
	"time"
	
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

func TestT808_0x8604_EncodeDecode(t *testing.T) {
	message := protocol.T808_0x8604{
		ID:        1,
		Attribute: 1,
		StartTime: time.Unix(time.Now().Unix(), 0),
		EndTime:   time.Unix(time.Now().Unix(), 0),
		MaxSpeed:  1024,
		Duration:  60,
		Vertexes: []protocol.Vertex{
			{Lat: decimal.NewFromFloat(123.234561), Lng: decimal.NewFromFloat(-23.432567)},
			{Lat: decimal.NewFromFloat(23.341098), Lng: decimal.NewFromFloat(-12.213435)},
		},
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x8604
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
	assert.True(t, reflect.DeepEqual(message, message2))
}
