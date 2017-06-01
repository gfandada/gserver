// 封装了tcp网关的操作
// 其实是实现了接口imodule
package gateway

import (
	"fmt"
	"gserver/src/lib/network"
)

// 定义网关结构体
type Gate struct {
	MaxConnNum    int    // 允许的最大的连接数
	PendingNum    int    // 允许的最大的客户端缓冲队列长度
	MaxMsgLen     uint32 // 允许的最大的消息长度
	ServerAddress string // tcp服务地址
	LenMsgLen     int    // tcp消息长度
}

// tcp服务器
var tcpServer *network.TcpServer

// 启动tcp网关
func (gate *Gate) Run(chClose chan bool) {
	fmt.Println("启动一个tcp网关")
	tcpServer.Start()
	<-chClose
	tcpServer.Close()
}

// 初始化
func (gate *Gate) OnInit() {
	fmt.Println("执行tcp网关的初始化")
	// FIXME 这里没有对gate的内容做验证
	if gate == nil {
		fmt.Println("tcp_gateway run failed, because gate is nil")
		return
	}
	tcpServer = new(network.TcpServer)
	tcpServer.ServerAddress = gate.ServerAddress
	tcpServer.MaxConnNum = gate.MaxConnNum
	tcpServer.PendingNum = gate.PendingNum
}

// 资源回收
func (gate *Gate) OnDestroy() {
	fmt.Println("tcp网关销毁中")
}
