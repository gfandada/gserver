package entity

type EntityDesc struct {
	Name       string // 类型名称
	UseAOI     bool   // 是否具有AOI，有-true
	Persistent bool   // 是否持久化，有-true
	Flag       int32  // 类型标识
}
