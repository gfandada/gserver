// entityid->*Entity
package entity

import (
	"sync"
)

var (
	_entitymanager *EntityManager
)

type EntityManager struct {
	entities map[EntityId]*Entity
	sync.RWMutex
}

func init() {
	_entitymanager = &EntityManager{
		entities: make(map[EntityId]*Entity),
	}
}

func (manager *EntityManager) put(entity *Entity) {
	manager.Lock()
	defer manager.Unlock()
	manager.entities[entity.Id] = entity
}

func (manager *EntityManager) del(entityid EntityId) {
	manager.Lock()
	defer manager.Unlock()
	delete(manager.entities, entityid)
}

func (manager *EntityManager) get(entityid EntityId) *Entity {
	manager.RLock()
	defer manager.RUnlock()
	return manager.entities[entityid]
}

func RegisterEntity(entity *Entity) {
	_entitymanager.put(entity)
}

func UnRegisterEntity(entity EntityId) {
	_entitymanager.del(entity)
}

func GetEntity(entity EntityId) *Entity {
	return _entitymanager.get(entity)
}
