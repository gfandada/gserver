package network

import (
	"net"

	"github.com/gorilla/websocket"
)

type Iagent interface {
	NewIagent() Iagent            // 拷贝构造器
	WriteMsg(msg interface{})     // 写消息
	LocalAddr() net.Addr          // 本地地址
	RemoteAddr() net.Addr         // 远程地址
	Close()                       // 关闭
	Start(*websocket.Conn)        // 启动
	GetUserData() interface{}     // 获取用户数据
	SetUserData(data interface{}) // 设置用户数据
}
