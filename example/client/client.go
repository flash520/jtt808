package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
	
	"github.com/shopspring/decimal"
	
	"github.com/flash520/jtt808"
	"github.com/flash520/jtt808/protocol"
)

var privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDi8hNwzJ8cgEZqFRBG3vQQhrRn+57enZMhc1F/C8jU2ywaZnWB
XCeWxdh5NMfoEMu0464Oetu2/94AvKEvr/C/tdIIZhqgvYOrGWOPAQ1XSCC29ldZ
0GpeaShiiQgKlmfyWYnWLQNmbgKTa6Wyu+nl1MuYWbKEuYsFC5pjdSIlMwIDAQAB
AoGBANZjpzfdJcZP67UVNu4setYN2umMS1Wz+CUWgnuJT2y9q9k4x3Kv8vo85rYB
xYOWMkos9+XX7C3hYwDBMWgSRlOcwEL8bjd4Tizwez9WTOAwjnwmZoZWQuDkeYzd
haWlSq6c2wc0H6G07nGWfrIbfAZ8GaqpzsaN6+mAZ1ZVJNn5AkEA/qxQNApmEFGC
KlEJCPYHxVSarbiPOgPOSt123n42f5wluPe81eU2vmFSDZ1Rrv4MDRDIqtWMJS1V
o5bdWYNkLwJBAOQgx4El4Ezvcw8I7lIuz8z3ssno3+EcP6mtLJ0ihQ2zsm4TcdYV
MYP45/+eKnOc2BBMPCQsreYuwXRCEPuhmj0CQQCxgomkvFrHpQiFVlZl2JcyA/aM
f8fVODHiHNtt2atC5yOj+Ym1zT6LFGqM8sqsnobn1HsKGC7G+wJmNBG1AtAhAkA/
omkkPFWCAHUe54XbBNXQPfPwYHY6y+9yPC0qs9tbhBmsnN3vMsA6KO9GHW+ICmM2
wJ0yFgh4Ieiyrk8gceadAkEAmS+XB0UWLMYlSjG2PGm06UAm/frvEFNu8742qQSJ
zS1yTivTUIflWm77kvUIiJJZP7UVwGTDSx3atOICtedimw==
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

func main() {
	privateKey, err := GetTestPrivateKey()
	if err != nil {
		panic(err)
	}
	
	tcpAddr, err := net.ResolveTCPAddr(
		"tcp", "127.0.0.1:8808")
	if err != nil {
		panic(err)
	}
	
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	
	// 终端鉴权
	message := protocol.Message{
		Header: protocol.Header{
			IccID:       19901234567,
			MsgSerialNo: 1,
		},
		Body: &protocol.T808_0x0102{
			AuthKey: "12345678",
		},
	}
	data, err := message.Encode()
	if err != nil {
		panic(data)
	}
	if _, err = conn.Write(data); err != nil {
		panic(err)
	}
	
	// 上报位置
	message = protocol.Message{
		Header: protocol.Header{
			IccID:       19901234567,
			MsgSerialNo: 2,
		},
		Body: &protocol.T808_0x0200{
			Alarm:     2342,
			Status:    0,
			Lat:       decimal.NewFromFloat(23.562345),
			Lng:       decimal.NewFromFloat(-128.323123),
			Altitude:  2345,
			Speed:     160,
			Direction: 72,
			Time:      time.Unix(time.Now().Unix(), 0),
		},
	}
	data, err = message.Encode()
	if err != nil {
		panic(data)
	}
	if _, err = conn.Write(data); err != nil {
		panic(err)
	}
	
	// 上传公钥
	message = protocol.Message{
		Header: protocol.Header{
			IccID:       19901234567,
			MsgSerialNo: 3,
		},
		Body: &protocol.T808_0x0A00{
			PublicKey: &privateKey.PublicKey,
		},
	}
	data, err = message.Encode()
	if err != nil {
		panic(data)
	}
	if _, err = conn.Write(data); err != nil {
		panic(err)
	}
	
	// 上传媒体文件
	offset := 0
	limit := 512
	mediaData := make([]byte, 1024*2)
	for {
		if offset >= len(mediaData) {
			break
		}
		
		if offset+limit > len(mediaData) {
			limit = (offset + limit) - len(mediaData)
		}
		
		sum := len(mediaData) / limit
		if len(mediaData)%limit > 0 {
			sum += 1
		}
		seq := (offset + limit) / limit
		if (offset+limit)%limit > 0 {
			seq += 1
		}
		
		header := protocol.Header{
			IccID:       19901234567,
			MsgSerialNo: 3,
		}
		header.Packet = &protocol.Packet{
			Sum: uint16(sum),
			Seq: uint16(seq),
		}
		message = protocol.Message{
			Header: header,
			Body: &protocol.T808_0x0801{
				MediaID:   1024,
				Type:      protocol.T808_0x0800_MediaTypeAudio,
				Coding:    protocol.T808_0x0800_MediaCodingJPEG,
				Event:     13,
				ChannelID: 28,
				Location: protocol.T808_0x0200{
					Alarm:     2342,
					Status:    0,
					Lat:       decimal.NewFromFloat(23.562345),
					Lng:       decimal.NewFromFloat(-128.323123),
					Altitude:  2345,
					Speed:     160,
					Direction: 72,
					Time:      time.Unix(time.Now().Unix(), 0),
				},
				Packet: bytes.NewReader(mediaData[offset : offset+limit]),
			},
		}
		data, err = message.Encode()
		if err != nil {
			panic(data)
		}
		if _, err = conn.Write(data); err != nil {
			panic(err)
		}
		offset += limit
	}
	
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go onMessageReceived(conn, &waitGroup)
	waitGroup.Wait()
}

func onMessageReceived(conn *net.TCPConn, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	
	var p jtt808.Protocol
	codec, err := p.NewCodec(conn)
	if err != nil {
		panic(err)
	}
	
	for {
		msg, err := codec.Receive()
		if err != nil {
			conn.Close()
			break
		}
		
		message := msg.(protocol.Message)
		if message.Header.MsgID == protocol.MsgT808_0x8104 {
			message = protocol.Message{
				Header: protocol.Header{
					IccID:       19901234567,
					MsgSerialNo: 1000,
				},
				Body: &protocol.T808_0x0104{
					ReplyMsgSerialNo: message.Header.MsgSerialNo,
					Params: []*protocol.Param{
						new(protocol.Param).SetByte(0x0084, 24),
						new(protocol.Param).SetBytes(0x0110, []byte{1, 2, 3, 4, 5, 6, 7, 8}),
						new(protocol.Param).SetUint16(0x0031, 100),
						new(protocol.Param).SetUint32(0x0046, 64000),
						new(protocol.Param).SetString(0x0083, "车牌号码"),
					},
				},
			}
			data, err := message.Encode()
			if err != nil {
				panic(data)
			}
			if _, err = conn.Write(data); err != nil {
				panic(err)
			}
		} else if message.Header.MsgID == protocol.MsgT808_0x8800 {
			fmt.Println("===========================> 媒体上传成功 <===========================")
		} else if message.Header.MsgID == protocol.MsgT808_0x8A00 {
			fmt.Println("===========================> 收到平台RSA公钥 <===========================")
		}
	}
}
