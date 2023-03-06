package protocol

import (
	"github.com/shopspring/decimal"
	"math"
	"time"
)

// 区域属性
type AreaAttribute uint16

// 设置南纬
func (attr *AreaAttribute) SetSouthLatitude(b bool) {
	SetBitUint16((*uint16)(attr), 6, b)
}

// 设置西经
func (attr *AreaAttribute) SetWestLongitude(b bool) {
	SetBitUint16((*uint16)(attr), 7, b)
}

// 获取纬度类型
func (attr AreaAttribute) GetLatitudeType() LatitudeType {
	if GetBitUint16(uint16(attr), 6) {
		return SouthLatitudeType
	}
	return NorthLatitudeType
}

// 获取经度类型
func (attr AreaAttribute) GetLongitudeType() LongitudeType {
	if GetBitUint16(uint16(attr), 7) {
		return SouthLatitudeType
	}
	return NorthLatitudeType
}

// 设置离开报警平台
func (attr *AreaAttribute) SetExitAlarm(b bool) {
	SetBitUint16((*uint16)(attr), 5, b)
}

// 设置进入报警平台
func (attr *AreaAttribute) SetEnterAlarm(b bool) {
	SetBitUint16((*uint16)(attr), 3, b)
}

// 区域动作
type AreaAction byte

var (
	AreaActionUpdate AreaAction = 0
	AreaActionAdd    AreaAction = 1
	AreaActionEdit   AreaAction = 2
)

// 圆形区域
type CircleArea struct {
	// 区域ID
	ID uint32
	// 区域属性
	Attribute AreaAttribute
	// 中心点纬度
	Lat decimal.Decimal
	// 中心点经度
	Lng decimal.Decimal
	// 半径
	Radius uint32
	// 起始时间
	StartTime time.Time
	// 结束时间
	EndTime time.Time
	// 最高速度
	MaxSpeed uint16
	// 超速持续时间
	Duration byte
}

// 设置圆形区域
type T808_0x8600 struct {
	// 设置属性
	Action AreaAction
	// 区域项
	Items []CircleArea
}

func (entity *T808_0x8600) MsgID() MsgID {
	return MsgT808_0x8600
}

func (entity *T808_0x8600) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入设置属性
	writer.WriteByte(byte(entity.Action))

	// 写入区域总数
	writer.WriteByte(byte(len(entity.Items)))

	// 写入区域信息
	for idx := range entity.Items {
		item := &entity.Items[idx]

		// 写入区域ID
		writer.WriteUint32(item.ID)

		// 计算经纬度
		mul := decimal.NewFromFloat(1000000)
		lat := item.Lat.Mul(mul).IntPart()
		if lat < 0 {
			item.Attribute.SetSouthLatitude(true)
		}
		lng := item.Lng.Mul(mul).IntPart()
		if lng < 0 {
			item.Attribute.SetWestLongitude(true)
		}

		// 写入区域属性
		writer.WriteUint16(uint16(item.Attribute))

		// 写入中心点纬度
		writer.WriteUint32(uint32(math.Abs(float64(lat))))

		// 写入中心点经度
		writer.WriteUint32(uint32(math.Abs(float64(lng))))

		// 写入半径
		writer.WriteUint32(item.Radius)

		// 写入时间参数
		if item.Attribute&1 == 0 {
			continue
		}

		// 写入开始时间
		writer.WriteBcdTime(item.StartTime)

		// 写入结束时间
		writer.WriteBcdTime(item.EndTime)

		// 写入最高速度
		writer.WriteUint16(item.MaxSpeed)

		// 写入持续时间
		writer.WriteByte(item.Duration)
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x8600) Decode(data []byte) (int, error) {
	if len(data) < 20 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取设置属性
	action, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	entity.Action = AreaAction(action)

	// 读取区域总数
	count, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取区域信息
	entity.Items = make([]CircleArea, 0, count)
	for i := 0; i < int(count); i++ {
		var area CircleArea

		// 读取区域ID
		area.ID, err = reader.ReadUint32()
		if err != nil {
			return 0, err
		}

		// 读取区域属性
		attribute, err := reader.ReadUint16()
		if err != nil {
			return 0, err
		}
		area.Attribute = AreaAttribute(attribute)

		// 读取中心点纬度
		lat, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}

		// 读取中心点经度
		lng, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		area.Lat, area.Lng = getGeoPoint(
			lat, area.Attribute.GetLatitudeType() == SouthLatitudeType,
			lng, area.Attribute.GetLongitudeType() == WestLongitudeType)

		// 读取半径
		area.Radius, err = reader.ReadUint32()
		if err != nil {
			return 0, err
		}

		// 读取时间参数
		if area.Attribute&1 == 0 {
			entity.Items = append(entity.Items, area)
			continue
		}

		// 读取开始时间
		area.StartTime, err = reader.ReadBcdTime()
		if err != nil {
			return 0, err
		}

		// 读取结束时间
		area.EndTime, err = reader.ReadBcdTime()
		if err != nil {
			return 0, err
		}

		// 读取最高速度
		area.MaxSpeed, err = reader.ReadUint16()
		if err != nil {
			return 0, err
		}

		// 读取持续时间
		area.Duration, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}
		entity.Items = append(entity.Items, area)
	}
	return len(data) - reader.Len(), nil
}
