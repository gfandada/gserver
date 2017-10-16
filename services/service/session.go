package service

import (
	"sync"

	"github.com/gfandada/gserver/network"
)

var (
	_sessionm SessionManger
)

type SessionManger struct {
	pool map[int32]*Session
	sync.RWMutex
}

type Session struct {
	MQ       chan network.Data_Frame // 返回给网关的异步消息
	Agent    *Agent                  // 代理器
	UserId   int32                   // 玩家ID
	Die      chan struct{}           // 会话关闭信号
	Flag     int32                   // 会话标记
	UserData map[string]interface{}  // 用户自定义sess数据
}

func init() {
	_sessionm.init()
}

func (s *SessionManger) init() {
	s.pool = make(map[int32]*Session)
}

func (s *SessionManger) add(id int32, v *Session) {
	s.Lock()
	defer s.Unlock()
	s.pool[id] = v
}

func (s *SessionManger) remove(id int32) {
	s.Lock()
	defer s.Unlock()
	delete(s.pool, id)
}

func (s *SessionManger) get(id int32) *Session {
	s.RLock()
	defer s.RUnlock()
	return s.pool[id]
}

func (s *SessionManger) count() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.pool)
}

/**************************自定义session(非线程安全)*****************************/

func (s *Session) AddData(key string, data interface{}) {
	s.UserData[key] = data
}

func (s *Session) DelData(key string) {
	delete(s.UserData, key)
}

func (s *Session) Get(key string) interface{} {
	return s.UserData[key]
}

func New() *Session {
	sess := new(Session)
	sess.Die = make(chan struct{})
	sess.MQ = make(chan network.Data_Frame, DEFAULT_CH_SIZE)
	return sess
}

func Add(id int32, v *Session) {
	_sessionm.add(id, v)
}

func Remove(id int32) {
	_sessionm.remove(id)
}

func Get(id int32) *Session {
	return _sessionm.get(id)
}

// for async ipc, not sync
func Send(id int32, msg network.RawMessage) {
	sess := Get(id)
	if sess == nil {
		return
	}
	sess.Agent.Send(msg)
}

func Count() int {
	return _sessionm.count()
}

func ForEachSend(msg network.RawMessage) {
	for userid := range _sessionm.pool {
		Send(userid, msg)
	}
}
