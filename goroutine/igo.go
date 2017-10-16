package goroutine

import (
	"time"
)

type Igo interface {
	Name() string                                                   // 设置进程别名
	Timer() time.Duration                                           // 绑定定时器
	Timer_work()                                                    // 定时器回调
	InitGo()                                                        // 初始化
	Handler(msg string, args []interface{}, ret chan []interface{}) // 执行体
	CloseGo()                                                       // 关闭
}
