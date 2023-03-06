package protocol

import (
	"github.com/shopspring/decimal"
	"math"
	"time"
)

// 顶点
type Vertex struct {
	// 顶点纬度
	Lat decimal.Decimal
	// 顶点经度
	Lng decimal.Decimal
}

// 设置多边形区域
type T808_0x8604 struct {
	// 区域 ID
	ID uint32
	// 区域属性
	Attribute AreaAttribute
	// 起始时间
	StartTime time.Time
	// 结束时间
	EndTime time.Time
	// 最高速度
	MaxSpeed uint16
	// 超速持续时间
	Duration byte
	// 顶点项
	Vertexes []Vertex
}

func (entity *T808_0x8604) MsgID() MsgID {
	return MsgT808_0x8604
}

func (entity *T808_0x8604) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入区域ID
	writer.WriteUint32(entity.ID)

	// 计算经纬度
	mul := decimal.NewFromFloat(1000000)
	vertexes := make([]Vertex, 0, len(entity.Vertexes))
	for j := range entity.Vertexes {
		vertex := &entity.Vertexes[j]
		lat := vertex.Lat.Mul(mul)
		if lat.Cmp(decimal.Zero) < 0 {
			entity.Attribute.SetSouthLatitude(true)
		}
		lng := vertex.Lng.Mul(mul)
		if lng.Cmp(decimal.Zero) < 0 {
			entity.Attribute.SetWestLongitude(true)
		}
		vertexes = append(vertexes, Vertex{
			Lat: lat,
			Lng: lng,
		})
	}

	// 写入区域属性
	writer.WriteUint16(uint16(entity.Attribute))

	// 写入时间参数
	if entity.Attribute&1 == 1 {
		// 写入开始时间
		writer.WriteBcdTime(entity.StartTime)

		// 写入结束时间
		writer.WriteBcdTime(entity.EndTime)

		// 写入最高速度
		writer.WriteUint16(entity.MaxSpeed)

		// 写入持续时间
		writer.WriteByte(entity.Duration)
	}

	// 写入顶点总数
	writer.WriteUint16(uint16(len(entity.Vertexes)))

	// 写入顶点信息
	for _, vertex := range vertexes {
		// 写入纬度
		writer.WriteUint32(uint32(math.Abs(float64(vertex.Lat.IntPart()))))

		// 写入经度
		writer.WriteUint32(uint32(math.Abs(float64(vertex.Lng.IntPart()))))
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x8604) Decode(data []byte) (int, error) {
	if len(data) < 16 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取区域ID
	var err error
	entity.ID, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	// 读取区域属性
	attribute, err := reader.ReadUint16()
	if err != nil {
		return 0, err
	}
	entity.Attribute = AreaAttribute(attribute)

	// 读取时间参数
	if entity.Attribute&1 == 1 {
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

		// 读取最高速度
		entity.MaxSpeed, err = reader.ReadUint16()
		if err != nil {
			return 0, err
		}

		// 读取持续时间
		entity.Duration, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}
	}

	// 读取顶点总数
	vertexes, err := reader.ReadUint16()
	if err != nil {
		return 0, err
	}
	entity.Vertexes = make([]Vertex, 0, int(vertexes))

	// 读取顶点列表
	for j := 0; j < int(vertexes); j++ {
		var vertex Vertex

		// 读取纬度
		lat, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}

		// 读取经度
		lng, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}

		vertex.Lat, vertex.Lng = getGeoPoint(
			lat, entity.Attribute.GetLatitudeType() == SouthLatitudeType,
			lng, entity.Attribute.GetLongitudeType() == WestLongitudeType)
		entity.Vertexes = append(entity.Vertexes, vertex)
	}
	return len(data) - reader.Len(), nil
}
