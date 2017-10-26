package entity

import (
	"sync"
)

type EntityAtt struct {
	att map[string]float32
	sync.RWMutex
}

func (att *EntityAtt) add(key string, value float32) {
	att.Lock()
	defer att.Unlock()
	att.att[key] = value
}

func (att *EntityAtt) del(key string) {
	att.Lock()
	defer att.Unlock()
	delete(att.att, key)
}

func (att *EntityAtt) increase(key string, value float32) float32 {
	att.Lock()
	defer att.Unlock()
	v, ok := att.att[key]
	if !ok {
		att.att[key] = value
	} else {
		att.att[key] = v + value
	}
	return att.att[key]
}

func (att *EntityAtt) decrease(key string, value float32) float32 {
	att.Lock()
	defer att.Unlock()
	v, ok := att.att[key]
	if !ok {
		att.att[key] = float32(0)
	} else if v-value >= 0.001 {
		att.att[key] = v - value
	} else {
		att.att[key] = float32(0)
	}
	return att.att[key]
}

func (att *EntityAtt) get(key string) float32 {
	att.RLocker()
	defer att.RUnlock()
	return att.att[key]
}
