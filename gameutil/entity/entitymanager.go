package entity

import (
	"math/rand"
	"reflect"
)

var (
	_registeredEntityTypes map[string]*EntityTypeDesc // key是Entity.TypeName
	_entityManager         *EntityManager
)

type EntityIDSet map[EntityID]struct{}

type EntityManager struct {
	entities           map[EntityID]*Entity     // entity容器
	ownerOfClient      map[*GameClient]EntityID // 网络层与逻辑层的映射
	registeredServices map[string]EntityIDSet   // 服务容器
}

func init() {
	_entityManager = newEntityManager()
	_registeredEntityTypes = make(map[string]*EntityTypeDesc)
}

func newEntityManager() *EntityManager {
	return &EntityManager{
		entities:           map[EntityID]*Entity{},
		ownerOfClient:      map[*GameClient]EntityID{},
		registeredServices: map[string]EntityIDSet{},
	}
}

func (em *EntityManager) del(entityID EntityID) {
	delete(em.entities, entityID)
}

func (em *EntityManager) get(entityID EntityID) *Entity {
	return em.entities[entityID]
}

// 随机获取一个指定Service的提供商
func (em *EntityManager) chooseServiceProvider(serviceName string) EntityID {
	eids, ok := em.registeredServices[serviceName]
	if !ok {
		return ""
	}
	r := rand.Intn(len(eids))
	for eid := range eids {
		if r == 0 {
			return eid
		}
		r -= 1
	}
	return ""
}

// 清理client-entity映射关系
func (em *EntityManager) onEntityLoseClient(client *GameClient) {
	delete(em.ownerOfClient, client)
}

// 指定的client-entity映射关系
func (em *EntityManager) onEntityGetClient(entityID EntityID, client *GameClient) {
	em.ownerOfClient[client] = entityID
}

// 注册Entity
// @params typeName:Entity类型名称 entityPtr:行为装载器 isPersistent:是否持久化 useAOI:是否具有aoi行为
// @return EntityTypeDesc引用
func RegisterEntity(typeName string, entityPtr Ientity, isPersistent bool, useAOI bool) *EntityTypeDesc {
	if _, ok := _registeredEntityTypes[typeName]; ok {
		return nil
	}
	entityVal := reflect.Indirect(reflect.ValueOf(entityPtr))
	entityType := entityVal.Type()
	rpcDescs := rpcDescMap{}
	entityTypeDesc := &EntityTypeDesc{
		isPersistent:    isPersistent,
		useAOI:          useAOI,
		entityType:      entityType,
		rpcDescs:        rpcDescs,
		clientAttrs:     StringSet{},
		allClientAttrs:  StringSet{},
		persistentAttrs: StringSet{},
	}
	_registeredEntityTypes[typeName] = entityTypeDesc

	entityPtrType := reflect.PtrTo(entityType)
	numMethods := entityPtrType.NumMethod()
	for i := 0; i < numMethods; i++ {
		method := entityPtrType.Method(i)
		rpcDescs.visit(method)
	}
	return entityTypeDesc
}
