package protocol

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// CAN总线数据上传
type T808_0x0705 struct {
	// CAN总线数据接收时间
	ReceiveTime time.Duration
	// CAN总线数据项
	Items []T808_0x0705_CAN
}

// CAN总线数据项
type T808_0x0705_CAN struct {
	// CAN ID
	// bit31 表示 CAN通道号， 0：CAN1，1：CAN2
	// bit30 表示帧类型， 0：标准帧， 1：扩展帧
	// bit29 表示数据采集方式， 0：原始数据， 1：采集区间的平均值
	// bit28-bit0 表示 CAN总线 ID
	ID [4]byte
	// CAN数据
	Data [8]byte
}

func (entity *T808_0x0705) MsgID() MsgID {
	return MsgT808_0x0705
}

func (entity *T808_0x0705) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入数据项个数
	writer.WriteUint16(uint16(len(entity.Items)))

	// 写入数据接收时间
	if entity.ReceiveTime > time.Hour*24 {
		entity.ReceiveTime = time.Hour * 24
	}
	t := time.Time{}.Add(entity.ReceiveTime)
	s := fmt.Sprintf("%02d%02d%02d%04d",
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/int(time.Millisecond))
	writer.Write(stringToBCD(s), 5)

	// 写入CAN总线数据项
	for _, item := range entity.Items {
		writer.Write(item.ID[:])
		writer.Write(item.Data[:])
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x0705) Decode(data []byte) (int, error) {
	if len(data) < 7 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取数据项个数
	count, err := reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取数据接收时间
	buf, err := reader.Read(5)
	if err != nil {
		return 0, err
	}
	s := bcdToString(buf, true)
	if len(s) != 10 {
		return 0, errors.New("invalid BTC time")
	}
	hour, _ := strconv.Atoi(s[0:2])
	minute, _ := strconv.Atoi(s[2:4])
	second, _ := strconv.Atoi(s[4:6])
	millisecond, _ := strconv.Atoi(s[6:10])
	entity.ReceiveTime = time.Duration(hour)*time.Hour +
		time.Duration(minute)*time.Minute +
		time.Duration(second)*time.Second +
		time.Duration(millisecond)*time.Millisecond

	// 读取CAN总线数据项
	entity.Items = make([]T808_0x0705_CAN, 0, count)
	for i := 0; i < int(count); i++ {
		id, err := reader.Read(4)
		if err != nil {
			return 0, err
		}

		data, err := reader.Read(8)
		if err != nil {
			return 0, err
		}

		var item T808_0x0705_CAN
		copy(item.ID[:], id)
		copy(item.Data[:], data)
		entity.Items = append(entity.Items, item)
	}
	return len(data) - reader.Len(), nil
}
