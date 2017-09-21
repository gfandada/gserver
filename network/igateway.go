package network

import (
	"net"

	"github.com/gorilla/websocket"
)

type Igateway interface {
	WriteMsg(msg interface{})
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Start(*websocket.Conn)
	GetUserData() interface{}
	SetUserData(data interface{})
}
