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

type WsServer struct {
	ServerAddress  string               // 服务器对外暴露的地址：localhost:9527
	MaxConnNum     int                  // 最大的连接数
	MaxMsgLen      int                  // client单消息的最大长度
	ServerListener net.Listener         // 服务器的监听器
	PendingNum     int                  // 允许的最大的客户端连接的缓冲队列长度
	Agent          func(*WsConn) Iagent // 客户端代理
	Handler        *WsHandler           // wshandler
	HTTPTimeout    time.Duration        // http超时时间
	CertFile       string               // wss参数
	KeyFile        string               // wss参数
	MsgParser      *MessageParser       // 消息解析器
}

type WsHandler struct {
	MaxConnNum int                  // 最大的连接数
	MaxMsgLen  int                  // 单消息的最大长度
	MsgParser  *MessageParser       // 消息解析器
	PendingNum int                  // 允许的最大的客户端连接的缓冲队列长度
	Agent      func(*WsConn) Iagent // 客户端代理
	Upgrader   websocket.Upgrader   // 用于升级http连接
	Conns      WsConnMap            // WS连接池
	MutexWG    sync.WaitGroup
}

// HTTP路由处理函数
func (handler *WsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		logger.Error(fmt.Sprintf("Method not allowed, %s", r.Method))
		return
	}
	// 升级http->websocket
	conn, err := handler.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("upgrade error, %v", err))
		return
	}
	conn.SetReadLimit(int64(handler.MaxMsgLen))
	//	conn.SetReadDeadline(time.Now().Add(pongWait))
	//	conn.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	handler.MutexWG.Add(1)
	defer handler.MutexWG.Done()
	if ok := AddWsConn(conn, handler.MaxConnNum); !ok {
		conn.Close()
		return
	}
	wsConn := InitWsConn(conn, handler.PendingNum, handler.MsgParser)
	agent := handler.Agent(wsConn)
	agent.Run()
	wsConn.Close()
	DeleteConn(conn)
	agent.OnClose()
}

// 启动一个httpserver
func (server *WsServer) Start() {
	listener := server.init()
	if listener == nil {
		logger.Error(fmt.Sprintf("http server start failed, %v"))
		return
	}
	go server.run(listener)
}

// 初始化
func (server *WsServer) init() net.Listener {
	listener, err := net.Listen("tcp", server.ServerAddress)
	if err != nil {
		logger.Error(fmt.Sprintf("net.Listen error: %v", err))
		return nil
	}
	if server.MaxConnNum <= 0 {
		server.MaxConnNum = 100
		logger.Warning(fmt.Sprintf("server.MaxConnNum <= 0, defalut 100"))
	}
	if server.PendingNum <= 0 {
		server.PendingNum = 100
		logger.Warning(fmt.Sprintf("server.PendingNum <= 0, defalut 100"))
	}
	if server.MaxMsgLen <= 0 {
		server.MaxMsgLen = 1024
		logger.Warning(fmt.Sprintf("server.MaxMsgLen <= 0, defalut 1024"))
	}
	if server.HTTPTimeout <= 0 {
		server.HTTPTimeout = 10 * time.Second
		logger.Warning(fmt.Sprintf("server.HTTPTimeout <= 0, defalut 10s"))
	}
	if server.Agent == nil {
		return nil
	}
	// 支持wss
	// TODO
	if server.CertFile != "" || server.KeyFile != "" {
		config := &tls.Config{}
		config.NextProtos = []string{"http/1.1"}
		var err error
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(server.CertFile, server.KeyFile)
		if err != nil {
			logger.Warning(fmt.Sprintf("wss error: %v", err))
		}
		listener = tls.NewListener(listener, config)
	}
	server.ServerListener = listener
	server.Handler = &WsHandler{
		MaxConnNum: server.MaxConnNum,
		MaxMsgLen:  server.MaxMsgLen,
		PendingNum: server.PendingNum,
		Agent:      server.Agent,
		Upgrader: websocket.Upgrader{
			HandshakeTimeout: server.HTTPTimeout,
			CheckOrigin:      func(_ *http.Request) bool { return true },
		},
	}
	server.MsgParser = NewMessageParser()
	server.MsgParser.SetMsgLen(2, uint32(server.MaxMsgLen), 1)
	server.Handler.MsgParser = server.MsgParser
	InitWsPool()
	NewSessionMap()
	return listener
}

// run httpserver
func (server *WsServer) run(listener net.Listener) {
	httpServer := &http.Server{
		Addr:           server.ServerAddress,
		Handler:        server.Handler,
		ReadTimeout:    server.HTTPTimeout,
		WriteTimeout:   server.HTTPTimeout,
		MaxHeaderBytes: 1024,
	}
	httpServer.Serve(listener)
}

// 优雅的关闭
func (server *WsServer) Close() {
	server.ServerListener.Close()
	CloseWsPool()
	server.Handler.MutexWG.Wait()
}
