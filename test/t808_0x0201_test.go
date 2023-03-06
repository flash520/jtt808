package test

import (
	"testing"
	"time"
	
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
	"github.com/flash520/jtt808/protocol/extra"
)

func TestT808_0x0201_EncodeDecode(t *testing.T) {
	areaID := uint32(121313)
	message := protocol.T808_0x0201{
		ReplyMsgSerialNo: 123,
		Result: protocol.T808_0x0200{
			Alarm:     2342,
			Status:    0,
			Lat:       decimal.NewFromFloat(23.562345),
			Lng:       decimal.NewFromFloat(-128.323123),
			Altitude:  2345,
			Speed:     160,
			Direction: 72,
			Time:      time.Unix(time.Now().Unix(), 0),
			Extras: []extra.Entity{
				extra.NewExtra_0x01(10),
				extra.NewExtra_0x2A(16),
				extra.NewExtra_0x2B(10086),
				extra.NewExtra_0x02(32),
				extra.NewExtra_0x03(64),
				extra.NewExtra_0x04(128),
				extra.NewExtra_0x11(extra.Extra_0x11_Value{
					Type:   1,
					AreaID: &areaID,
				}),
				extra.NewExtra_0x12(extra.Extra_0x12_Value{
					Type:      1,
					AreaID:    234,
					Direction: 1,
				}),
				extra.NewExtra_0x13(extra.Extra_0x13_Value{
					RoadID:   1006545,
					Duration: 234,
					Result:   1,
				}),
				extra.NewExtra_0x25(3344322),
				extra.NewExtra_0x30(12),
				extra.NewExtra_0x31(23),
			},
		},
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x0201
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
}
