package entity

import (
	"fmt"

	"github.com/gfandada/gserver/util"
)

type Entity struct {
	Id        EntityId    // id标识
	Desc      *EntityDesc // 特征描述
	I         Ientity     // entity装载器
	Space     *Space      // 属于哪个space
	Destroyed bool        // 销毁标记:true-已销毁
	Client    *GameClient // 网络对象
	aoi       aoi         // aoi描述数据
	Att       *EntityAtt  // 属性
}

func NewEntity(desc *EntityDesc) *Entity {
	e := &Entity{
		Id:        EntityId(util.NewV4().String()),
		Desc:      desc,
		Destroyed: false,
		Client:    new(GameClient),
		Att:       new(EntityAtt),
	}
	initAOI(&e.aoi)
	return e
}

// 绑定Ientity
func (e *Entity) BindIentity(ientity Ientity) {
	e.I = ientity
}

// 字符串
func (e *Entity) String() string {
	return fmt.Sprintf("Entity<%s:%s>", e.Desc.Name, e.Id)
}

// 是否具有aoi
func (e *Entity) IsUseAOI() bool {
	return e.Desc.UseAOI
}

// 获取当前位置
func (e *Entity) GetPosition() Vector3 {
	return e.aoi.pos
}

func (e *Entity) init(entityID EntityId, entityInstance Ientity, desc *EntityDesc) {
	e.Id = entityID
	e.I = entityInstance
	e.Desc = desc
	initAOI(&e.aoi)
	e.I.OnInit()
}

/************************************AOI交互***********************************/

// 将指定Entity添加为邻居
func (e *Entity) interest(other *Entity) {
	e.aoi.interest(other)
}

// 将指定Entity从邻居中移除
func (e *Entity) uninterest(other *Entity) {
	e.aoi.uninterest(other)
}

// 获取邻居列表
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

/*********************************与space交互********************************/

// 进入space
func (e *Entity) EnterSpace(spaceID SpaceId, pos Vector3) {
	if e.Space != nil {
		return
	}
	space := GetSpace(spaceID)
	if space != nil {
		space.enter(e, pos, false)
	} else {
		fmt.Println("entity EnterSpace space is nil")
	}
}

// 离开space
func (e *Entity) LeaveSpace() {
	if e.Space == nil {
		return
	}
	e.Space = GetSpace(e.Space.Id)
	if e.Space != nil {
		e.Space.leave(e)
	} else {
		fmt.Println("entity LeaveSpace space is nil")
	}
}

// 在space中移动
func (e *Entity) MoveSpace(pos Vector3) {
	if e.Space == nil {
		return
	}
	e.Space = GetSpace(e.Space.Id)
	if e.Space != nil {
		e.Space.move(e, pos)
	} else {
		fmt.Println("entity MoveSpace space is nil")
	}
}

/*********************************与client交互********************************/

// 绑定clientid
// FIXME 绑定适当的client可以让entity具有网络传输能力
func (e *Entity) BindGameClient(clientid int32) {
	e.Client.clientid = clientid
}

// 异步消息
func (e *Entity) Post(users []*Entity, msg interface{}) {
	defer func() {
		if r := recover(); r != nil {
			// ...
			fmt.Println("post error ", r)
		}
	}()
	go e.Client.Post(users, msg)
}

/*********************************与att属性交互********************************/

func (e *Entity) Decrease(key string, value float32) float32 {
	return e.Att.decrease(key, value)
}

func (e *Entity) Increase(key string, value float32) float32 {
	return e.Att.increase(key, value)
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
	return false
}

func (e *Entity) Flag() int32 {
	return 0
}
