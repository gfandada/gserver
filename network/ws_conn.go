package network

import (
	"fmt"
	"net"
	"time"

	"github.com/gfandada/gserver/logger"
	"github.com/gorilla/websocket"
)

type WsConn struct {
	conn *websocket.Conn
}

func (conn *WsConn) ReadMsg() ([]byte, error) {
	_, data, err := conn.conn.ReadMessage()
	return data, err
}

func (conn *WsConn) WriteMsg(arg []byte) error {
	return conn.conn.WriteMessage(websocket.BinaryMessage, arg)
}

func (conn *WsConn) LocalAddr() net.Addr {
	return conn.conn.LocalAddr()
}

func (conn *WsConn) RemoteAddr() net.Addr {
	return conn.conn.RemoteAddr()
}

func (conn *WsConn) SetReadDeadline(t time.Time) error {
	return conn.conn.SetReadDeadline(t)
}

func (conn *WsConn) SetWriteDeadline(t time.Time) error {
	return conn.conn.SetWriteDeadline(t)
}

func (conn *WsConn) Close() {
	logger.Debug(fmt.Sprintf("websocket close conn %v", conn.RemoteAddr()))
	conn.conn.Close()
}
