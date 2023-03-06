package extra

import (
	"bytes"
	"encoding/binary"
	
	"github.com/flash520/jtt808/errors"
)

// 路段行驶报警
type Extra_0x13 struct {
	serialized []byte
	value      Extra_0x13_Value
}

type Extra_0x13_Value struct {
	RoadID   uint32
	Duration uint16
	Result   byte
}

func NewExtra_0x13(value Extra_0x13_Value) *Extra_0x13 {
	extra := Extra_0x13{
		value: value,
	}
	
	var temp [4]byte
	buffer := bytes.NewBuffer(nil)
	binary.BigEndian.PutUint32(temp[:4], value.RoadID)
	buffer.Write(temp[:])
	
	binary.BigEndian.PutUint16(temp[:2], value.Duration)
	buffer.Write(temp[:2])
	
	buffer.WriteByte(value.Result)
	
	extra.serialized = buffer.Bytes()
	return &extra
}

func (Extra_0x13) ID() byte {
	return byte(TypeExtra_0x13)
}

func (extra Extra_0x13) Data() []byte {
	return extra.serialized
}

func (extra Extra_0x13) Value() interface{} {
	return extra.value
}

func (extra *Extra_0x13) Decode(data []byte) (int, error) {
	if len(data) < 7 {
		return 0, errors.ErrInvalidExtraLength
	}
	
	extra.value.RoadID = binary.BigEndian.Uint32(data[:4])
	extra.value.Duration = binary.BigEndian.Uint16(data[4:6])
	extra.value.Result = data[6]
	return 7, nil
}
