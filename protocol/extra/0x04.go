package extra

import (
	"encoding/binary"
	
	"github.com/flash520/jtt808/errors"
)

// 报警确认
type Extra_0x04 struct {
	serialized []byte
	value      uint16
}

func NewExtra_0x04(val uint16) *Extra_0x04 {
	extra := Extra_0x04{
		value: val,
	}
	
	var temp [2]byte
	binary.BigEndian.PutUint16(temp[:2], val)
	extra.serialized = temp[:2]
	return &extra
}

func (Extra_0x04) ID() byte {
	return byte(TypeExtra_0x04)
}

func (extra Extra_0x04) Data() []byte {
	return extra.serialized
}

func (extra Extra_0x04) Value() interface{} {
	return extra.value
}

func (extra *Extra_0x04) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, errors.ErrInvalidExtraLength
	}
	extra.value = binary.BigEndian.Uint16(data)
	return 2, nil
}
