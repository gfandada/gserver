package network

import (
	"net"
	"time"
)

type Iconn interface {
	ReadMsg() ([]byte, error)           // 读取
	WriteMsg(arg []byte) error          // 写入
	LocalAddr() net.Addr                // 本地地址
	RemoteAddr() net.Addr               // 远程地址
	SetReadDeadline(t time.Time) error  // 读超时
	SetWriteDeadline(t time.Time) error // 写超时
	Close()                             // 关闭
}
