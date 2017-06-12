// 启动一个tcpserver
package network

import (
	"fmt" // FIXME 第一个debug版本不使用持久化的日志方案
	"lib/util"
	"net"
)

type TcpServer struct {
	ServerAddress  string             // 服务器对外暴露的地址：localhost:9527
	MaxConnNum     int                // 最大的连接数
	ServerListener net.Listener       // 服务器的监听器
	PendingNum     int                // 允许的最大的客户端缓冲队列长度
	MsgParser      *MessageParser     // 消息解析器
	Agent          func(*Conn) Iagent // 客户端代理
}

// tcpserver启动入口
func (server *TcpServer) Start() {
	server.init()
	go server.run()
}

// 服务器初始化
func (server *TcpServer) init() {
	// 必要的检查
	if server == nil {
		fmt.Println("init server is nil")
		return
	}
	// 创建服务器监听器
	listener, err := net.Listen("tcp", server.ServerAddress)
	if err != nil {
		fmt.Println("net.Listen error:", err.Error())
		return
	}
	// 必要的检查
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
	Init()
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
		// 更新连接池
		if ok := server.AddConn(conn, server.MaxConnNum); !ok {
			continue
		}
		// 初始化客户端conn
		tcpConn := InitConn(conn, server.PendingNum, server.MsgParser)
		agent := server.Agent(tcpConn)
		go func() {
			fmt.Println("启动一个代理携程循环执行run:", util.GetPid())
			// 循环反序列化并路由消息
			agent.Run()
			conn.Close()
			// 更新连接池
			server.DeleteConn(conn)
			// agent的清理工作
			agent.OnClose()
		}()
	}
}

// 优雅的关闭
func (server *TcpServer) Close() {
	// 不再接受连接
	server.ServerListener.Close()
	// 关闭已有的连接
	Close()
}
