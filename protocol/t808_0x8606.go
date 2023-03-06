package protocol

import (
	"github.com/shopspring/decimal"
	"math"
	"time"
)

// 设置路线
type T808_0x8606 struct {
	// 路线 ID
	ID uint32
	// 路线属性
	Attribute uint16
	// 起始时间
	StartTime time.Time
	// 结束时间
	EndTime time.Time
	// 拐点项
	Points []T808_0x8606_Point
}

// 路段属性
type T808_0x8606_SectionAttribute byte

// 是否限时
func (attr *T808_0x8606_SectionAttribute) HasTimeLimit() bool {
	return GetBitByte(byte(*attr), 0)
}

// 是否限速
func (attr *T808_0x8606_SectionAttribute) HasSpeedLimit() bool {
	return GetBitByte(byte(*attr), 1)
}

// 设置限时
func (attr *T808_0x8606_SectionAttribute) SetTimeLimit(b bool) {
	SetBitByte((*byte)(attr), 0, b)
}

// 设置限速
func (attr *T808_0x8606_SectionAttribute) SetSpeedLimit(b bool) {
	SetBitByte((*byte)(attr), 1, b)
}

// 设置南纬
func (attr *T808_0x8606_SectionAttribute) SetSouthLatitude(b bool) {
	SetBitByte((*byte)(attr), 2, b)
}

// 设置西经
func (attr *T808_0x8606_SectionAttribute) SetWestLongitude(b bool) {
	SetBitByte((*byte)(attr), 3, b)
}

// 获取纬度类型
func (attr T808_0x8606_SectionAttribute) GetLatitudeType() LatitudeType {
	if GetBitByte(byte(attr), 2) {
		return SouthLatitudeType
	}
	return NorthLatitudeType
}

// 获取经度类型
func (attr T808_0x8606_SectionAttribute) GetLongitudeType() LongitudeType {
	if GetBitByte(byte(attr), 3) {
		return WestLongitudeType
	}
	return EastLongitudeType
}

// 路线拐点
type T808_0x8606_Point struct {
	// 拐点 ID
	ID uint32
	// 路段 ID
	SectionID uint32
	// 拐点纬度
	Lat decimal.Decimal
	// 拐点经度
	Lng decimal.Decimal
	// 路段宽度
	Width byte
	// 路段属性
	Attribute T808_0x8606_SectionAttribute
	// 路段行驶过长阈值
	TimeTooLong *uint16
	// 路段行驶过短阈值
	TimeTooShort *uint16
	// 路段最高速度
	MaxSpeed *uint16
	// 路段超速持续时间
	OverSpeedDuration *byte
}

func (entity *T808_0x8606) MsgID() MsgID {
	return MsgT808_0x8606
}

func (entity *T808_0x8606) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入路线ID
	writer.WriteUint32(entity.ID)

	// 写入路线属性
	writer.WriteUint16(entity.Attribute)

	// 写入开始时间
	writer.WriteBcdTime(entity.StartTime)

	// 写入结束时间
	writer.WriteBcdTime(entity.EndTime)

	// 写入拐点总数
	writer.WriteUint16(uint16(len(entity.Points)))

	// 写入拐点列表
	for _, point := range entity.Points {
		//  写入拐点ID
		writer.WriteUint32(point.ID)

		// 写入路段ID
		writer.WriteUint32(point.SectionID)

		// 计算经纬度
		mul := decimal.NewFromFloat(1000000)
		lat := point.Lat.Mul(mul).IntPart()
		if lat < 0 {
			point.Attribute.SetSouthLatitude(true)
		}
		lng := point.Lng.Mul(mul).IntPart()
		if lng < 0 {
			point.Attribute.SetWestLongitude(true)
		}

		// 写入拐点经度
		writer.WriteUint32(uint32(math.Abs(float64(lat))))

		// 写入拐点纬度
		writer.WriteUint32(uint32(math.Abs(float64(lng))))

		// 设置拐点属性
		point.Attribute.SetTimeLimit(point.TimeTooLong != nil && point.TimeTooShort != nil)
		point.Attribute.SetSpeedLimit(point.MaxSpeed != nil && point.OverSpeedDuration != nil)

		// 写入路段宽度
		writer.WriteByte(point.Width)

		// 写入路段属性
		writer.WriteByte(byte(point.Attribute))

		// 写入路段行驶阈值
		if point.TimeTooLong != nil && point.TimeTooShort != nil {
			writer.WriteUint16(*point.TimeTooLong)
			writer.WriteUint16(*point.TimeTooShort)
		}

		// 写入超速报警设置
		if point.MaxSpeed != nil && point.OverSpeedDuration != nil {
			writer.WriteUint16(*point.MaxSpeed)
			writer.WriteByte(*point.OverSpeedDuration)
		}
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x8606) Decode(data []byte) (int, error) {
	if len(data) < 20 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取路线ID
	var err error
	entity.ID, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	// 读取路线属性
	entity.Attribute, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取开始时间
	entity.StartTime, err = reader.ReadBcdTime()
	if err != nil {
		return 0, err
	}

	// 读取结束时间
	entity.EndTime, err = reader.ReadBcdTime()
	if err != nil {
		return 0, err
	}

	// 读取拐点总数
	count, err := reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取拐点列表
	for i := 0; i < int(count); i++ {
		var point T808_0x8606_Point

		// 读取拐点ID
		point.ID, err = reader.ReadUint32()
		if err != nil {
			return 0, err
		}

		// 读取路段ID
		point.SectionID, err = reader.ReadUint32()
		if err != nil {
			return 0, err
		}

		// 读取拐点纬度
		latitude, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}

		// 读取拐点经度
		longitude, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}

		// 读取路段宽度
		point.Width, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}

		// 读取路段属性
		attribute, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}
		point.Attribute = T808_0x8606_SectionAttribute(attribute)
		point.Lat, point.Lng = getGeoPoint(
			latitude, point.Attribute.GetLatitudeType() == SouthLatitudeType,
			longitude, point.Attribute.GetLongitudeType() == WestLongitudeType)

		// 读取路段行驶阈值
		if point.Attribute.HasTimeLimit() {
			timeTooLong, err := reader.ReadUint16()
			if err != nil {
				return 0, err
			}

			timeTooShort, err := reader.ReadUint16()
			if err != nil {
				return 0, err
			}

			point.TimeTooLong = &timeTooLong
			point.TimeTooShort = &timeTooShort
		}

		// 读取超速报警设置
		if point.Attribute.HasSpeedLimit() {
			maxSpeed, err := reader.ReadUint16()
			if err != nil {
				return 0, err
			}

			overSpeedDuration, err := reader.ReadByte()
			if err != nil {
				return 0, err
			}

			point.MaxSpeed = &maxSpeed
			point.OverSpeedDuration = &overSpeedDuration
		}
		entity.Points = append(entity.Points, point)
	}
	return len(data) - reader.Len(), nil
}
