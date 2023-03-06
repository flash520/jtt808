# 部标 JT/T808(2013) 协议服务器
`go808` 是使用 Golang 语言开发的部标 JT/T808(2013) 协议服务器。此项目的目的是为了减少重复工作，加快开发速度。开发者无需关心 JT/T808 协议封包解包，自动进行字节流协议与结构体互转，只需专注具体业务的逻辑处理。

**特点**

易于扩展，开发者可自行添加自定义协议。

使用简单，几行代码即可启用一个支持 JT/T808 协议的服务器。

## 安装
```
$ go get github.com/flash520/jtt808
```

## 示例

### 1. 上报位置
```go
package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/flash520/jtt808"
	"github.com/flash520/jtt808/protocol"
	"github.com/flash520/jtt808/protocol/extra"
)

// 处理上报位置
func handleReportLocation(session *jtt808.Session, message *protocol.Message) {
	// 打印消息
	entity := message.Body.(*protocol.T808_0x0200)
	fields := log.Fields{
		"IccID": message.Header.IccID,
		"警告":    fmt.Sprintf("0x%x", entity.Alarm),
		"状态":    fmt.Sprintf("0x%x", entity.Status),
		"纬度":    entity.Lat,
		"经度":    entity.Lng,
		"海拔":    entity.Altitude,
		"速度":    entity.Speed,
		"方向":    entity.Direction,
		"时间":    entity.Time,
	}

	for _, ext := range entity.Extras {
		switch ext.ID() {
		case extra.Extra_0x01{}.ID():
			fields["行驶里程"] = ext.(*extra.Extra_0x01).Value()
		case extra.Extra_0x02{}.ID():
			fields["剩余油量"] = ext.(*extra.Extra_0x02).Value()
		}
	}
	log.WithFields(fields).Info("上报终端位置信息")

	// 回复平台应答
	session.Reply(message, protocol.T808_0x8100_ResultSuccess)
}

func main() {
	server, _ := go808.NewServer(go808.Options{
		Keepalive:       60,
	})
	server.AddHandler(protocol.MsgT808_0x0200, handleReportLocation)
	server.Run("tcp", 8808)
}
```

### 2. 下发->应答回调
```go
package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/flash520/jtt808"
	"github.com/flash520/jtt808/protocol"
)

// 处理终端鉴权
func handleAuthentication(session *go808.Session, message *protocol.Message) {
	// 回复平台应答
	session.Reply(message, protocol.T808_0x8100_ResultSuccess)

	// 查询终端参数
	session.Request(new(protocol.T808_0x8104), func(answer *protocol.Message) {
		response := answer.Body.(*protocol.T808_0x0104)
		for _, param := range response.Params {
			fmt.Println("参数ID", param.ID())
		}
	})
}

func main() {
	server, _ := go808.NewServer(go808.Options{
		Keepalive:       60,
	})
	server.AddHandler(protocol.MsgT808_0x0200, handleAuthentication)
	server.Run("tcp", 8808)
}
```

### 3. 分包合并
```go
package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/flash520/jtt808"
	"github.com/flash520/jtt808/protocol"
)

// 处理上传媒体
func handleUploadMediaPacket(session *go808.Session, message *protocol.Message) {
	entity := message.Body.(*protocol.T808_0x0801)

	// 读取完整数据包
	fullPacket := make([]byte, 1024*1024)
	entity.Packet.Read(fullPacket[:])

	session.Send(&protocol.T808_0x8800{
		MediaID: entity.MediaID,
	})
}

func main() {
	server, _ := go808.NewServer(go808.Options{
		Keepalive:       60,
		AutoMergePacket: true, // 自动合并分包
	})
	server.AddHandler(protocol.MsgT808_0x0801, handleUploadMediaPacket)
	server.Run("tcp", 8808)
}
```

### 4. 自定义协议
TODO...

## 开发进度

| 消息ID        | 消息名称| 是否完成 |
| :---: | :------: | :------: |
| 0x0001        |终端通用应答                  | √         |
| 0x8001        |平台通用应答                  | √         |
| 0x0002        |终端心跳                      | √         |
| 0x8003        |补传分包请求                  | √         |
| 0x0100        |终端注册                      | √         |
| 0x8100        |终端注册应答                  | √         |
| 0x0003        |终端注销                      | √         |
| 0x0102        |终端鉴权                      | √         |
| 0x8103        |设置终端参数                  | √         |
| 0x8104        |查询终端参数                  | √         |
| 0x0104        |查询终端参数应答              | √         |
| 0x8105        |终端控制                      | √         |
| 0x8106        |查询指定终端参数              | √         |
| 0x8107        |查询终端属性                  | √         |
| 0x0107        |查询终端属性应答              | √         |
| 0x8108        |下发终端升级包                | √         |
| 0x0108        |终端升级结果通知              | √         |
| 0x0200        |位置信息汇报                  | √         |
| 0x8201        |位置信息查询                  | √         |
| 0x0201        |位置信息查询应答              | √         |
| 0x8202        |临时位置跟踪控制              | √         |
| 0x8203        |人工确认报警消息              | √         |
| 0x8300        |文本信息下发                  | √         |
| 0x8301        |事件设置                      | √         |
| 0x0301        |事件报告                      | √         |
| 0x8302        |提问下发                      | √         |
| 0x0302        |提问应答                      | √         |
| 0x8303        |信息点播菜单设置              | √         |
| 0x0303        |信息点播/取消                 | √         |
| 0x8304        |信息服务                      | √         |
| 0x8400        |电话回拨                      | √         |
| 0x8401        |设置电话本                    | √         |
| 0x8500        |车辆控制                      | √         |
| 0x0500        |车辆控制应答                  | √         |
| 0x8600        |设置圆形区域                  | √         |
| 0x8601        |删除圆形区域                  | √         |
| 0x8602        |设置矩形区域                  | √         |
| 0x8603        |删除矩形区域                  | √         |
| 0x8604        |设置多边形区域                | √         |
| 0x8605        |删除多边形区域                | √         |
| 0x8606        |设置路线                      | √         |
| 0x8607        |删除路线                      | √         |
| 0x8700        |行驶记录仪数据采集命令        | √         |
| 0x0700        |行驶记录仪数据上传            | √         |
| 0x8701        |行驶记录仪参数下传命令        | √         |
| 0x0701        |电子运单上报                  | √         |
| 0x0702        |驾驶员身份信息采集上报        | √         |
| 0x8702        |上报驾驶员身份信息请求        | √         |
| 0x0704        |定位数据批量上传              | √         |
| 0x0705        |CAN 总线数据上传              | √         |
| 0x0800        |多媒体事件信息上传            | √         |
| 0x0801        |多媒体数据上传                | √         |
| 0x8800        |多媒体数据上传应答            | √         |
| 0x8801        |摄像头立即拍摄命令            | √         |
| 0x0805        |摄像头立即拍摄命令应答        | √         |
| 0x8802        |存储多媒体数据检索            | √         |
| 0x0802        |存储多媒体数据检索应答        | √         |
| 0x8803        |存储多媒体数据上传            | √         |
| 0x8804        |录音开始命令                  | √         |
| 0x8805        |单条存储多媒体数据检索上传命令| √         |
| 0x8900        |数据下行透传                  | √         |
| 0x0900        |数据上行透传                  | √         |
| 0x0901        |数据压缩上报                  | √         |
| 0x8A00        |平台 RSA 公钥                 | √         |
| 0x0A00        |终端 RSA 公钥                 | √         |
