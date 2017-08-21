package goroutine

import (
	"errors"
	"time"

	"github.com/gfandada/gserver/util"
)

type Goroutine struct {
	ChanSend    chan []interface{} // 发送通道
	ChanRecv    chan []interface{} // 接收通道
	ChanControl chan struct{}      // 内置的控制通道
	Pending     int                // 排队的消息数量
}

// 新建并运行一个本地携程
// @params igo：携程的装载器  async：true异步 false同步
// @return 携程id
func Start(igo Igo, async bool) (uint64, error) {
	var pid uint64
	loop := func() {
		pid1, v := initGo(igo)
		pid = pid1
		defer closeGo(pid, igo, v)
		for {
			select {
			case input := <-v.ChanSend:
				if input == nil {
					v.Pending--
					break
				}
				output := igo.handler(input)
				if !async {
					v.ChanRecv <- output
				}
				v.Pending--
			// 暂不开放更多的携程控制权限
			case <-v.ChanControl:
				return
			}
		}
	}
	go loop()
	time.Sleep(time.Duration(50) * time.Microsecond)
	if pid != 0 {
		return pid, nil
	}
	return 0, errors.New("create goroutine failed")
}

// 停止指定携程
// id: 携程的pid
func Stop(id uint64) {
	defer func() {
		recover()
	}()
	v := Query(id)
	if v != nil {
		v.(*Goroutine).ChanControl <- struct{}{}
	}
}

// 同步请求
// input:请求参数  timeout:超时时间（秒）
func Call(pid uint64, input []interface{}, timeout int) ([]interface{}, error) {
	defer func() {
		recover()
	}()
	v := Query(pid)
	if v == nil {
		return nil, errors.New("goroutine is not exist")
	}
	g := v.(*Goroutine)
	g.ChanSend <- input
	select {
	case ret := <-g.ChanRecv:
		return ret, nil
	case <-time.After(time.Duration(timeout) * time.Second):
		return nil, errors.New("time out")
	}
}

// 异步请求
func Cast(pid uint64, input []interface{}) error {
	defer func() {
		recover()
	}()
	v := Query(pid)
	if v == nil {
		return errors.New("goroutine is not exist")
	}
	g := v.(*Goroutine)
	g.ChanSend <- input
	return nil
}

func initGo(igo Igo) (uint64, *Goroutine) {
	id := util.GetPid()
	v := &Goroutine{
		ChanSend:    make(chan []interface{}, 1),
		ChanRecv:    make(chan []interface{}, 1),
		ChanControl: make(chan struct{}, 1),
	}
	Register(id, v)
	igo.initGo()
	return id, v
}

func closeGo(id uint64, igo Igo, v *Goroutine) {
	Unregister(id)
	close(v.ChanControl)
	close(v.ChanRecv)
	close(v.ChanSend)
	igo.closeGo()
}
