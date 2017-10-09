package entity

import (
	"fmt"
)

// Space定义
type Space struct {
	Entity                  // space也是一种Entity
	entities EntitySet      // space中的entity容器
	Kind     int            // space类型
	I        Ispace         // space装载器
	aoiCalc  Iaoicalculator // aoi计算器
}

// 检查是不是空场景
func (space *Space) IsNil() bool {
	return space.Kind == 0
}

// entity离开space
func (space *Space) leave(entity *Entity) {
	if entity.Space != space {
		return
	}
	if space.IsNil() {
		return
	}
	// 和邻居相互取消关注
	for neighbor := range entity.aoi.neighbors {
		entity.uninterest(neighbor)
		neighbor.uninterest(entity)
	}
	space.aoiCalc.Leave(&entity.aoi)
	space.entities.Del(entity)
	entity.Space = nil
	space.I.OnEntityLeaveSpace(entity)
	entity.I.OnLeaveSpace(space)
	// TODO 通知client
}

// entity进入space
// TODO 数据恢复
// @params entity:对象 pos:位置数据 isRestore:数据恢复？(true-是)
func (space *Space) enter(entity *Entity, pos Vector3, isRestore bool) {
	if entity.Space != nil {
		return
	}
	if space.IsNil() || !entity.IsUseAOI() {
		return
	}
	entity.Space = space
	space.entities.Add(entity)
	space.aoiCalc.Enter(&entity.aoi, pos)
	if !isRestore {
		// TODO 通知client
		enter, _ := space.aoiCalc.Adjust(&entity.aoi)
		// 和第一次的新邻居互相关注
		for _, naoi := range enter {
			neighbor := naoi.getEntity()
			entity.interest(neighbor)
			neighbor.interest(entity)
		}
		space.I.OnEntityEnterSpace(entity)
		entity.I.OnEnterSpace()
	} else {
		enter, _ := space.aoiCalc.Adjust(&entity.aoi)
		for _, naoi := range enter {
			neighbor := naoi.getEntity()
			entity.aoi.interest(neighbor)
			neighbor.aoi.interest(entity)
		}
	}
}

// 移动
// @params entity:移动对象 newPos:目标位置
func (space *Space) move(entity *Entity, newPos Vector3) {
	if space.IsNil() {
		return
	}
	space.aoiCalc.Move(&entity.aoi, newPos)
	enter, leave := space.aoiCalc.Adjust(&entity.aoi)
	// 和失效的邻居相互取消关注
	for _, naoi := range leave {
		neighbor := naoi.getEntity()
		entity.uninterest(neighbor)
		neighbor.uninterest(entity)
	}
	// 和新的邻居相互关注
	for _, naoi := range enter {
		neighbor := naoi.getEntity()
		entity.interest(neighbor)
		neighbor.interest(entity)
	}
}

func (space *Space) OnInit() {
	space.entities = EntitySet{}
	space.I = space.Entity.I.(Ispace)
	space.aoiCalc = newXZListAOICalculator()
	space.I.OnSpaceInit()
}

func (space *Space) OnDestroy() {
	for e := range space.entities {
		e.Destroy()
	}
	_spaceManager.delSpace(space.ID)
}

func (space *Space) String() string {
	if space.Kind != 0 {
		return fmt.Sprintf("Space<%d|%s>", space.Kind, space.ID)
	}
	return "Space<nil>"
}

func (space *Space) OnCreated() {
	space.onSpaceCreated()
	space.I.OnSpaceCreated()
}

func (space *Space) onSpaceCreated() {
	_spaceManager.putSpace(space)
}

/*********************************实现Ispace接口*******************************/

func (space *Space) OnSpaceInit() {

}

func (space *Space) OnSpaceCreated() {

}

func (space *Space) OnSpaceDestroy() {

}

func (space *Space) OnEntityEnterSpace(entity *Entity) {

}

func (space *Space) OnEntityLeaveSpace(entity *Entity) {

}
