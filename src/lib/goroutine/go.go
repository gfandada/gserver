// 封装了可被管理的携程
package goroutine

import (
	"container/list"
	"fmt"
	"sync"
)

type Go struct {
	ChanCb    chan func() // 基于函数的chan
	PendingGo int         // 排队的数量
}

type LinearGo struct {
	Fun      func() // 执行的方法
	CallBack func() // 执行的回调
}

// 上下文环境定义
type LinearContext struct {
	Gop            *Go        // Go
	LinearGop      *list.List // linearGo的容器
	MutexLinearGo  sync.Mutex // 互斥的访问linearGo容器
	MutexExecution sync.Mutex // 互斥的执行方法
}

// 新建一个gogoroutine go
func NewGo(length int) *Go {
	goS := new(Go)
	goS.ChanCb = make(chan func(), length)
	return goS
}

// FIXME 以携程的方式异步运行函数fun，同时注册其回调
// 非携程安全
func (g *Go) Go(fun func(), callback func()) {
	if g == nil || g.ChanCb == nil {
		fmt.Println("Go is nil")
		return
	}
	fmt.Println("goroutine go")
	// FIXME 为了保证安全，不在g.ChanCb <- callback处自增
	g.PendingGo++
	go func() {
		defer func() {
			g.ChanCb <- callback
			if r := recover(); r != nil {
				fmt.Println("goroutine go panic:", r)
			}
		}()
		fun()
	}()
}

// FIXME 以携程的方式异步运行函数fun，同时注册其回调
// 携程安全
func (linearContext *LinearContext) Go(fun func(), callback func()) {
	if linearContext == nil || linearContext.Gop == nil {
		fmt.Println("Go is nil")
		return
	}
	// 锁代码
	linearContext.MutexLinearGo.Lock()
	linearContext.LinearGop.PushBack(&LinearGo{fun, callback})
	linearContext.MutexLinearGo.Unlock()
	linearContext.Gop.PendingGo++
	go func() {
		// 锁方法
		linearContext.MutexExecution.Lock()
		defer linearContext.MutexExecution.Unlock()
		linearContext.MutexLinearGo.Lock()
		linearGop := linearContext.LinearGop.Remove(linearContext.LinearGop.Front()).(*LinearGo)
		linearContext.MutexLinearGo.Unlock()
		defer func() {
			linearContext.Gop.ChanCb <- linearGop.CallBack
			if r := recover(); r != nil {
				fmt.Println("linearContext Go error:", r)
			}
		}()
		linearGop.Fun()
	}()
}

// FIXME 以函数调用的方式运行回调函数
func (g *Go) CallBack(callback func()) {
	if g == nil || g.ChanCb == nil || callback != nil {
		fmt.Println("params is nil")
		return
	}
	fmt.Println("goroutine CallBack")
	defer func() {
		g.PendingGo--
		if r := recover(); r != nil {
			fmt.Println("goroutine CallBack:", r)
		}
	}()
	callback()
}

// 优雅的关闭
func (g *Go) Close() {
	if g == nil || g.ChanCb == nil || g.PendingGo <= 0 {
		return
	}
	for g.PendingGo > 0 {
		g.CallBack(<-g.ChanCb)
	}
}

// 是否空闲，true：空闲 false：不空闲
// FIXME 关于空闲的定义是完全取决于pending，而非负载
func (g *Go) Idle() bool {
	if g == nil {
		fmt.Println("g is nil")
	}
	return g.PendingGo == 0
}

// 基于Go,新建一个上下文环境
func (g *Go) NewLinearContext() *LinearContext {
	linearContext := new(LinearContext)
	linearContext.Gop = g
	linearContext.LinearGop = list.New()
	return linearContext
}
