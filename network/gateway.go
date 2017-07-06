// 封装了网关的操作
package network

import (
	"fmt"
	"net"
	"time"

	"github.com/gfandada/gserver/logger"
	"github.com/gfandada/gserver/network/protobuff"
)

// 定义网关结构体
type Gate struct {
	MaxConnNum       int      // 允许的最大的连接数
	PendingNum       int      // 最大发送队列长度（server -> client）
	MaxMsgLen        int      // 允许的服务器接收的最大的消息长度
	MessageProcessor Imessage // 用于消息体的处理

	ServerAddress string // tcp服务地址
	LenMsgLen     int    // tcp消息长度

	WsServerAddress string        // websocket服务地址
	HTTPTimeout     time.Duration // http超时时间
	CertFile        string        // wws
	KeyFile         string        // wws
}

type Agent struct {
	Conn     Iconn       // 套接字套作接口
	Gate     *Gate       // 网关配置数据
	UserData interface{} // 用户数据
}

// tcp服务器
var tcpServer *TcpServer

// ws服务器
var wsServer *WsServer

/***************************实现imodule接口*******************************/

// 启动网关
func (gate *Gate) Run(chClose chan bool) {
	switch {
	case tcpServer != nil:
		tcpServer.Start()
		<-chClose
		tcpServer.Close()
	case wsServer != nil:
		wsServer.Start()
		<-chClose
		wsServer.Close()
	}
}

// 初始化
func (gate *Gate) OnInit() {
	if gate == nil {
		logger.Error("tcp_gateway run failed, because gate is nil")
		return
	}
	switch {
	case gate.ServerAddress != "":
		tcpServer = new(TcpServer)
		tcpServer.ServerAddress = gate.ServerAddress
		tcpServer.MaxConnNum = gate.MaxConnNum
		tcpServer.PendingNum = gate.PendingNum
		tcpServer.Agent = func(conn *Conn) Iagent {
			arg := &Agent{Conn: conn, Gate: gate}
			return arg
		}
	case gate.WsServerAddress != "":
		wsServer = new(WsServer)
		wsServer.ServerAddress = gate.WsServerAddress
		wsServer.MaxConnNum = gate.MaxConnNum
		wsServer.MaxMsgLen = gate.MaxMsgLen
		wsServer.PendingNum = gate.PendingNum
		wsServer.HTTPTimeout = gate.HTTPTimeout
		wsServer.CertFile = gate.CertFile
		wsServer.KeyFile = gate.KeyFile
		wsServer.Agent = func(conn *WsConn) Iagent {
			arg := &Agent{Conn: conn, Gate: gate}
			return arg
		}
	}
}

// 资源回收
func (gate *Gate) OnDestroy() {
	logger.Error(fmt.Sprintf("gateway OnDestroy, %v", gate))
}

/****************************实现了Iagent接口**********************************/

func (agent *Agent) Run() {
	if agent.Gate == nil {
		logger.Error(fmt.Sprintf("agent Run params is nil, %v", agent))
		return
	}
	for {
		msg, err := agent.Conn.ReadMsg()
		if err != nil {
			fmt.Println(err)
			break
		}
		if agent.Gate.MessageProcessor != nil {
			realMsg, errs := agent.Gate.MessageProcessor.Deserialize(msg)
			if errs != nil {
				logger.Error(fmt.Sprintf("Deserialize err:%v", errs))
				break
			}
			if err := agent.Gate.MessageProcessor.Router(realMsg, agent); err != nil {
				logger.Error(fmt.Sprintf("msg route err:%v", errs))
				break
			}
		}
	}
}

func (agent *Agent) OnClose() {
	logger.Info(fmt.Sprintf("agent OnClose:%v", agent))
}

/****************************实现了Igateway接口**********************************/

func (agent *Agent) WriteMsg(msg protobuff.RawMessage) {
	if agent.Gate.MessageProcessor != nil {
		data, err := agent.Gate.MessageProcessor.Serialize(msg)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = agent.Conn.WriteMsg(data...)
		if err != nil {
			logger.Error(fmt.Sprintf("write message %v error: %v", msg.MsgId, err))
		}
	}
}

func (agent *Agent) LocalAddr() net.Addr {
	if agent.Conn != nil {
		return agent.Conn.LocalAddr()
	}
	return nil
}

func (agent *Agent) RemoteAddr() net.Addr {
	if agent.Conn != nil {
		return agent.Conn.RemoteAddr()
	}
	return nil
}

func (agent *Agent) Close() {
	if agent.Conn != nil {
		agent.Conn.Close()
	}
}

func (agent *Agent) Destroy() {
	if agent.Conn != nil {
		agent.Conn.Destroy()
	}
}

func (agent *Agent) GetUserData() interface{} {
	return agent.UserData
}

func (agent *Agent) SetUserData(data interface{}) {

}

/******************************实现了Iack接口*********************************/

// TODO 需要优化
func (agent *Agent) Ack(data []interface{}) {
	if data == nil {
		return
	}
	switch len(data) {
	case 1:
		// ack自己
		agent.WriteMsg(data[0].(protobuff.RawMessage))
		return
	case 2:
		// ack自己
		agent.WriteMsg(data[0].(protobuff.RawMessage))
		// 更新session
		SetSession(data[1].(*Session))
		return
	}
}
