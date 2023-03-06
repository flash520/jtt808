package extra

import (
	"encoding/binary"
	
	"github.com/flash520/jtt808/errors"
)

// 扩展车辆信号状态位
type Extra_0x25 struct {
	serialized []byte
	value      uint32
}

func NewExtra_0x25(val uint32) *Extra_0x25 {
	extra := Extra_0x25{
		value: val,
	}
	
	var temp [4]byte
	binary.BigEndian.PutUint32(temp[:4], val)
	extra.serialized = temp[:4]
	return &extra
}

func (Extra_0x25) ID() byte {
	return byte(TypeExtra_0x25)
}

func (extra Extra_0x25) Data() []byte {
	return extra.serialized
}

func (extra Extra_0x25) Value() interface{} {
	return extra.value
}

func (extra *Extra_0x25) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, errors.ErrInvalidExtraLength
	}
	extra.value = binary.BigEndian.Uint32(data)
	return 4, nil
}
