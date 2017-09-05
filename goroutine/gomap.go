package goroutine

import (
	"sync"
)

type GoMap struct {
	records map[uint64]*Goroutine // [id]*Goroutine
	names   map[string]uint64     // [name]id
	sync.RWMutex
}

var (
	_default_registry GoMap
)

func init() {
	_default_registry.init()
}

func (s *GoMap) init() {
	s.records = make(map[uint64]*Goroutine)
	s.names = make(map[string]uint64)
}

func (s *GoMap) registerById(id uint64, v *Goroutine) {
	s.Lock()
	defer s.Unlock()
	s.records[id] = v
}

func (s *GoMap) unregisterById(id uint64) {
	s.Lock()
	defer s.Unlock()
	delete(s.records, id)
}

func (s *GoMap) registerByName(id uint64, name string, v *Goroutine) {
	s.Lock()
	defer s.Unlock()
	s.records[id] = v
	s.names[name] = id
}

func (s *GoMap) unregisterByName(name string) {
	s.Lock()
	defer s.Unlock()
	delete(s.names, name)
	delete(s.records, s.names[name])
}

func (s *GoMap) queryById(id uint64) *Goroutine {
	s.RLock()
	defer s.RUnlock()
	return s.records[id]
}

func (s *GoMap) queryByName(name string) *Goroutine {
	s.RLock()
	defer s.RUnlock()
	id := s.names[name]
	return s.records[id]
}

func (s *GoMap) count() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.records)
}

/******************************对外提供的接口*******************************/

// 注册一组值
func Register(id uint64, name string, v *Goroutine) {
	if name == "" {
		_default_registry.registerById(id, v)
	} else {
		_default_registry.registerByName(id, name, v)
	}
}

// 反注册
func Unregister(id uint64, name string) {
	if name == "" {
		_default_registry.unregisterById(id)
	} else {
		_default_registry.unregisterByName(name)
	}
}

// 通过进程id查询
func QueryById(id uint64) *Goroutine {
	return _default_registry.queryById(id)
}

// 通过进程name查询
func QueryByName(name string) *Goroutine {
	return _default_registry.queryByName(name)
}

// 统计计数
func Count() int {
	return _default_registry.count()
}
