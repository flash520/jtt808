package jtt808

import (
	"crypto/rsa"
	"runtime/debug"
	"sync"
	"sync/atomic"

	"github.com/flash520/link"
	log "github.com/sirupsen/logrus"

	"github.com/flash520/jtt808/protocol"
)

// 请求上下文
type requestContext struct {
	msgID    uint16
	serialNo uint16
	callback func(answer *protocol.Message)
}

// 终端会话
type Session struct {
	next    uint32
	iccID   uint64
	server  *Server
	session *link.Session

	mux      sync.Mutex
	requests []requestContext

	UserData interface{}
}

// 创建Session
func newSession(server *Server, sess *link.Session) *Session {
	return &Session{
		server:  server,
		session: sess,
	}
}

// 获取ID
func (session *Session) ID() uint64 {
	return session.session.ID()
}

// 获取服务实例
func (session *Session) GetServer() *Server {
	return session.server
}

// 获取RSA公钥
func (session *Session) GetPublicKey() *rsa.PublicKey {
	codec, ok := session.session.Codec().(*ProtocolCodec)
	if !ok || codec == nil {
		return nil
	}
	return codec.GetPublicKey()
}

// 设置RSA公钥
func (session *Session) SetPublicKey(publicKey *rsa.PublicKey) {
	codec, ok := session.session.Codec().(*ProtocolCodec)
	if !ok || codec == nil {
		return
	}
	codec.SetPublicKey(publicKey)
}

// 发送消息
func (session *Session) Send(entity protocol.Entity) (uint16, error) {
	message := protocol.Message{
		Body: entity,
		Header: protocol.Header{
			MsgID:       entity.MsgID(),
			IccID:       atomic.LoadUint64(&session.iccID),
			MsgSerialNo: session.nextID(),
		},
	}

	err := session.session.Send(message)
	if err != nil {
		return 0, err
	}
	return message.Header.MsgSerialNo, nil
}

// 回复消息
func (session *Session) Reply(msg *protocol.Message, result protocol.Result) (uint16, error) {
	entity := protocol.T808_0x8001{
		ReplyMsgSerialNo: msg.Header.MsgSerialNo,
		ReplyMsgID:       msg.Header.MsgID,
		Result:           result,
	}
	return session.Send(&entity)
}

// Request 发起请求
func (session *Session) Request(entity protocol.Entity, cb func(answer *protocol.Message)) (uint16, error) {
	serialNo, err := session.Send(entity)
	if err != nil {
		return 0, err
	}

	if cb != nil {
		session.addRequestContext(requestContext{
			msgID:    uint16(entity.MsgID()),
			serialNo: serialNo,
			callback: cb,
		})
	}
	return serialNo, nil
}

// Close 关闭连接
func (session *Session) Close() error {
	return session.session.Close()
}

// 获取消息ID
func (session *Session) nextID() uint16 {
	var id uint32
	for {
		id = atomic.LoadUint32(&session.next)
		if id == 0xff {
			if atomic.CompareAndSwapUint32(&session.next, id, 1) {
				id = 1
				break
			}
		} else if atomic.CompareAndSwapUint32(&session.next, id, id+1) {
			id += 1
			break
		}
	}
	return uint16(id)
}

// 消息接收事件
func (session *Session) message(message *protocol.Message) {
	if message.Header.IccID > 0 {
		old := atomic.LoadUint64(&session.iccID)
		if old != 0 && old != message.Header.IccID {
			log.WithFields(log.Fields{
				"id":  session.ID(),
				"old": old,
				"new": message.Header.IccID,
			}).Warn("[JT/T 808] terminal IccID is inconsistent")
		}
		atomic.StoreUint64(&session.iccID, message.Header.IccID)
	}

	var msgSerialNo uint16
	switch message.Header.MsgID {
	case protocol.MsgT808_0x0001:
		// 终端通用应答
		msgSerialNo = message.Body.(*protocol.T808_0x0001).ReplyMsgSerialNo
	case protocol.MsgT808_0x0104:
		// 查询终端参数应答
		msgSerialNo = message.Body.(*protocol.T808_0x0104).ReplyMsgSerialNo
	case protocol.MsgT808_0x0201:
		// 位置信息查询应答
		msgSerialNo = message.Body.(*protocol.T808_0x0201).ReplyMsgSerialNo
	case protocol.MsgT808_0x0302:
		// 提问应答
		msgSerialNo = message.Body.(*protocol.T808_0x0302).ReplyMsgSerialNo
	case protocol.MsgT808_0x0500:
		// 车辆控制应答
		msgSerialNo = message.Body.(*protocol.T808_0x0500).ReplyMsgSerialNo
	case protocol.MsgT808_0x0700:
		// 行驶记录数据上传
		msgSerialNo = message.Body.(*protocol.T808_0x0700).ReplyMsgSerialNo
	case protocol.MsgT808_0x0802:
		// 存储多媒体数据检索应答
		msgSerialNo = message.Body.(*protocol.T808_0x0802).ReplyMsgSerialNo
	case protocol.MsgT808_0x0805:
		// 摄像头立即拍摄命令应答
		msgSerialNo = message.Body.(*protocol.T808_0x0805).ReplyMsgSerialNo
	}
	if msgSerialNo == 0 {
		return
	}

	ctx, ok := session.takeRequestContext(msgSerialNo)
	if ok {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		ctx.callback(message)
	}
}

// 添加请求上下文
func (session *Session) addRequestContext(ctx requestContext) {
	session.mux.Lock()
	defer session.mux.Unlock()

	for idx, item := range session.requests {
		if item.msgID == ctx.msgID {
			session.requests[idx] = ctx
			return
		}
	}
	session.requests = append(session.requests, ctx)
}

// 取出请求上下文
func (session *Session) takeRequestContext(msgSerialNo uint16) (requestContext, bool) {
	session.mux.Lock()
	defer session.mux.Unlock()

	for idx, item := range session.requests {
		if item.serialNo == msgSerialNo {
			session.requests[idx] = session.requests[len(session.requests)-1]
			session.requests = session.requests[:len(session.requests)-1]
			return item, true
		}
	}
	return requestContext{}, false
}
