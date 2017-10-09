// entity可以表示player，monter，npc等逻辑体
// 不同的entity可以通过Ientity转载器表现不一样的特性
// 本模块内置了aoi等基础功能，是以上逻辑体的功能扩展的基础
package entity

import (
	"fmt"
	"reflect"
	"time"
	"unsafe"
)

/*********************************私有属性的访问********************************/

// 是否具有aoi行为
func (e *Entity) IsUseAOI() bool {
	return e.typeDesc.useAOI
}

// 获取当前位置
func (e *Entity) GetPosition() Vector3 {
	return e.aoi.pos
}

// 设置当前位置(移动)
func (e *Entity) SetPosition(pos Vector3) {
	e.Space.move(e, pos)
}

/*********************************实现Ientity接口********************************/

func (e *Entity) OnInit() {
}

func (e *Entity) OnCreated() {
}

func (e *Entity) OnDestroy() {
}

func (e *Entity) OnMigrateOut() {
}

func (e *Entity) OnMigrateIn() {
}

func (e *Entity) OnRestored() {
}

func (e *Entity) OnEnterSpace() {
}

func (e *Entity) OnLeaveSpace(space *Space) {
}

func (e *Entity) IsPersistent() bool {
	return e.typeDesc.isPersistent
}

/**********************************************************************/

func (e *Entity) String() string {
	return fmt.Sprintf("%s<%s>", e.TypeName, e.ID)
}

// 销毁entity
func (e *Entity) Destroy() {
	if e.destroyed {
		return
	}
	e.Space.leave(e)
	_entityManager.del(e.ID)
	e.destroyed = true
	e.client.destory()
}

func (e *Entity) IsDestroyed() bool {
	return e.destroyed
}

// TODO 持久化
func (e *Entity) Save() {
	if !e.I.IsPersistent() {
		return
	}
}

// 判断当前Entity是不是space
func (e *Entity) IsSpaceEntity() bool {
	return e.TypeName == _SPACE_ENTITY_TYPE
}

// 将Entity转换为space
func (e *Entity) ToSpace() *Space {
	if !e.IsSpaceEntity() {
		return nil
	}
	return (*Space)(unsafe.Pointer(e))
}

func (e *Entity) init(typeName string, entityID EntityID, entityInstance reflect.Value) {
	e.ID = entityID
	e.IV = entityInstance
	e.I = entityInstance.Interface().(Ientity)
	e.TypeName = typeName
	e.typeDesc = _registeredEntityTypes[typeName]
	// 初始化Entity-AOI
	initAOI(&e.aoi)
	e.I.OnInit()
}

// 将指定Entity添加为邻居
func (e *Entity) interest(other *Entity) {
	e.aoi.interest(other)
	// TODO client notify
}

// 将指定Entity从邻居中移除
func (e *Entity) uninterest(other *Entity) {
	e.aoi.uninterest(other)
	// TODO client notify
}

// 获取邻居
func (e *Entity) Neighbors() EntitySet {
	return e.aoi.neighbors
}

// 判断指定Entity是不是邻居
func (e *Entity) IsNeighbor(other *Entity) bool {
	return e.aoi.neighbors.Contains(other)
}

// 计算两个Entity之间的距离
func (e *Entity) DistanceTo(other *Entity) Coord {
	return e.aoi.pos.DistanceTo(other.aoi.pos)
}

// 以go的方式运行cb
func (e *Entity) Post(cb func()) {
	go cb()
}

// 调用另一个Entity的某个方法
// TODO 暂时只支持本地调用
func (e *Entity) Call(id EntityID, method string, args ...interface{}) {
	other := _entityManager.get(id)
	if other != nil {
		other.Post(func() {
			other.onCallFromLocal(method, args)
		})
	}
}

// TODO 调用指定的service
// TODO 暂时只支持本地调用
func (e *Entity) CallService(serviceName string, method string, args ...interface{}) {
	serviceEid := _entityManager.chooseServiceProvider(serviceName)
	e.Call(serviceEid, method, args)
}

// 本地调用Entity的方法
func (e *Entity) onCallFromLocal(methodName string, args []interface{}) {
	defer func() {
		if r := recover(); r != nil {
			// TODO
			fmt.Println("onCallFromLocal error ", r, " entity ", e)
		}
	}()
	rpcDesc := e.typeDesc.rpcDescs[methodName]
	if rpcDesc == nil {
		return
	}
	if rpcDesc.NumArgs < len(args) {
		fmt.Println("onCallFromLocal args error")
		return
	}
	methodType := rpcDesc.MethodType
	fmt.Println("methodType ", methodType)
	in := make([]reflect.Value, rpcDesc.NumArgs+1)
	in[0] = reflect.ValueOf(e.I)
	for i, v := range args {
		in[i+1] = reflect.ValueOf(v)
	}
	rpcDesc.Func.Call(in)
}

func (e *Entity) DeclareService(serviceName string) {
	e.declaredServices[serviceName] = struct{}{}
	// TODO notify client????
}

/********************************与client的交互*********************************/

// 获取绑定的client对象
func (e *Entity) GetClient() *GameClient {
	return e.client
}

// 设置client
func (e *Entity) SetClient(client *GameClient) {
	oldClient := e.client
	if oldClient.clientid == client.clientid {
		return
	}
	e.client = client
	if oldClient != nil {
		// 清理旧映射
		_entityManager.onEntityLoseClient(oldClient)
		for neighbor := range e.Neighbors() {
			// TODO 通知邻居entity的销毁
			fmt.Println("邻居", neighbor)
		}
		// TODO 通知e
	}
	if client != nil {
		// 保存新映射
		_entityManager.onEntityGetClient(e.ID, client)
		// TODO 通知client，entity的创建
		for neighbor := range e.Neighbors() {
			// TODO 通知邻居entity的创建
			fmt.Println("邻居", neighbor)
		}
	}
}

// 将e的client给other
func (e *Entity) GiveClientTo(other *Entity) {
	if e.client == nil {
		return
	}
	client := e.client
	e.SetClient(nil)
	other.SetClient(client)
}

// 让e和e的邻居执行f方法
// TODO
func (e *Entity) ForAllClients(f func(client *GameClient)) {
	if e.client != nil {
		f(e.client)
	}
	for neighbor := range e.Neighbors() {
		if neighbor.client != nil {
			f(neighbor.client)
		}
	}
}

/********************************与属性的交互*********************************/

/********************************与道具的交互*********************************/

/********************************与space的交互*********************************/

// 正在进入space
func (e *Entity) isEnteringSpace() bool {
	now := time.Now().UnixNano()
	return now < (e.enteringSpaceRequest.RequestTime + int64(ENTER_SPACE_REQUEST_TIMEOUT))
}

// 进入space
func (e *Entity) EnterSpace(spaceID EntityID, pos Vector3) {
	if e.isEnteringSpace() {
		e.I.OnEnterSpace()
		return
	}
	localSpace := _spaceManager.getSpace(spaceID)
	if localSpace != nil { // target space is local, just enter
		e.enterLocalSpace(localSpace, pos)
	} else {
		// TODO 迁移
	}
}

// 进入本地space
// 离开旧场景进入新场景
func (e *Entity) enterLocalSpace(space *Space, pos Vector3) {
	if space == e.Space {
		return
	}
	e.enteringSpaceRequest.SpaceID = space.ID
	e.enteringSpaceRequest.EnterPos = pos
	e.enteringSpaceRequest.RequestTime = time.Now().UnixNano()
	e.Post(func() {
		e.clearEnteringSpaceRequest()
		if space.IsDestroyed() {
			return
		}
		e.Space.leave(e)
		space.enter(e, pos, false)
	})
}

func (e *Entity) clearEnteringSpaceRequest() {
	e.enteringSpaceRequest.SpaceID = ""
	e.enteringSpaceRequest.EnterPos = Vector3{}
	e.enteringSpaceRequest.RequestTime = 0
}

/********************************同步函数************************************/

// 收集器：收集同步信息统一通知全局entity-client
func (e *Entity) CollectEntitySyncInfos() {

}
