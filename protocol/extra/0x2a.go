package extra

import (
	"encoding/binary"
	
	"github.com/flash520/jtt808/errors"
)

// IO状态位
type Extra_0x2A struct {
	serialized []byte
	value      uint16
}

func NewExtra_0x2A(val uint16) *Extra_0x2A {
	extra := Extra_0x2A{
		value: val,
	}
	
	var temp [2]byte
	binary.BigEndian.PutUint16(temp[:2], val)
	extra.serialized = temp[:2]
	return &extra
}

func (Extra_0x2A) ID() byte {
	return byte(TypeExtra_0x2a)
}

func (extra Extra_0x2A) Data() []byte {
	return extra.serialized
}

func (extra Extra_0x2A) Value() interface{} {
	return extra.value
}

func (extra *Extra_0x2A) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, errors.ErrInvalidExtraLength
	}
	extra.value = binary.BigEndian.Uint16(data)
	return 2, nil
}
