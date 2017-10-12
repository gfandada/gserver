// space统一维护aoi数据的变化
// 不可也没有必要并发调用space接口
package entity

import (
	"fmt"

	"github.com/gfandada/gserver/util"
)

type Space struct {
	Id       SpaceId        // 场景id
	Type     int            // space类型
	I        Ispace         // space装载器
	aoiCalc  Iaoicalculator // aoi计算器
	entities EntitySet      // space中的entity容器
}

func NewSpace(spaceType int, ispace Ispace) *Space {
	return &Space{
		Id:       SpaceId(util.NewV4().String()),
		Type:     spaceType,
		I:        ispace,
		aoiCalc:  newXZListAOICalculator(),
		entities: EntitySet{},
	}
}

func (space *Space) String() string {
	if space.Type != 0 {
		return fmt.Sprintf("Space<%d:%s>", space.Type, space.Id)
	}
	return "Space<nil>"
}

// 检查是不是空场景
func (space *Space) IsNil() bool {
	return space.Type == NIL_SPACE
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
