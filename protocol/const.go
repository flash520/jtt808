package protocol

const (
	// 标志位
	PrefixID = byte(0x7e)

	// 转义符
	EscapeByte = byte(0x7d)

	// 0x7d < ———— > 0x7d 后紧跟一个0x01
	EscapeByteSufix1 = byte(0x01)

	// 0x7e < ———— > 0x7d 后紧跟一个0x02
	EscapeByteSufix2 = byte(0x02)

	// 消息头大小
	MessageHeaderSize = 12
)
