// 通用的tcp服务器
package network

import (
	"fmt"
	"net"
	"sync"

	"github.com/gfandada/gserver/logger"
)

type TcpServer struct {
	serverAddress  string
	maxConnNum     int
	maxMsgLen      int
	minMsgLen      int
	pendingNum     int
	readTimeout    int
	writeTimeout   int
	msgParser      Imessage
	serverListener net.Listener
	gate           Iagent
	conns          map[net.Conn]struct{}
	mutexConns     sync.Mutex
	wgConns        sync.WaitGroup
	mutexWG        sync.WaitGroup
}

func StartTcp(config *Config) *TcpServer {
	server := new(TcpServer)
	server.serverAddress = config.ServerAddress
	server.maxConnNum = config.MaxConnNum
	server.maxMsgLen = config.MaxMsgLen
	server.minMsgLen = config.MinMsgLen
	server.pendingNum = config.PendingNum
	server.readTimeout = config.ReadDeadline
	server.writeTimeout = config.WriteDeadline
	server.msgParser = config.MsgParser
	server.gate = config.Gate
	server.conns = make(map[net.Conn]struct{}, config.MaxConnNum)
	server.start()
	return server
}

func (server *TcpServer) start() {
	listener := server.init()
	if listener == nil {
		logger.Error(fmt.Sprintf("http server start failed, %v"))
		return
	}
	go server.run(listener)
}

func (server *TcpServer) init() net.Listener {
	listener, err := net.Listen("tcp", server.serverAddress)
	if err != nil {
		logger.Error(fmt.Sprintf("net.Listen error: %v", err))
		return nil
	}
	if server.maxConnNum <= 0 {
		server.maxConnNum = 1000
		logger.Warning(fmt.Sprintf("server.maxConnNum <= 0, defalut 1000"))
	}
	if server.pendingNum <= 0 {
		server.pendingNum = 100
		logger.Warning(fmt.Sprintf("server.pendingNum <= 0, defalut 100"))
	}
	if server.maxMsgLen <= 0 {
		server.maxMsgLen = 512
		logger.Warning(fmt.Sprintf("server.maxMsgLen <= 0, defalut 512"))
	}
	if server.minMsgLen < 0 {
		server.minMsgLen = 0
		logger.Warning(fmt.Sprintf("server.minMsgLen < 0, defalut 0"))
	}
	if server.readTimeout <= 0 {
		server.readTimeout = 10
		logger.Warning(fmt.Sprintf("server.readTimeout <= 0, defalut 10s"))
	}
	if server.writeTimeout <= 0 {
		server.writeTimeout = 10
		logger.Warning(fmt.Sprintf("server.writeTimeout <= 0, defalut 10s"))
	}
	server.serverListener = listener
	if server.msgParser == nil {
		logger.Error("server.msgParser is nil %v", server)
		return nil
	}
	return listener
}

func (server *TcpServer) run(listener net.Listener) {
	server.mutexWG.Add(1)
	defer server.mutexWG.Done()
	for {
		conn, err := server.serverListener.Accept()
		if err != nil {
			return
		}
		server.mutexConns.Lock()
		if len(server.conns) >= server.maxConnNum {
			server.mutexConns.Unlock()
			conn.Close()
			continue
		}
		server.conns[conn] = struct{}{}
		server.mutexConns.Unlock()
		go func() {
			server.wgConns.Add(1)
			defer server.wgConns.Done()
			defer func() {
				if r := recover(); r != nil {
					// TODO
					fmt.Println("tcp agent error:", r)
				}
			}()
			// for agent
			server.gate.NewIagent().Start(&TcpConn{conn: conn})
		}()
	}
}

func (server *TcpServer) Close() {
	server.serverListener.Close()
	server.mutexWG.Wait()
	server.mutexConns.Lock()
	for conn := range server.conns {
		conn.Close()
	}
	server.conns = nil
	server.mutexConns.Unlock()
	server.wgConns.Wait()
}
