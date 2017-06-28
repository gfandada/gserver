// 封装了连接池的操作
package network

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WsConnMap map[*websocket.Conn]struct{}

var wsconnPool WsConnMap
var wsconnMapMux sync.Mutex

// 初始化
func InitWsPool() {
	wsconnPool = make(WsConnMap)
}

// 关闭
func CloseWsPool() {
	wsconnMapMux.Lock()
	defer wsconnMapMux.Unlock()
	for conn := range wsconnPool {
		conn.Close()
	}
	wsconnPool = nil
}

// 添加一个新的连接
func AddWsConn(conn *websocket.Conn, maxNum int) bool {
	wsconnMapMux.Lock()
	defer wsconnMapMux.Unlock()
	// 检查是否达到连接池上线
	if wsconnPool == nil || len(wsconnPool) >= maxNum {
		return false
	}
	wsconnPool[conn] = struct{}{}
	return true
}

// 删除一个失效的连接
func DeleteConn(conn *websocket.Conn) {
	wsconnMapMux.Lock()
	defer wsconnMapMux.Unlock()
	delete(wsconnPool, conn)
}
