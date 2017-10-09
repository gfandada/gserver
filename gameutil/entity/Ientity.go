package entity

// entity装载器接口
type Ientity interface {
	OnInit()                   // 初始化
	OnCreated()                // 创建
	OnDestroy()                // 销毁
	OnMigrateOut()             // space迁移
	OnMigrateIn()              // space迁移
	OnRestored()               // 恢复
	OnEnterSpace()             // 进入space
	OnLeaveSpace(space *Space) // 离开space
	IsPersistent() bool        // 是否需要数据持久化
}
