// 启动一个tcpserver
package network

import (
	"fmt" // FIXME 第一个debug版本不使用持久化的日志方案
	"lib/util"
	"net"
	"sync"
)

type TcpServer struct {
	ServerAddress  string             // 服务器对外暴露的地址：localhost:9527
	MaxConnNum     int                // 最大的连接数
	ServerListener net.Listener       // 服务器的监听器
	PendingNum     int                // 最大发送队列长度（server -> client）
	MsgParser      *MessageParser     // 消息解析器
	Agent          func(*Conn) Iagent // 客户端代理
	MutexWG        sync.WaitGroup
}

// tcpserver启动入口
func (server *TcpServer) Start() {
	if ok := server.init(); !ok {
		fmt.Println("tcp server start failed")
		return
	}
	go server.run()
}

// 服务器初始化
func (server *TcpServer) init() bool {
	if server == nil {
		fmt.Println("init server is nil")
		return false
	}
	listener, err := net.Listen("tcp", server.ServerAddress)
	if err != nil {
		fmt.Println("net.Listen error:", err.Error())
		return false
	}
	if server.MaxConnNum <= 0 {
		server.MaxConnNum = 100
		fmt.Println("server.MaxConnNum <= 0, defalut 100")
	}
	if server.PendingNum <= 0 {
		server.PendingNum = 100
		fmt.Println("server.PendingNum <= 0, defalut 100")
	}
	server.ServerListener = listener
	server.MsgParser = NewMessageParser()
	server.MsgParser.SetMsgLen(2, 1024*5, 1)
	Init()
	return true
}

// 处理客户端的连接
func (server *TcpServer) run() {
	fmt.Println("tcp server pid:", util.GetPid())
	for {
		fmt.Println("loop accept")
		conn, err := server.ServerListener.Accept()
		// FIXME 这里需要对不同的错误情况不同处理，本版本暂时直接关闭
		if err != nil {
			fmt.Println("server.ServerListener.Accept error:", err.Error())
			server.Close()
			return
		}
		if ok := server.AddConn(conn, server.MaxConnNum); !ok {
			continue
		}
		tcpConn := InitConn(conn, server.PendingNum, server.MsgParser)
		agent := server.Agent(tcpConn)
		go func() {
			server.MutexWG.Add(1)
			defer server.MutexWG.Done()
			fmt.Println("启动一个代理携程循环执行代理:", util.GetPid())
			agent.Run()
			conn.Close()
			server.DeleteConn(conn)
			agent.OnClose()
		}()
	}
}

// 优雅的关闭
func (server *TcpServer) Close() {
	server.ServerListener.Close()
	Close()
	server.MutexWG.Wait()
}
