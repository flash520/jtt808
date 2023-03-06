package extra

import (
	"encoding/binary"
	
	"github.com/flash520/jtt808/errors"
)

// 模拟量
type Extra_0x2B struct {
	serialized []byte
	value      uint32
}

func NewExtra_0x2B(val uint32) *Extra_0x2B {
	extra := Extra_0x2B{
		value: val,
	}
	
	var temp [4]byte
	binary.BigEndian.PutUint32(temp[:4], val)
	extra.serialized = temp[:4]
	return &extra
}

func (Extra_0x2B) ID() byte {
	return byte(TypeExtra_0x2b)
}

func (extra Extra_0x2B) Data() []byte {
	return extra.serialized
}

func (extra Extra_0x2B) Value() interface{} {
	return extra.value
}

func (extra *Extra_0x2B) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, errors.ErrInvalidExtraLength
	}
	extra.value = binary.BigEndian.Uint32(data)
	return 4, nil
}
