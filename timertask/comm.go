package timertask

import "time"

type TimerTaskFunc func(args ...interface{})

// 定时任务管理器
type TimerTaskManager struct {
	curIndex  int                         // 当前下标
	slots     [3600]map[string]*TimerTask // 环形slots
	closed    chan struct{}               // 关闭通道
	taskClose chan struct{}               // 关闭任务
	timeClose chan struct{}               // 时间关闭
	startTime time.Time                   // 启动时间
}

// 任务
type TimerTask struct {
	cycleNum int           // 循环次数
	exec     TimerTaskFunc // 执行的函数
	params   []interface{} // 参数
}
