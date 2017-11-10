// entityid->*Entity
package entity

import (
	"sync"
)

var (
	_entitymanager *EntityManager
)

type EntityManager struct {
	entities sync.Map
}

func init() {
	_entitymanager = &EntityManager{}
}

func (manager *EntityManager) put(entity *Entity) {
	manager.entities.Store(entity.Id, entity)
}

func (manager *EntityManager) del(entityid EntityId) {
	manager.entities.Delete(entityid)
}

func (manager *EntityManager) get(entityid EntityId) *Entity {
	entity, _ := manager.entities.Load(entityid)
	return entity.(*Entity)
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
