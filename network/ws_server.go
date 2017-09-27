// 通用的websocket服务器
package network

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gfandada/gserver/logger"

	"github.com/gorilla/websocket"
)

type Config struct {
	ServerAddress string         // 服务地址
	MaxConnNum    int            // 最大连接数
	MaxMsgLen     int            // client->gateway message上限
	MinMsgLen     int            // client->gateway message下限
	ReadDeadline  int            // gateway->client读超时
	WriteDeadline int            // gateway->client写超时
	MaxHeader     int            // header上限(for websocket)
	HttpTimeout   int            // http-get超时(for websocket)
	CertFile      string         // for ssl
	KeyFile       string         // for ssl
	MsgParser     Imessage       // for message
	Parser        *MessageParser // for 报文
	Gate          Iagent         // 网关
	PendingNum    int            // gateway->client异步ipc队列上限
	Rpm           int            // client->gateway流量上限/min
	AsyncMQ       int            // service->gateway异步ipc队列上限
	GateWayIds    uint16         // gateway本地路由id段(当前路由规则是简单的id分段规则)
}

type WsServer struct {
	serverAddress  string
	maxHeader      int
	maxConnNum     int
	maxMsgLen      int
	minMsgLen      int
	pendingNum     int
	readTimeout    int
	writeTimeout   int
	httpTimeout    int
	certFile       string
	keyFile        string
	msgParser      Imessage
	handler        *wsHandler
	serverListener net.Listener
	gate           Iagent
}

type wsHandler struct {
	maxConnNum   int
	maxMsgLen    int
	minMsgLen    int
	pendingNum   int
	readTimeout  int
	writeTimeout int
	upgrader     websocket.Upgrader
	gate         Iagent
	conns        map[*websocket.Conn]struct{}
	mutexConns   sync.Mutex
	wgConns      sync.WaitGroup
	mutexWG      sync.WaitGroup
}

func StartWs(config *Config) *WsServer {
	server := new(WsServer)
	server.serverAddress = config.ServerAddress
	server.maxHeader = config.MaxHeader
	server.maxConnNum = config.MaxConnNum
	server.maxMsgLen = config.MaxMsgLen
	server.minMsgLen = config.MinMsgLen
	server.pendingNum = config.PendingNum
	server.readTimeout = config.ReadDeadline
	server.writeTimeout = config.WriteDeadline
	server.httpTimeout = config.HttpTimeout
	server.certFile = config.CertFile
	server.keyFile = config.KeyFile
	server.msgParser = config.MsgParser
	server.gate = config.Gate
	server.start()
	return server
}

func (server *WsServer) start() {
	listener := server.init()
	if listener == nil {
		logger.Error(fmt.Sprintf("websocket-server start failed %v", server))
		return
	}
	go server.run(listener)
}

func (handler *wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		logger.Error("websocket-server Method not allowed, %s", r.Method)
		return
	}
	conn, err := handler.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("websocket-server upgrade error, %v", err)
		return
	}
	conn.SetReadLimit(int64(handler.maxMsgLen + 8))
	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("websocket-server agent error %v", r))
		}
	}()
	handler.mutexConns.Lock()
	if len(handler.conns) >= handler.maxConnNum {
		handler.mutexConns.Unlock()
		conn.Close()
		return
	}
	handler.conns[conn] = struct{}{}
	handler.mutexConns.Unlock()
	handler.mutexWG.Add(1)
	defer handler.mutexWG.Done()
	handler.gate.NewIagent().Start(&WsConn{conn: conn})
}

func (server *WsServer) init() net.Listener {
	listener, err := net.Listen("tcp", server.serverAddress)
	if err != nil {
		logger.Error(fmt.Sprintf("websocket-server net.Listen error %v", err))
		return nil
	}
	if server.maxHeader <= 0 {
		server.maxHeader = 512
		logger.Warning(fmt.Sprintf("websocket-server server.maxHeader <= 0, defalut 512"))
	}
	if server.maxConnNum <= 0 {
		server.maxConnNum = 1024
		logger.Warning(fmt.Sprintf("websocket-server server.maxConnNum <= 0, defalut 1024"))
	}
	if server.pendingNum <= 0 {
		server.pendingNum = 64
		logger.Warning(fmt.Sprintf("websocket-server server.pendingNum <= 0, defalut 64"))
	}
	if server.maxMsgLen <= 0 {
		server.maxMsgLen = 512
		logger.Warning(fmt.Sprintf("websocket-server server.maxMsgLen <= 0, defalut 512"))
	}
	if server.minMsgLen < 0 {
		server.minMsgLen = 0
		logger.Warning(fmt.Sprintf("websocket-server server.minMsgLen < 0, defalut 0"))
	}
	if server.httpTimeout <= 0 {
		server.httpTimeout = 10
		logger.Warning(fmt.Sprintf("websocket-server server.httpTimeout <= 0, defalut 10s"))
	}
	if server.readTimeout <= 0 {
		server.readTimeout = 10
		logger.Warning(fmt.Sprintf("websocket-server server.readTimeout <= 0, defalut 10s"))
	}
	if server.writeTimeout <= 0 {
		server.writeTimeout = 10
		logger.Warning(fmt.Sprintf("websocket-server server.writeTimeout <= 0, defalut 10s"))
	}
	// for ssl
	if server.certFile != "" || server.keyFile != "" {
		config := &tls.Config{}
		config.NextProtos = []string{"http/1.1"}
		var err error
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(server.certFile, server.keyFile)
		if err != nil {
			logger.Warning(fmt.Sprintf("websocket-server ssl error %v", err))
		}
		listener = tls.NewListener(listener, config)
	}
	server.serverListener = listener
	server.handler = &wsHandler{
		maxConnNum:   server.maxConnNum,
		maxMsgLen:    server.maxMsgLen,
		minMsgLen:    server.minMsgLen,
		pendingNum:   server.pendingNum,
		readTimeout:  server.readTimeout,
		writeTimeout: server.writeTimeout,
		conns:        make(map[*websocket.Conn]struct{}, server.maxConnNum),
		upgrader: websocket.Upgrader{
			HandshakeTimeout: time.Duration(server.httpTimeout) * time.Second,
			CheckOrigin:      func(_ *http.Request) bool { return true },
		},
	}
	if server.msgParser == nil {
		logger.Error("websocket-server server.msgParser is nil")
		return nil
	}
	server.handler.gate = server.gate
	return listener
}

func (server *WsServer) run(listener net.Listener) {
	httpServer := &http.Server{
		Addr:           server.serverAddress,
		Handler:        server.handler,
		ReadTimeout:    time.Duration(server.readTimeout) * time.Second,
		WriteTimeout:   time.Duration(server.writeTimeout) * time.Second,
		MaxHeaderBytes: server.maxHeader,
	}
	httpServer.Serve(listener)
}

func (server *WsServer) Close() {
	logger.Info(fmt.Sprintf("websocket-server close %v Accept conns %d",
		server, len(server.handler.conns)))
	server.serverListener.Close()
	server.handler.mutexWG.Wait()
	server.handler.mutexConns.Lock()
	for conn := range server.handler.conns {
		conn.Close()
	}
	server.handler.conns = nil
	server.handler.mutexConns.Unlock()
	server.handler.wgConns.Wait()
}
