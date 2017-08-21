package goroutine

type Igo interface {
	initGo()                             // 初始化
	handler([]interface{}) []interface{} // 执行体
	closeGo()                            // 关闭
}
