package extra

import (
	"bytes"
	"encoding/binary"
	
	"github.com/flash520/jtt808/errors"
)

// 进出报警
type Extra_0x12 struct {
	serialized []byte
	value      Extra_0x12_Value
}

type Extra_0x12_Value struct {
	Type      byte
	AreaID    uint32
	Direction byte
}

func NewExtra_0x12(value Extra_0x12_Value) *Extra_0x12 {
	extra := Extra_0x12{
		value: value,
	}
	
	var temp [4]byte
	buffer := bytes.NewBuffer(nil)
	buffer.WriteByte(value.Type)
	
	binary.BigEndian.PutUint32(temp[:], value.AreaID)
	buffer.Write(temp[:])
	
	buffer.WriteByte(value.Direction)
	
	extra.serialized = buffer.Bytes()
	return &extra
}

func (Extra_0x12) ID() byte {
	return byte(TypeExtra_0x12)
}

func (extra Extra_0x12) Data() []byte {
	return extra.serialized
}

func (extra Extra_0x12) Value() interface{} {
	return extra.value
}

func (extra *Extra_0x12) Decode(data []byte) (int, error) {
	if len(data) < 6 {
		return 0, errors.ErrInvalidExtraLength
	}
	
	extra.value.Type = data[0]
	extra.value.AreaID = binary.BigEndian.Uint32(data[1:5])
	extra.value.Direction = data[5]
	return 6, nil
}
