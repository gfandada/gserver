// 封装了客户端会话的操作
package network

import (
	"fmt"
	"net"
	"sync"
)

// conn数据
type Conn struct {
	Conn      net.Conn       // 客户端的连接
	ChanSend  chan []byte    // 用来保存服务器需要发送给客户端的数据
	ChanRecv  chan []byte    // FIXME 用来保存接收到的客户端数据
	MsgParser *MessageParser // 消息解析器
	ChanMut   sync.Mutex     // 保证ChanSend携程安全
	Online    bool           // 是否在线的标记:true在线 false离线
}

// 初始化客户端conn
func InitConn(conn net.Conn, pendingNum int, msgParser *MessageParser) *Conn {
	tcpConn := new(Conn)
	tcpConn.Conn = conn
	tcpConn.ChanSend = make(chan []byte, pendingNum)
	tcpConn.MsgParser = msgParser
	tcpConn.Online = true
	go func() {
		fmt.Println("循环从chanSend读取msg，发送至client--sta")
		for sendData := range tcpConn.ChanSend {
			if sendData == nil {
				break
			}
			// 从chan中读出数据，发送给客户端conn
			_, err := conn.Write(sendData)
			if err != nil {
				break
			}
		}
		fmt.Println("循环从chanSend读取msg，发送至client--end")
		conn.Close()
	}()
	return tcpConn
}

/****************************实现了iagent接口******************************/

func (conn *Conn) Run() {
	for {
		// 解析消息
		_, err := conn.ReadMsg()
		// 严格的消息限制
		if err != nil {
			break
		}
		// TODO
		// 1 反序列化消息
		// 2 消息路由
	}
}

func (conn *Conn) OnClose() {

}

/****************************实现io标准的read******************************/

func (conn *Conn) Read(b []byte) (int, error) {
	return conn.Conn.Read(b)
}

func (conn *Conn) Write(b []byte) {
	conn.ChanMut.Lock()
	defer conn.ChanMut.Unlock()
	if !conn.Online || b == nil {
		return
	}
	if len(conn.ChanSend) >= cap(conn.ChanSend) {
		conn.doDestroy()
		return
	}
	conn.ChanSend <- b
}

func (conn *Conn) doDestroy() {
	// 强制关闭socket
	conn.Conn.(*net.TCPConn).SetLinger(0)
	conn.Conn.Close()
	close(conn.ChanSend)
	// 更新标记
	if conn.Online {
		conn.Online = false
	}
}

/****************************实现了iconn接口******************************/

func (conn *Conn) ReadMsg() ([]byte, error) {
	return conn.MsgParser.Read(conn)
}

func (conn *Conn) WriteMsg(args ...[]byte) error {
	return conn.MsgParser.Write(conn, args...)
}

func (conn *Conn) LocalAddr() net.Addr {
	return conn.Conn.LocalAddr()
}

func (conn *Conn) RemoteAddr() net.Addr {
	return conn.Conn.RemoteAddr()
}

func (conn *Conn) Close() {

}

func (conn *Conn) Destroy() {

}
