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
	Start(*websocket.Conn, *Config)
	GetUserData() interface{}
	SetUserData(data interface{})
}
