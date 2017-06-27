package gateway

import (
	"net"
)

type Igateway interface {
	WriteMsg(msg interface{})
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	GetUserData() interface{}
	SetUserData(data interface{})
}
