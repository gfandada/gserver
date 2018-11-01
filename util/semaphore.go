/*
    通用的限流器（基于线程并发数）
 */
package util

import (
	"context"
	"sync"
	"time"
)

type Semaphore struct {
	permits int           // 总的许可数量
	avail   int           // 当前可用的许可数量
	channel chan struct{} // 内部同步通道
	aMutex  sync.Mutex    // 用来同步外部的申请请求
	rMutex  sync.Mutex    // 用来同步外部的释放请求
	pMutex  sync.RWMutex  // 用来同步avail的数量
}

// 构建一个新的并发限流器
func New(permits int) *Semaphore {
	if permits < 1 {
		panic("Invalid number of permits. Less than 1")
	}
	// 填充内部channel
	channel := make(chan struct{}, permits)
	for i := 0; i < permits; i++ {
		channel <- struct{}{}
	}
	return &Semaphore{
		permits: permits,
		avail:   permits,
		channel: channel,
	}
}

// 获取许可
// 如果不可用，会一直block直到可用
func (s *Semaphore) Acquire() {
	s.aMutex.Lock()
	defer s.aMutex.Unlock()
	s.pMutex.Lock()
	s.avail--
	s.pMutex.Unlock()
	<-s.channel
}

// 获取许可，和Acquire不同的是，可以一次获取n个许可
// 如果不可用，会一直block直到可用
func (s *Semaphore) AcquireMany(n int) {
	if n > s.permits {
		n = s.permits
	}
	for ; n > 0; n-- {
		s.Acquire()
	}
}

// 获取许可，和AcquireMany不同的是，增加了context.Context，方便上下文同步
// 使用场景父子进程上下文
// 返回true则获取成功，false则失败
func (s *Semaphore) AcquireContext(ctx context.Context, n int) bool {
	acquired := make(chan struct{}, 1)
	reverse := make(chan bool, 1)
	go func() {
		s.AcquireMany(n)
		acquired <- struct{}{}
		if <-reverse {
			s.ReleaseMany(n)
		}
		close(acquired)
		close(reverse)
	}()

	select {
	case <-ctx.Done():
		reverse <- true
		return false
	case <-acquired:
		reverse <- false
		return true
	}
}

// 获取许可，和AcquireMany不同的是，不会一直block，增加了超时时间d
// 返回true则获取成功，false则超时失败
func (s *Semaphore) AcquireWithin(n int, d time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	return s.AcquireContext(ctx, n)
}

// 释放一个权限
// 如果不可用，会一直block直到可用
func (s *Semaphore) Release() {
	s.rMutex.Lock()
	defer s.rMutex.Unlock()
	s.channel <- struct{}{}
	s.pMutex.Lock()
	s.avail++
	s.pMutex.Unlock()
}

// 释放多个权限
// 如果不可用，会一直block直到可用
func (s *Semaphore) ReleaseMany(n int) {
	if n > s.permits {
		n = s.permits
	}
	for ; n > 0; n-- {
		s.Release()
	}
}

// 获取当前可用的许可数量
func (s *Semaphore) AvailablePermits() int {
	s.pMutex.RLock()
	defer s.pMutex.RUnlock()

	if s.avail < 0 {
		return 0
	}
	return s.avail
}

// 获取当前剩余的全部许可
func (s *Semaphore) DrainPermits() int {
	n := s.AvailablePermits()
	if n > 0 {
		s.AcquireMany(n)
	}
	return n
}
