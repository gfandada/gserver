package goroutine

import (
	"time"
)

type Igo interface {
	name() string                                                   // 设置进程别名
	timer() time.Duration                                           // 绑定定时器
	timer_work()                                                    // 定时器回调
	initGo()                                                        // 初始化
	handler(msg string, args []interface{}, ret chan []interface{}) // 执行体
	closeGo()                                                       // 关闭
}
