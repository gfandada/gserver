package entity

import (
	"sync"
)

var (
	_spacemanager *SpaceManager
)

type SpaceManager struct {
	spaces map[SpaceId]*Space
	sync.RWMutex
}

func init() {
	_spacemanager = &SpaceManager{
		spaces: make(map[SpaceId]*Space),
	}
}

func (manager *SpaceManager) put(space *Space) {
	manager.Lock()
	defer manager.Unlock()
	manager.spaces[space.Id] = space
}

func (manager *SpaceManager) del(spaceid SpaceId) {
	manager.Lock()
	defer manager.Unlock()
	delete(manager.spaces, spaceid)
}

func (manager *SpaceManager) get(spaceid SpaceId) *Space {
	manager.RLock()
	defer manager.RUnlock()
	return manager.spaces[spaceid]
}

func RegisterSpace(space *Space) {
	_spacemanager.put(space)
}

func UnRegisterSpace(space SpaceId) {
	_spacemanager.del(space)
}

func GetSpace(space SpaceId) *Space {
	return _spacemanager.get(space)
}
