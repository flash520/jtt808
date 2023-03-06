package extra

import (
	"bytes"
	"encoding/binary"
	
	"github.com/flash520/jtt808/errors"
)

// 超速报警
type Extra_0x11 struct {
	serialized []byte
	value      Extra_0x11_Value
}

type Extra_0x11_Value struct {
	Type   byte
	AreaID *uint32
}

func NewExtra_0x11(value Extra_0x11_Value) *Extra_0x11 {
	extra := Extra_0x11{
		value: value,
	}
	buffer := bytes.NewBuffer(nil)
	buffer.WriteByte(value.Type)
	
	if value.Type != 0 && value.AreaID != nil {
		var temp [4]byte
		binary.BigEndian.PutUint32(temp[:], uint32(*value.AreaID))
		buffer.Write(temp[:])
	}
	extra.serialized = buffer.Bytes()
	return &extra
}

func (Extra_0x11) ID() byte {
	return byte(TypeExtra_0x11)
}

func (extra Extra_0x11) Data() []byte {
	return extra.serialized
}

func (extra Extra_0x11) Value() interface{} {
	return extra.value
}

func (extra *Extra_0x11) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, errors.ErrInvalidExtraLength
	}
	extra.value.AreaID = nil
	extra.value.Type = data[0]
	if extra.value.Type == 0 {
		return 1, nil
	}
	
	if len(data[1:]) < 4 {
		return 0, errors.ErrInvalidExtraLength
	}
	areaID := binary.BigEndian.Uint32(data[1:])
	extra.value.AreaID = &areaID
	return 5, nil
}
