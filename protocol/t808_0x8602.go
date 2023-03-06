package protocol

import (
	"github.com/shopspring/decimal"
	"math"
	"time"
)

// 矩形区域
type RectArea struct {
	// 区域ID
	ID uint32
	// 区域属性
	Attribute AreaAttribute
	// 左上点纬度
	LeftTopLat decimal.Decimal
	// 左上点经度
	LeftTopLon decimal.Decimal
	// 右下点纬度
	RightBottomLat decimal.Decimal
	// 右下点经度
	RightBottomLon decimal.Decimal
	// 起始时间
	StartTime time.Time
	// 结束时间
	EndTime time.Time
	// 最高速度
	MaxSpeed uint16
	// 超速持续时间
	Duration byte
}

// 设置矩形区域
type T808_0x8602 struct {
	// 设置属性
	Action AreaAction
	// 区域项
	Items []RectArea
}

func (entity *T808_0x8602) MsgID() MsgID {
	return MsgT808_0x8602
}

func (entity *T808_0x8602) Encode() ([]byte, error) {
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
		leftLat := item.LeftTopLat.Mul(mul).IntPart()
		if leftLat < 0 {
			item.Attribute.SetSouthLatitude(true)
		}
		leftLon := item.LeftTopLon.Mul(mul).IntPart()
		if leftLon < 0 {
			item.Attribute.SetWestLongitude(true)
		}
		rightLat := item.RightBottomLat.Mul(mul).IntPart()
		if rightLat < 0 {
			item.Attribute.SetSouthLatitude(true)
		}
		rightLon := item.RightBottomLon.Mul(mul).IntPart()
		if rightLon < 0 {
			item.Attribute.SetWestLongitude(true)
		}

		// 写入区域属性
		writer.WriteUint16(uint16(item.Attribute))

		// 写入左上角纬度
		writer.WriteUint32(uint32(math.Abs(float64(leftLat))))

		// 写入左上角经度
		writer.WriteUint32(uint32(math.Abs(float64(leftLon))))

		// 写入右下角纬度
		writer.WriteUint32(uint32(math.Abs(float64(rightLat))))

		// 写入右下角经度
		writer.WriteUint32(uint32(math.Abs(float64(rightLon))))

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

func (entity *T808_0x8602) Decode(data []byte) (int, error) {
	if len(data) < 24 {
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
	entity.Items = make([]RectArea, 0, count)
	for i := 0; i < int(count); i++ {
		var area RectArea

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

		// 读取左上角纬度
		leftTopLat, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}

		// 读取左上角经度
		leftTopLon, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}

		// 读取右下角纬度
		rightBottomLat, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}

		// 读取右下角经度
		rightBottomLon, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}

		area.LeftTopLat, area.LeftTopLon = getGeoPoint(
			leftTopLat, area.Attribute.GetLatitudeType() == SouthLatitudeType,
			leftTopLon, area.Attribute.GetLongitudeType() == WestLongitudeType)
		area.RightBottomLat, area.RightBottomLon = getGeoPoint(
			rightBottomLat, area.Attribute.GetLatitudeType() == SouthLatitudeType,
			rightBottomLon, area.Attribute.GetLongitudeType() == WestLongitudeType)

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
