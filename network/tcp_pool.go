// 封装了连接池的操作
package network

import (
	"fmt"
	"net"
	"sync"

	"github.com/gfandada/gserver/logger"
)

type ConnMap map[net.Conn]struct{}

var connPool ConnMap
var connMapMux sync.Mutex

// 初始化
func Init() {
	connPool = make(ConnMap)
}

// 关闭
func Close() {
	connMapMux.Lock()
	defer connMapMux.Unlock()
	for conn := range connPool {
		conn.Close()
	}
	connPool = nil
}

// 添加一个新的连接
func (server *TcpServer) AddConn(conn net.Conn, maxNum int) bool {
	connMapMux.Lock()
	defer connMapMux.Unlock()
	if len(connPool) >= maxNum {
		logger.Warning(fmt.Sprintf("The number of connections has reached the upper limit, %d", maxNum))
		return false
	}
	connPool[conn] = struct{}{}
	return true
}

// 删除一个失效的连接
func (server *TcpServer) DeleteConn(conn net.Conn) {
	connMapMux.Lock()
	defer connMapMux.Unlock()
	delete(connPool, conn)
}
