package goroutine

type Igo interface {
	name() string                                                   // 设置进程别名
	initGo()                                                        // 初始化
	handler(msg string, args []interface{}, ret chan []interface{}) // 执行体
	closeGo()                                                       // 关闭
}
