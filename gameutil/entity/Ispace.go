package entity

type Ispace interface {
	OnSpaceInit()                      // space初始化
	OnSpaceCreated()                   // space创建
	OnSpaceDestroy()                   // space销毁
	OnEntityEnterSpace(entity *Entity) // 当任意Entity进入space时
	OnEntityLeaveSpace(entity *Entity) // 当任意Entity离开space时
}
