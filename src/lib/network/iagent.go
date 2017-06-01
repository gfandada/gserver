// 定义了agent接口
// 主要用于代理用户套接字操作
package network

type Agent interface {
	Run()     // 启动一个代理
	OnClose() // 关闭一个代理
}
