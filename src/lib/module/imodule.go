// 定义了模块的接口
package module

type Imodule interface {
	OnInit()               // 模块的初始化
	OnDestroy()            // 模块的销毁
	Run(ChClose chan bool) // 以插件化的方式运行模块
}
