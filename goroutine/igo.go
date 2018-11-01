package goroutine

import (
	"time"
)

type Igo interface {
	Init()                                                          // 初始化
	Name() string                                                   // 设置进程别名
	Handler(msg string, args []interface{}, ret chan []interface{}) // 执行体
	SetTimer() time.Duration                                        // 绑定定时器
	TimerWork()                                                     // 定时器回调
	Close()                                                         // 强制关闭                                                      // 关闭
}
