package extra

import (
	"github.com/flash520/jtt808/errors"
)

// 无线通信网络信号强度
type Extra_0x30 struct {
	serialized []byte
	value      byte
}

func NewExtra_0x30(val byte) *Extra_0x30 {
	extra := Extra_0x30{
		value: val,
	}
	extra.serialized = []byte{val}
	return &extra
}

func (Extra_0x30) ID() byte {
	return byte(TypeExtra_0x30)
}

func (extra Extra_0x30) Data() []byte {
	return extra.serialized
}

func (extra Extra_0x30) Value() interface{} {
	return extra.value
}

func (extra *Extra_0x30) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, errors.ErrInvalidExtraLength
	}
	extra.value = data[0]
	return 1, nil
}
