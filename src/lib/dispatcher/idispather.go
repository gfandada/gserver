package dispather

type Idispather interface {
	Load()   // 装载配置
	UnLoad() // 卸载配置
}
