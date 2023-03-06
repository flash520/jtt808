package extra

import (
	"github.com/flash520/jtt808/errors"
)

// GNSS定位卫星数
type Extra_0x31 struct {
	serialized []byte
	value      byte
}

func NewExtra_0x31(val byte) *Extra_0x31 {
	extra := Extra_0x31{
		value: val,
	}
	extra.serialized = []byte{val}
	return &extra
}

func (Extra_0x31) ID() byte {
	return byte(TypeExtra_0x31)
}

func (extra Extra_0x31) Data() []byte {
	return extra.serialized
}

func (extra Extra_0x31) Value() interface{} {
	return extra.value
}

func (extra *Extra_0x31) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, errors.ErrInvalidExtraLength
	}
	extra.value = data[0]
	return 1, nil
}
