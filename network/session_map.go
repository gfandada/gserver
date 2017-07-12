// 全局session数据管理
package network

import (
	"errors"
	"sync"

	"github.com/gfandada/gserver/logger"
)

// 用户id-用户数据
type SessionMap map[uint64]*Agent

var SessionConn SessionMap
var sessionMapMux sync.Mutex

// 操作类型
const (
	Update = iota
	Delete
	Find
)

type SessionData struct {
	Key   string
	Value interface{}
}

func NewSessionMap() {
	SessionConn = make(SessionMap)
}

func FindSessionConn(userId uint64) (*Agent, error) {
	sessionMapMux.Lock()
	data, ok := SessionConn[userId]
	sessionMapMux.Unlock()
	if !ok {
		return nil, errors.New("no user")
	}
	return data, nil
}

func AddSessionConn(userId uint64, agent *Agent) {
	sessionMapMux.Lock()
	defer sessionMapMux.Unlock()
	SessionConn[userId] = agent
	logger.Debug("AddSessionConn userid %d agent %v", userId, agent)
}

func DeleteSessionConn(userId uint64) {
	sessionMapMux.Lock()
	defer sessionMapMux.Unlock()
	delete(SessionConn, userId)
	logger.Debug("DeleteSessionConn userid %d", userId)
}
