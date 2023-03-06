package jtt808

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net"
	"runtime/debug"
	"strconv"
	"sync"
	
	"github.com/flash520/link"
	log "github.com/sirupsen/logrus"
	
	"github.com/flash520/jtt808/protocol"
)

// Options 服务器选项
type Options struct {
	Keepalive       int64
	AutoMergePacket bool
	CloseHandler    func(*Session)
	PrivateKey      *rsa.PrivateKey
}

// Server 协议服务器
type Server struct {
	server     *link.Server
	handler    sessionHandler
	timer      *CountdownTimer
	privateKey *rsa.PrivateKey
	
	mutex    sync.Mutex
	sessions map[uint64]*Session
	
	closeHandler    func(*Session)
	messageHandlers sync.Map
}

type MessageHandler func(*Session, *protocol.Message)

// NewServer 创建服务
func NewServer(options Options) (*Server, error) {
	if options.Keepalive <= 0 {
		options.Keepalive = 60
	}
	
	if options.PrivateKey != nil && options.PrivateKey.Size() != 128 {
		return nil, errors.New("RSA key must be 1024 bits")
	}
	
	server := Server{
		closeHandler: options.CloseHandler,
		sessions:     make(map[uint64]*Session),
		privateKey:   options.PrivateKey,
	}
	server.handler.server = &server
	server.handler.autoMergePacket = options.AutoMergePacket
	server.timer = NewCountdownTimer(options.Keepalive, server.handleReadTimeout)
	return &server, nil
}

// Run 运行服务
func (server *Server) Run(network string, port int) error {
	if server.server != nil {
		return errors.New("server already running")
	}
	
	address := fmt.Sprintf("0.0.0.0:%d", port)
	listen, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	
	p := Protocol{
		privateKey: server.privateKey,
	}
	server.server = link.NewServer(listen, &p, 24, server.handler)
	log.Infof("[JT/T 808] protocol server started on %s", address)
	return server.server.Serve()
}

// 停止服务
func (server *Server) Stop() {
	if server.server != nil {
		server.server.Stop()
		server.server = nil
	}
}

// 获取Session
func (server *Server) GetSession(id uint64) (*Session, bool) {
	server.mutex.Lock()
	defer server.mutex.Unlock()
	session, ok := server.sessions[id]
	if !ok {
		return nil, false
	}
	return session, true
}

// 获取平台私钥
func (server *Server) GetPrivateKey() *rsa.PrivateKey {
	return server.privateKey
}

// 广播消息
func (server *Server) Broadcast(entity protocol.Entity) int {
	server.mutex.Lock()
	sessions := make([]*Session, 0, len(server.sessions))
	for _, session := range server.sessions {
		sessions = append(sessions, session)
	}
	server.mutex.Unlock()
	
	count := 0
	for _, session := range sessions {
		if _, err := session.Send(entity); err == nil {
			count++
		}
	}
	return count
}

// 添加消息处理
func (server *Server) AddHandler(msgID protocol.MsgID, handler MessageHandler) {
	if handler != nil {
		server.messageHandlers.Store(msgID, handler)
	}
}

// 处理关闭
func (server *Server) handleClose(session *Session) {
	server.mutex.Lock()
	delete(server.sessions, session.ID())
	server.mutex.Unlock()
	
	server.timer.Remove(strconv.FormatUint(session.ID(), 10))
	if server.closeHandler != nil {
		func() {
			defer func() {
				if err := recover(); err != nil {
					debug.PrintStack()
				}
			}()
			server.closeHandler(session)
		}()
	}
	
	log.WithFields(log.Fields{
		"id": session.ID(),
	}).Debug("[JT/T 808] session closed")
}

// 处理读超时
func (server *Server) handleReadTimeout(key string) {
	sessionID, err := strconv.ParseUint(key, 10, 64)
	if err != nil {
		return
	}
	
	session, ok := server.GetSession(sessionID)
	if !ok {
		return
	}
	session.Close()
	
	log.WithFields(log.Fields{
		"id": sessionID,
	}).Debug("[JT/T 808] session read timeout")
}

// 分派消息
func (server *Server) dispatchMessage(session *Session, message *protocol.Message) {
	log.WithFields(log.Fields{
		"id": fmt.Sprintf("0x%x", message.Header.MsgID),
	}).Debug("[JT/T 808] dispatch message")
	
	handler, ok := server.messageHandlers.Load(message.Header.MsgID)
	if !ok {
		log.WithFields(log.Fields{
			"id": fmt.Sprintf("0x%x", message.Header.MsgID),
		}).Info("[JT/T 808] dispatch message canceled, handler not found")
		return
	}
	
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()
	server.timer.Update(strconv.FormatUint(session.ID(), 10))
	handler.(MessageHandler)(session, message)
}
