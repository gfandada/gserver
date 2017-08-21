package goroutine

import (
	"sync"
)

type GoMap struct {
	records map[uint64]interface{}
	sync.RWMutex
}

var (
	_default_registry GoMap
)

func init() {
	_default_registry.init()
}

func (s *GoMap) init() {
	s.records = make(map[uint64]interface{})
}

func (s *GoMap) register(id uint64, v interface{}) {
	s.Lock()
	defer s.Unlock()
	s.records[id] = v
}

func (s *GoMap) unregister(id uint64) {
	s.Lock()
	defer s.Unlock()
	delete(s.records, id)
}

func (s *GoMap) query(id uint64) interface{} {
	s.RLock()
	defer s.RUnlock()
	return s.records[id]
}

func (s *GoMap) count() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.records)
}

/******************************对外提供的接口*******************************/

// 注册一组值
func Register(id uint64, v interface{}) {
	_default_registry.register(id, v)
}

// 反注册
func Unregister(id uint64) {
	_default_registry.unregister(id)
}

// 查询
func Query(id uint64) interface{} {
	return _default_registry.query(id)
}

// 统计计数
func Count() int {
	return _default_registry.count()
}
