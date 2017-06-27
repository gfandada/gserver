// 封装了全局的timer
package gservices

// 客户端发送的message
type TimerMessage struct {
	F          interface{}   // 消息handler
	Args       []interface{} // 函数调用参数
	Date       interface{}   // 定时器执行的时间
	NoticeChan chan bool     // 通知通道
}

// 定时器服务器
type TimerServer struct {
	MonthlyChan chan *TimerMessage // 固定时间定时队列
	WeeklyChan  chan *TimerMessage // 每周时间定时队列
	DailyChan   chan *TimerMessage // 每日时间定时队列
}
