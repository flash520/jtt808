package extra

// 附加消息类型
type Type byte

const (
	// 里程
	TypeExtra_0x01 Type = 0x01
	// 油量
	TypeExtra_0x02 Type = 0x02
	// 速度
	TypeExtra_0x03 Type = 0x03
	// 报警确认
	TypeExtra_0x04 Type = 0x04
	// 超速报警
	TypeExtra_0x11 Type = 0x11
	// 进出区域报警
	TypeExtra_0x12 Type = 0x12
	// 路段行驶时间报警
	TypeExtra_0x13 Type = 0x13
	//扩展车辆信号状态位
	TypeExtra_0x25 Type = 0x25
	// IO状态位
	TypeExtra_0x2a Type = 0x2a
	// 模拟量
	TypeExtra_0x2b Type = 0x2b
	// 无线通信网络信号强度
	TypeExtra_0x30 Type = 0x30
	// GNSS定位卫星数
	TypeExtra_0x31 Type = 0x31
)

// 消息实体映射
var entityMapper = map[byte]func() Entity{
	byte(TypeExtra_0x01): func() Entity {
		return new(Extra_0x01)
	},
	byte(TypeExtra_0x02): func() Entity {
		return new(Extra_0x02)
	},
	byte(TypeExtra_0x03): func() Entity {
		return new(Extra_0x03)
	},
	byte(TypeExtra_0x04): func() Entity {
		return new(Extra_0x04)
	},
	byte(TypeExtra_0x11): func() Entity {
		return new(Extra_0x11)
	},
	byte(TypeExtra_0x12): func() Entity {
		return new(Extra_0x12)
	},
	byte(TypeExtra_0x13): func() Entity {
		return new(Extra_0x13)
	},
	byte(TypeExtra_0x25): func() Entity {
		return new(Extra_0x25)
	},
	byte(TypeExtra_0x2a): func() Entity {
		return new(Extra_0x2A)
	},
	byte(TypeExtra_0x2b): func() Entity {
		return new(Extra_0x2B)
	},
	byte(TypeExtra_0x30): func() Entity {
		return new(Extra_0x30)
	},
	byte(TypeExtra_0x31): func() Entity {
		return new(Extra_0x31)
	},
}

// 类型注册
func Register(typ byte, creator func() Entity) {
	entityMapper[typ] = creator
}
