// 封装了tcp网关的操作
// 其实是实现了接口imodule
package gateway

import (
	"fmt"
	"lib/chanrpc"
	"lib/network"
	"lib/util"
)

// 定义网关结构体
type Gate struct {
	MaxConnNum       int               // 允许的最大的连接数
	PendingNum       int               // 允许的最大的客户端缓冲队列长度
	MaxMsgLen        uint32            // 允许的最大的消息长度
	ServerAddress    string            // tcp服务地址
	LenMsgLen        int               // tcp消息长度
	AgentChanRPC     *chanrpc.Server   // 用于路由的rpcserver
	MessageProcessor network.Processor // 用于消息体的处理
}

type agent struct {
	conn     network.Conn //客户端连接
	gate     *Gate        // 网关数据
	userData interface{}  // 用户数据
}

// tcp服务器
var tcpServer *network.TcpServer

/***************************实现imodule接口*******************************/

// 启动tcp网关
func (gate *Gate) Run(chClose chan bool) {
	tcpServer.Start()
	<-chClose
	tcpServer.Close()
}

// 初始化
func (gate *Gate) OnInit() {
	// FIXME 这里没有对gate的内容做验证
	if gate == nil {
		fmt.Println("tcp_gateway run failed, because gate is nil")
		return
	}
	tcpServer = new(network.TcpServer)
	tcpServer.ServerAddress = gate.ServerAddress
	tcpServer.MaxConnNum = gate.MaxConnNum
	tcpServer.PendingNum = gate.PendingNum
	tcpServer.Agent = func(conn *network.Conn) network.Iagent {
		fmt.Println("NewAgent:", util.GetPid())
		arg := &agent{conn: *conn, gate: gate}
		if gate.AgentChanRPC != nil {
			gate.AgentChanRPC.Go("NewAgent", arg)
		}
		return arg
	}
}

// 资源回收
func (gate *Gate) OnDestroy() {
	fmt.Println("tcp网关销毁中")
}

/****************************实现了Iagent接口**********************************/

func (agent *agent) Run() {
	fmt.Println("客户端代理执行的run函数")
	if agent.gate == nil {
		fmt.Println("Run params is nil")
		return
	}
	for {
		// 读取消息
		msg, err := agent.conn.ReadMsg()
		if err != nil {
			fmt.Println(err)
			break
		}
		// 反序列化
		realMsg, errs := agent.gate.MessageProcessor.Deserialize(msg)
		if errs != nil {
			fmt.Println(errs)
			break
		}
		// 消息路由
		if err := agent.gate.MessageProcessor.Route(realMsg, agent); err != nil {
			fmt.Println("msg route err:", err)
			break
		}
	}
}

func (agent *agent) OnClose() {

}
