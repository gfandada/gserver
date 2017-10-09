package entity

import (
	"reflect"
)

var (
	_spaceManager *SpaceManager
	_spaceType    reflect.Type
)

type SpaceManager struct {
	spaces map[EntityID]*Space // space-id和space的映射关系
}

func init() {
	_spaceManager = newSpaceManager()
}

func newSpaceManager() *SpaceManager {
	return &SpaceManager{
		spaces: map[EntityID]*Space{},
	}
}

func (spmgr *SpaceManager) putSpace(space *Space) {
	spmgr.spaces[space.ID] = space
}

func (spmgr *SpaceManager) delSpace(id EntityID) {
	delete(spmgr.spaces, id)
}

func (spmgr *SpaceManager) getSpace(id EntityID) *Space {
	return spmgr.spaces[id]
}

// 注册场景(场景是一种特殊的entity)
// @params spacePtr:场景装载器
func RegisterSpace(spacePtr Ispace) {
	spaceVal := reflect.Indirect(reflect.ValueOf(spacePtr))
	_spaceType = spaceVal.Type()
	RegisterEntity(_SPACE_ENTITY_TYPE, spacePtr.(Ientity), false, false)
}
