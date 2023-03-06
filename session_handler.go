package jtt808

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	
	"github.com/flash520/link"
	log "github.com/sirupsen/logrus"
	
	"github.com/flash520/jtt808/protocol"
)

// Session处理
type sessionHandler struct {
	server          *Server
	autoMergePacket bool
}

func (handler sessionHandler) HandleSession(sess *link.Session) {
	log.WithFields(log.Fields{
		"id": sess.ID(),
	}).Debug("[JT/T 808] new session created")
	
	// 创建Session
	session := newSession(handler.server, sess)
	handler.server.mutex.Lock()
	handler.server.sessions[sess.ID()] = session
	handler.server.mutex.Unlock()
	handler.server.timer.Update(strconv.FormatUint(session.ID(), 10))
	sess.AddCloseCallback(nil, nil, func() {
		handler.server.handleClose(session)
	})
	
	for {
		// 接收消息
		msg, err := sess.Receive()
		if err != nil {
			sess.Close()
			break
		}
		
		// 分发消息
		message := msg.(protocol.Message)
		if message.Body == nil || reflect.ValueOf(message.Body).IsNil() {
			session.Reply(&message, protocol.T808_0x8001ResultUnsupported)
			continue
		}
		
		if !handler.autoMergePacket || !message.Header.Property.IsEnablePacket() {
			session.message(&message)
			handler.server.dispatchMessage(session, &message)
			continue
		}
		
		// 处理分包消息
		entityPacket, ok := interface{}(message.Body).(protocol.EntityPacket)
		if !ok {
			session.message(&message)
			handler.server.dispatchMessage(session, &message)
			continue
		}
		
		multipartFile := MultipartFile{
			IccID: message.Header.IccID,
			MsgID: message.Header.MsgID,
			Tag:   entityPacket.GetTag(),
			Sum:   message.Header.Packet.Sum,
		}
		buf, err := ioutil.ReadAll(entityPacket.GetReader())
		if err != nil {
			log.WithFields(log.Fields{
				"iccid":  message.Header.IccID,
				"msgid":  fmt.Sprintf("0x%x", message.Header.MsgID),
				"seq":    message.Header.Packet.Seq,
				"reason": err,
			}).Warn("[JT/T 808] failed to read packet data")
			session.Reply(&message, protocol.T808_0x8001ResultFail)
			continue
		}
		
		err = multipartFile.Write(message.Header.Packet.Seq, buf)
		if err != nil {
			log.WithFields(log.Fields{
				"iccid":  message.Header.IccID,
				"msgid":  fmt.Sprintf("0x%x", message.Header.MsgID),
				"seq":    message.Header.Packet.Seq,
				"reason": err,
			}).Warn("[JT/T 808] failed to write packet data to file")
			session.Reply(&message, protocol.T808_0x8001ResultFail)
			continue
		}
		
		session.Reply(&message, protocol.T808_0x8001ResultSuccess)
		if message.Header.Packet.Seq != message.Header.Packet.Sum || !multipartFile.IsFull() {
			continue
		}
		
		reader, err := multipartFile.Merge()
		if err != nil {
			log.WithFields(log.Fields{
				"iccid":  message.Header.IccID,
				"msgid":  fmt.Sprintf("0x%x", message.Header.MsgID),
				"reason": err,
			}).Warn("[JT/T 808] failed to merge packet file parts")
			continue
		}
		
		// 分发分包消息
		entityPacket.SetReader(reader)
		session.message(&message)
		handler.server.dispatchMessage(session, &message)
		reader.Close()
	}
}
