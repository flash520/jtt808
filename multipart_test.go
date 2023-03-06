package jtt808

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/flash520/jtt808/protocol"
)

func TestMultipartFile(t *testing.T) {
	multipartFile := MultipartFile{
		IccID: 19923476579,
		MsgID: protocol.MsgT808_0x0100,
		Sum:   3,
		Tag:   123456,
	}
	multipartFile.Write(2, []byte{0x03, 0x04})
	multipartFile.Write(1, []byte{0x01, 0x02})
	multipartFile.Write(3, []byte{0x05, 0x06})
	assert.True(t, multipartFile.IsFull())
	reader, err := multipartFile.Merge()
	if err != nil {
		assert.Error(t, err)
	}
	reader.Close()
}
