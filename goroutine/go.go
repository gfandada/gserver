/***
	可被管理的轻量级进程服务，作用域是进程内
	进程特性：
	1 通过name/id作为唯一标识[id是进程的唯一标识，name是进程别名]
	2 可以方便的向指定进程发送消息
	3 进程间的消息分两类：业务消息和控制消息
	4 业务能力由进程自定义，控制消息目前只支持停止stop
	5 进程间没有直接或者间接联系，暂未实现link/restart机制
***/
package goroutine

import (
	"errors"
	"fmt"
	"time"

	"github.com/gfandada/gserver/util"
	"kubernetes/pkg/kubelet/kubeletconfig/util/log"
)

const (
	defaultTimeOut    = 5   // 超时时间s
	defaultMaxPending = 1e4 // pending上限
)

const (
	genServerStartFailed  = "genserver start failed"
	genServerNotFound     = "genserver is not found"
	genServerHandlerPanic = "genserver handler panic"
	genServerTimerPanic   = "genserver timer panic"
	callTimeOut           = "call time out"
)

type Goroutine struct {
	chanMsg     chan *message // 消息异步通道
	chanControl chan struct{} // 内置的控制通道
}

type message struct {
	chanRecv chan []interface{} // 接受消息
	msg      string             // 消息标识
	args     []interface{}      // 消息参数
}

// 新建并运行一个本地服务进程
// @params igo：进程的装载器
// @return 进程id，错误描述
func Start(igo Igo) (pid string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			log.Errorf(genServerStartFailed+" %v", err)
		}
	}()
	done := make(chan string, 1)
	goroutine := NewGoroutine()
	if igo.SetTimer() <= 0 {
		goroutine.doWithNoTimer(igo, done)
	} else {
		goroutine.doWithTimer(igo, done)
	}
	pid = <-done
	return pid, nil
}

// link型服务进程
func StartLink() {

}

// 停止指定进程
// @params flag: 进程标示，可以是id也可以是别名
func Stop(flag string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	goroutine, errS := GetOneGoroutine(flag)
	if errS != nil {
		return errS
	}
	if goroutine != nil {
		goroutine.chanControl <- struct{}{}
	}
	return
}

// 同步请求
// @params flag:进程标示（id或者别名）msg:请求消息 args:消息参数 timeout:超时时间（秒）
// @return 请求结果，错误描述
func Call(flag string, msg string, args []interface{}, timeout int) (retU []interface{}, err error) {
	msgS := &message{
		chanRecv: make(chan []interface{}, 1),
		msg:      msg,
		args:     args,
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
		close(msgS.chanRecv)
	}()
	goroutine, errS := GetOneGoroutine(flag)
	if errS != nil {
		return nil, errS
	}
	goroutine.chanMsg <- msgS
	if timeout == 0 {
		timeout = defaultTimeOut
	}
	select {
	case ret := <-msgS.chanRecv:
		return ret, nil
	case <-time.After(time.Duration(timeout) * time.Second):
		return nil, errors.New(callTimeOut)
	}
}

// 异步请求
// @params flag:进程标示（id或者别名) msg:请求消息 args:消息参数
// @return 错误描述
func Cast(flag string, msg string, args []interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	goroutine, errS := GetOneGoroutine(flag)
	if errS != nil {
		return errS
	}
	select {
	case goroutine.chanMsg <- &message{
		msg:  msg,
		args: args,
	}:
	default: // 默认丢包
	}
	return nil
}

// 获取排队消息数量
// @params flag:进程标示（id或者别名)
// @return 排队消息数量
func Pending(flag string) int {
	goroutine, err := GetOneGoroutine(flag)
	if err != nil {
		return 0
	}
	return len(goroutine.chanMsg)
}

// 根据server标示获取server实体
func GetOneGoroutine(flag string) (goroutine *Goroutine, err error) {
	goroutine = QueryById(flag)
	if goroutine == nil {
		goroutine = QueryByName(flag)
		if goroutine == nil {
			err = errors.New(genServerNotFound)
			return
		}
	}
	return
}

// 检查进程是否存活
// @params flag:进程标示（id或者别名)
// @return true|false
func IsAlive(flag string) bool {
	_, err := GetOneGoroutine(flag)
	if err != nil {
		return false
	}
	return true
}

// 构建server结构
func NewGoroutine() *Goroutine {
	return &Goroutine{}
}

// 使用loop定时器
func (goroutine *Goroutine) doWithTimer(igo Igo, done chan string) {
	timer := time.NewTimer(igo.SetTimer())
	loop := func() {
		pid := goroutine.init(igo)
		done <- pid
		defer goroutine.close(pid, igo)
		for {
			select {
			case input := <-goroutine.chanMsg:
				if input == nil {
					break
				}
				goroutine.handler(igo, input)
			case <-timer.C:
				timer.Reset(igo.SetTimer())
				goroutine.timerWork(igo)
			case <-goroutine.chanControl:
				return
			}
		}
	}
	go loop()
}

// 未使用定时器
func (goroutine *Goroutine) doWithNoTimer(igo Igo, done chan string) {
	loop := func() {
		pid := goroutine.init(igo)
		done <- pid
		defer goroutine.close(pid, igo)
		for {
			select {
			case input := <-goroutine.chanMsg:
				if input == nil {
					break
				}
				goroutine.handler(igo, input)
			case <-goroutine.chanControl:
				return
			}
		}
	}
	go loop()
}

func (goroutine *Goroutine) init(igo Igo) string {
	id := util.NewV4().ToString()
	// default 10000
	goroutine.chanMsg = make(chan *message, defaultMaxPending)
	goroutine.chanControl = make(chan struct{}, 1)
	Register(id, igo.Name(), goroutine)
	igo.Init()
	return id
}

func (goroutine *Goroutine) close(id string, igo Igo) {
	Unregister(id, igo.Name())
	close(goroutine.chanControl)
	close(goroutine.chanMsg)
	igo.Close()
}

func (goroutine *Goroutine) handler(igo Igo, input *message) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf(genServerHandlerPanic+" input %v error: %v", input, r)
		}
	}()
	igo.Handler(input.msg, input.args, input.chanRecv)
}

func (goroutine *Goroutine) timerWork(igo Igo) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf(genServerTimerPanic+" error: %v", r)
		}
	}()
	igo.TimerWork()
}
