package network

import (
	"net"
	"sync"

	"github.com/gorilla/websocket"
)

// WsConn数据
type WsConn struct {
	Conn     *websocket.Conn // 客户端的连接
	ChanSend chan []byte     // 用来保存服务器需要发送给客户端的数据
	//	ChanRecv  chan []byte     //  用来保存接收到的客户端数据
	ChanMut   sync.Mutex     // 保证ChanSend携程安全
	Online    bool           // 是否在线的标记:true在线 false离线
	MsgParser *MessageParser // 消息解析器
}

// 初始化客户端conn
func InitWsConn(conn *websocket.Conn, pendingNum int, msgParser *MessageParser) *WsConn {
	wsConn := new(WsConn)
	wsConn.Conn = conn
	wsConn.ChanSend = make(chan []byte, pendingNum)
	wsConn.Online = true
	wsConn.MsgParser = msgParser
	go func() {
		for sendData := range wsConn.ChanSend {
			if sendData == nil {
				break
			}
			err := conn.WriteMessage(websocket.BinaryMessage, sendData)
			if err != nil {
				break
			}
		}
		wsConn.Close()
	}()
	return wsConn
}

func (wsConn *WsConn) doDestroy() {
	wsConn.Conn.UnderlyingConn().(*net.TCPConn).SetLinger(0)
	wsConn.Conn.Close()
	if wsConn.Online {
		close(wsConn.ChanSend)
		wsConn.Online = false
	}
}
func (wsConn *WsConn) Write(b []byte) {
	wsConn.ChanMut.Lock()
	defer wsConn.ChanMut.Unlock()
	if !wsConn.Online || b == nil {
		return
	}
	if len(wsConn.ChanSend) == cap(wsConn.ChanSend) {
		wsConn.doDestroy()
		return
	}
	wsConn.ChanSend <- b
}

/****************************实现了iconn接口******************************/

// 非携程安全
func (wsConn *WsConn) ReadMsg() ([]byte, error) {
	return wsConn.MsgParser.ReadWs(wsConn)
}

// 写操作
// 线程安全
func (wsConn *WsConn) WriteMsg(args ...[]byte) error {
	return wsConn.MsgParser.WriteWs(wsConn, args...)
}

func (wsConn *WsConn) LocalAddr() net.Addr {
	return wsConn.Conn.LocalAddr()
}

func (wsConn *WsConn) RemoteAddr() net.Addr {
	return wsConn.Conn.RemoteAddr()
}

func (wsConn *WsConn) Close() {
	if !wsConn.Online {
		return
	}
	wsConn.ChanMut.Lock()
	defer wsConn.ChanMut.Unlock()
	wsConn.Write(nil)
	wsConn.Online = false
}

func (wsConn *WsConn) Destroy() {
	wsConn.ChanMut.Lock()
	defer wsConn.ChanMut.Unlock()
	wsConn.doDestroy()
}
