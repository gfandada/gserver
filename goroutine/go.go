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

	"github.com/gfandada/gserver/logger"
	"github.com/gfandada/gserver/util"
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

// 新建并运行一个本地进程
// @params igo：进程的装载器
// @return 携程id，错误描述
func Start(igo Igo) (pid uint64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	done := make(chan uint64, 1)
	timeD := igo.Timer()
	var loop func()
	var v *Goroutine
	if timeD <= 0 {
		loop = func() {
			pid, v = initGo(igo)
			done <- pid
			defer closeGo(pid, igo, v)
			for {
				select {
				case input := <-v.chanMsg:
					if input == nil {
						break
					}
					handler(igo, input)
				case <-v.chanControl:
					return
				}
			}
		}
	} else {
		timer := time.NewTimer(timeD)
		loop = func() {
			pid, v = initGo(igo)
			done <- pid
			defer closeGo(pid, igo, v)
			for {
				select {
				case input := <-v.chanMsg:
					if input == nil {
						break
					}
					handler(igo, input)
				case <-timer.C:
					//timer = time.After(igo.Timer())
					timer.Reset(igo.Timer())
					timer_work(igo)
				case <-v.chanControl:
					return
				}
			}
		}
	}
	go loop()
	pid = <-done
	if pid != 0 {
		return pid, nil
	}
	return 0, errors.New("create goroutine failed")
}

// 通过id停止指定进程
// @params id: 进程id
func StopById(id uint64) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	v := QueryById(id)
	if v != nil {
		v.chanControl <- struct{}{}
	}
	return
}

// 通过别名停止指定进程
// @params name: 进程别名
func StopByName(name string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	v := QueryByName(name)
	if v != nil {
		v.chanControl <- struct{}{}
	}
	return
}

// 同步请求
// @params pid:进程id msg:请求消息 args:消息参数 timeout:超时时间（秒）
// @return 请求结果，错误描述
func Call(pid uint64, msg string, args []interface{}, timeout int) ([]interface{}, error) {
	var err error
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
	v := QueryById(pid)
	if v == nil {
		err = errors.New("goroutine is not exist")
		return nil, err
	}
	v.chanMsg <- msgS
	select {
	case ret := <-msgS.chanRecv:
		return ret, err
	case <-time.After(time.Duration(timeout) * time.Second):
		return nil, errors.New("call time out")
	}
}

// 异步请求
// @params pid:进程id msg:请求消息 args:消息参数
// @return 错误描述
func Cast(pid uint64, msg string, args []interface{}) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	v := QueryById(pid)
	if v == nil {
		return errors.New("goroutine is not exist")
	}
	select {
	case v.chanMsg <- &message{
		msg:  msg,
		args: args,
	}:
	default: // 默认丢包
	}
	return nil
}

// 异步请求
// @params name:进程名称 msg:请求消息 args:消息参数
// @return 错误描述
func CastByName(name string, msg string, args []interface{}) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	v := QueryByName(name)
	if v == nil {
		return errors.New("goroutine is not exist")
	}
	select {
	case v.chanMsg <- &message{
		msg:  msg,
		args: args,
	}:
	default: // 默认丢包
	}
	return nil
}

// 获取排队消息数量
// @params pid:进程id
// @return 排队消息数量
func Pending(pid uint64) int {
	v := QueryById(pid)
	if v == nil {
		return 0
	}
	return len(v.chanMsg)
}

// 检查进程是否存活
// @params pid:进程id
// @return true|false
func IsAlive(pid uint64) bool {
	v := QueryById(pid)
	if v == nil {
		return false
	}
	return true
}

func initGo(igo Igo) (uint64, *Goroutine) {
	id := util.GetPid()
	v := &Goroutine{
		chanMsg:     make(chan *message, 10000),
		chanControl: make(chan struct{}, 1),
	}
	Register(id, igo.Name(), v)
	igo.InitGo()
	return id, v
}

func closeGo(id uint64, igo Igo, v *Goroutine) {
	Unregister(id, igo.Name())
	close(v.chanControl)
	close(v.chanMsg)
	igo.CloseGo()
}

func handler(igo Igo, input *message) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("goroutine handler panic input %v error: %v", input, r)
		}
	}()
	igo.Handler(input.msg, input.args, input.chanRecv)
}

func timer_work(igo Igo) {
	igo.Timer_work()
}
