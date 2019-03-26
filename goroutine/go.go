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
	server := NewServer()
	if igo.SetTimer() <= 0 {
		server.doWithNoTimer(igo, done)
	} else {
		server.doWithTimer(igo, done)
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
	server, errS := GetOneServer(flag)
	if errS != nil {
		return errS
	}
	if server != nil {
		server.chanControl <- struct{}{}
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
	server, errS := GetOneServer(flag)
	if errS != nil {
		return nil, errS
	}
	server.chanMsg <- msgS
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
	server, errS := GetOneServer(flag)
	if errS != nil {
		return errS
	}
	select {
	case server.chanMsg <- &message{
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
	server, err := GetOneServer(flag)
	if err != nil {
		return 0
	}
	return len(server.chanMsg)
}

// 根据server标示获取server实体
func GetOneServer(flag string) (server *Goroutine, err error) {
	server = QueryById(flag)
	if server == nil {
		server = QueryByName(flag)
		if server == nil {
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
	_, err := GetOneServer(flag)
	if err != nil {
		return false
	}
	return true
}

// 构建server结构
func NewServer() *Goroutine {
	return &Goroutine{}
}

// 使用loop定时器
func (s *Goroutine) doWithTimer(iServer Igo, done chan string) {
	timer := time.NewTimer(iServer.SetTimer())
	loop := func() {
		pid := s.init(iServer)
		done <- pid
		defer s.close(pid, iServer)
		for {
			select {
			case input := <-s.chanMsg:
				if input == nil {
					break
				}
				s.handler(iServer, input)
			case <-timer.C:
				timer.Reset(iServer.SetTimer())
				s.timerWork(iServer)
			case <-s.chanControl:
				return
			}
		}
	}
	go loop()
}

// 未使用定时器
func (s *Goroutine) doWithNoTimer(iServer Igo, done chan string) {
	loop := func() {
		pid := s.init(iServer)
		done <- pid
		defer s.close(pid, iServer)
		for {
			select {
			case input := <-s.chanMsg:
				if input == nil {
					break
				}
				s.handler(iServer, input)
			case <-s.chanControl:
				return
			}
		}
	}
	go loop()
}

func (s *Goroutine) init(iServer Igo) string {
	id := string(util.NewV4().Bytes())
	// default 10000
	s.chanMsg = make(chan *message, defaultMaxPending)
	s.chanControl = make(chan struct{}, 1)
	Register(id, iServer.Name(), s)
	iServer.Init()
	return id
}

func (s *Goroutine) close(id string, iServer Igo) {
	Unregister(id, iServer.Name())
	close(s.chanControl)
	close(s.chanMsg)
	iServer.Close()
}

func (s *Goroutine) handler(iServer Igo, input *message) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf(genServerHandlerPanic+" input %v error: %v", input, r)
		}
	}()
	iServer.Handler(input.msg, input.args, input.chanRecv)
}

func (s *Goroutine) timerWork(iServer Igo) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf(genServerTimerPanic+" error: %v", r)
		}
	}()
	iServer.TimerWork()
}
