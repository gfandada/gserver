// 定义了agent接口
package network

type Iagent interface {
	Run()     // 启动一个代理，执行相应的逻辑
	OnClose() // 关闭一个代理，清理资源
}
