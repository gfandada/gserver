// 全局session数据管理
package network

import (
	"sync"
)

// 用户id-用户数据
type SessionMap map[uint32]map[string]interface{}

var Session SessionMap
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
	Session = make(SessionMap)
}

func OptSession(userId uint32, opt int, data []SessionData) {
	sessionMapMux.Lock()
	defer sessionMapMux.Unlock()
	switch opt {
	case Update:
		updateS(userId, data)
	case Delete:
		deleteS(userId, data)
	case Find:
		findS(userId, data)
	}
}

func findS(userId uint32, dataNew []SessionData) []*SessionData {
	data, ok := Session[userId]
	if !ok {
		return nil
	}
	ret := make([]*SessionData, len(dataNew))
	for index, value := range dataNew {
		ret[index] = &SessionData{
			Key:   value.Key,
			Value: data[value.Key],
		}
	}
	return ret
}

func updateS(userId uint32, dataNew []SessionData) {
	data, ok := Session[userId]
	if !ok {
		mapData := make(map[string]interface{})
		for _, value := range dataNew {
			mapData[value.Key] = value.Value
		}
		Session[userId] = mapData
		return
	}
	for _, value := range dataNew {
		data[value.Key] = value.Value
	}
}

func deleteS(userId uint32, dataNew []SessionData) {
	data, ok := Session[userId]
	if !ok {
		return
	}
	for _, value := range dataNew {
		delete(data, value.Key)
	}
}
