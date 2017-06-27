// 封装了tcp网关的操作
// 其实是实现了接口imodule
package gateway

import (
	"fmt"
	"lib/network"
	"lib/network/protobuff"
	"lib/util"
	"net"
	"time"
)

// 定义网关结构体
type Gate struct {
	MaxConnNum       int              // 允许的最大的连接数
	PendingNum       int              // 最大发送队列长度（server -> client）
	MaxMsgLen        int              // 允许的最大的消息长度
	MessageProcessor network.Imessage // 用于消息体的处理

	ServerAddress string // tcp服务地址
	LenMsgLen     int    // tcp消息长度

	WsServerAddress string        // websocket服务地址
	HTTPTimeout     time.Duration // http超时时间
	CertFile        string        // wws
	KeyFile         string        // wws
}

type Agent struct {
	Conn     network.Iconn // 套接字套作接口
	Gate     *Gate         // 网关配置数据
	UserData interface{}   // 用户数据
}

// tcp服务器
var tcpServer *network.TcpServer

// ws服务器
var wsServer *network.WsServer

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
		fmt.Println("tcp_gateway run failed, because gate is nil")
		return
	}
	switch {
	case gate.ServerAddress != "":
		tcpServer = new(network.TcpServer)
		tcpServer.ServerAddress = gate.ServerAddress
		tcpServer.MaxConnNum = gate.MaxConnNum
		tcpServer.PendingNum = gate.PendingNum
		tcpServer.Agent = func(conn *network.Conn) network.Iagent {
			fmt.Println("NewAgent:", util.GetPid())
			arg := &Agent{Conn: conn, Gate: gate}
			return arg
		}
	case gate.WsServerAddress != "":
		wsServer = new(network.WsServer)
		wsServer.ServerAddress = gate.WsServerAddress
		wsServer.MaxConnNum = gate.MaxConnNum
		wsServer.MaxMsgLen = gate.MaxMsgLen
		wsServer.PendingNum = gate.PendingNum
		wsServer.HTTPTimeout = gate.HTTPTimeout
		wsServer.CertFile = gate.CertFile
		wsServer.KeyFile = gate.KeyFile
		wsServer.Agent = func(conn *network.WsConn) network.Iagent {
			arg := &Agent{Conn: conn, Gate: gate}
			return arg
		}
	}
}

// 资源回收
func (gate *Gate) OnDestroy() {
	fmt.Println("网关销毁中")
}

/****************************实现了Iagent接口**********************************/

func (agent *Agent) Run() {
	fmt.Println("客户端代理执行的run函数:", util.GetPid())
	if agent.Gate == nil {
		fmt.Println("Run params is nil")
		return
	}
	for {
		// 读取消息
		msg, err := agent.Conn.ReadMsg()
		if err != nil {
			fmt.Println(err)
			break
		}
		if agent.Gate.MessageProcessor != nil {
			// 反序列化
			realMsg, errs := agent.Gate.MessageProcessor.Deserialize(msg)
			if errs != nil {
				fmt.Println(errs)
				break
			}
			// 消息路由
			if err := agent.Gate.MessageProcessor.Router(realMsg, agent); err != nil {
				fmt.Println("msg route err:", err)
				break
			}
		}
	}
}

func (agent *Agent) OnClose() {
	fmt.Println("代理被销毁")
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
			fmt.Printf("write message %v error: %v \n", msg.MsgId, err)
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
