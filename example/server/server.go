package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	
	log "github.com/sirupsen/logrus"
	
	"github.com/flash520/jtt808"
	"github.com/flash520/jtt808/protocol"
	"github.com/flash520/jtt808/protocol/extra"
)

var privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDYKDhr0jA+/Fvcf/Fcv7Xn8PC2lkHCUjqE2LdMXAN+4RVve4PC
hh1dllgECC4doxK5ukEUijJeCwUrzLW2hzoxgN//Z9B6+LuW6twjF4vAVy1M6Abo
R69KNp2o34tDXm+Vo+4dgcHgBrH9s2ANqY++ImV0se+y6cd4Eo6bqfop/QIDAQAB
AoGBAMJBuxri5WrlfmS2MqIoxACy3pEojeZl4aNb47bjBl0zSQFMXkgmISPnJihR
dag60mxJP42G+ObdPoNzUGa+NoN5Gn8ro8esrBrOiUBdFIl7e24xQKd6skxoelIX
1dc2SzN5vOKFLAbu7r6GCohIpKsi0Icsldgh3REWJhmvkKKhAkEA+JSNOO2kxX5g
9kbZ8rfhUYN/xWWsCRXBcewocVlEUuit7QlZZeYYhAzoF2UlX7YjbPhWKEJIGPua
vLiMg+9+iQJBAN6b6tE59CRHd3JeSfs2x4Czudeb0hdOxdR2wnecOexP/LncWrRI
HOGqqdZFkAsrFyhZ3k2+Eff3fuV7GEg+UtUCQBUtGobB/+pvJLV2PbTmo0Q9bpIT
Yj934f3hf2SAlUh21/I8fKgonOgK7W6oyDFKI+Rxl21gkCHItVrkYdwPd/kCQQCU
fWTRU9srKBDhVUv8KrpBe6GH1QT7TyxfYSivKKLqoyBtyjMm9sNtNK49pAFFseSs
oeXL7fGGeq1G3imAZzJRAkEA7rIZYUjRmmqiNTFILkFo6OZ9cYd7AyZWg2KckUQ4
pwdWR1prOW0BORYozJc3ZloSN4JZP5Lu2gCerRNGRdenDw==
-----END RSA PRIVATE KEY-----
`

func GetTestPrivateKey() (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("get private key error")
	}
	
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return privateKey, nil
	}
	
	privateKey2, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey2.(*rsa.PrivateKey), nil
}

// 处理终端鉴权
func handleAuthentication(session *jtt808.Session, message *protocol.Message) {
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

// 处理上报位置
func handleReportLocation(session *jtt808.Session, message *protocol.Message) {
	// 打印消息
	entity := message.Body.(*protocol.T808_0x0200)
	fields := log.Fields{
		"IccID": message.Header.IccID,
		"警告":  fmt.Sprintf("0x%x", entity.Alarm),
		"状态":  fmt.Sprintf("0x%x", entity.Status),
		"纬度":  entity.Lat,
		"经度":  entity.Lng,
		"海拔":  entity.Altitude,
		"速度":  entity.Speed,
		"方向":  entity.Direction,
		"时间":  entity.Time,
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

// 处理上传媒体
func handleUploadMediaPacket(session *jtt808.Session, message *protocol.Message) {
	entity := message.Body.(*protocol.T808_0x0801)
	
	// 读取完整数据包
	fullPacket := make([]byte, 1024*1024)
	n, _ := entity.Packet.Read(fullPacket[:])
	fmt.Println(n)
	
	session.Send(&protocol.T808_0x8800{
		MediaID: entity.MediaID,
	})
}

// 处理终端 RSA公钥
func handleUploadRsaPublicKey(session *jtt808.Session, message *protocol.Message) {
	// 设置终端公钥
	entity := message.Body.(*protocol.T808_0x0A00)
	session.SetPublicKey(entity.PublicKey)
	
	// 回复平台应答
	session.Reply(message, protocol.T808_0x8100_ResultSuccess)
	
	// 下发平台公钥
	key := session.GetServer().GetPrivateKey()
	if key != nil {
		session.Send(&protocol.T808_0x8A00{
			PublicKey: &key.PublicKey,
		})
	}
}

func main() {
	privateKey, err := GetTestPrivateKey()
	if err != nil {
		panic(err)
	}
	
	server, _ := jtt808.NewServer(jtt808.Options{
		Keepalive:       60,
		AutoMergePacket: true,
		CloseHandler:    nil,
		PrivateKey:      privateKey,
	})
	server.AddHandler(protocol.MsgT808_0x0102, handleAuthentication)
	server.AddHandler(protocol.MsgT808_0x0200, handleReportLocation)
	server.AddHandler(protocol.MsgT808_0x0801, handleUploadMediaPacket)
	server.AddHandler(protocol.MsgT808_0x0A00, handleUploadRsaPublicKey)
	server.Run("tcp", 8808)
}
