// 定义了客户端套接字的接口
package network

import (
	"net"
)

type Iconn interface {
	ReadMsg() ([]byte, error)      // 从客户端的套接字中获取数据
	WriteMsg(args ...[]byte) error // 往客户端的套接字中写入数据
	LocalAddr() net.Addr           // 获取客户端连接的本地地址
	RemoteAddr() net.Addr          // 获取客户端连接的远端地址
	Close()                        // 关闭客户端连接
	Destroy()                      // 销毁客户端连接
}
