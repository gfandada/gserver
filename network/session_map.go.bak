// 全局session数据管理
package network

import (
	"sync"
)

// 用户id-用户数据
type SessionMap map[uint64]map[string]interface{}

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

func OptSession(userId uint64, opt int, data []SessionData) []*SessionData {
	sessionMapMux.Lock()
	defer sessionMapMux.Unlock()
	switch opt {
	case Update:
		updateS(userId, data)
	case Delete:
		deleteS(userId, data)
	case Find:
		return findS(userId, data)
	}
	return nil
}

func findS(userId uint64, dataNew []SessionData) []*SessionData {
	if dataNew == nil || len(dataNew) == 0 {
		return nil
	}
	data, ok := Session[userId]
	if !ok {
		return nil
	}
	ret := make([]*SessionData, len(dataNew))
	for index, value := range dataNew {
		if _, ok := data[value.Key]; !ok {
			ret[index] = &SessionData{
				Key:   value.Key,
				Value: nil,
			}
			continue
		}
		ret[index] = &SessionData{
			Key:   value.Key,
			Value: data[value.Key],
		}
	}
	return ret
}

func updateS(userId uint64, dataNew []SessionData) {
	if dataNew == nil || len(dataNew) == 0 {
		return
	}
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
	return
}

func deleteS(userId uint64, dataNew []SessionData) {
	if dataNew == nil || len(dataNew) == 0 {
		return
	}
	data, ok := Session[userId]
	if !ok {
		return
	}
	for _, value := range dataNew {
		delete(data, value.Key)
	}
	return
}
