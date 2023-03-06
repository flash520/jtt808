package test

import (
	"reflect"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

func TestT808_0x0702_EncodeDecode(t *testing.T) {
	message := protocol.T808_0x0702{
		State:        21,
		Time:         time.Unix(time.Now().Unix(), 0),
		ICCardResult: 1,
		DriverName:   "王师傅",
		Number:       "NS12345678",
		CompanyName:  "机构名称",
		ExpiryDate:   time.Unix(time.Now().Unix(), 0),
	}
	data, err := message.Encode()
	if err != nil {
		assert.Error(t, err, "encode error")
	}
	
	var message2 protocol.T808_0x0702
	_, err = message2.Decode(data)
	if err != nil {
		assert.Error(t, err, "decode error")
	}
	message2.ExpiryDate = message.ExpiryDate
	assert.True(t, reflect.DeepEqual(message, message2))
}
