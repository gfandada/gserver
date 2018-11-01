package util

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var done = make(chan int, 3)
var g = &sync.WaitGroup{}

func Test(t *testing.T) {
	l := 7
	s := New(10)
	// 并发7个线程去处理
	for i := 0; i < l; i++ {
		g.Add(1)
		go aq(s, i+1)
	}
	// 第8个线程
	go func() {
		g.Add(1)
		// 获取5个锁
		if s.AcquireWithin(5, time.Second*3) {
			s.ReleaseMany(5)
			fmt.Println("8号获取并释放了5个许可")
		} else {
			fmt.Println("8号获取5个许可，超时了")
		}
		g.Done()
	}()
	fmt.Println("主线程开始wait")
	g.Wait()
	fmt.Println("主线程完成wait")

	// 这个时候肯定还有10个许可
	if n := s.DrainPermits(); n != 10 {
		t.Fail()
		s.ReleaseMany(n)
	} else {
		s.ReleaseMany(10)
	}
}

// 1-7号线程
func aq(s *Semaphore, i int) {
	fmt.Println(i, "号等待获取许可", ",当前剩余的许可:", s.AvailablePermits())
	s.AcquireMany(i)
	fmt.Println(i, "号获取", i, "个许可", ",当前剩余的许可:", s.AvailablePermits())
	time.Sleep(time.Second * 3)
	s.ReleaseMany(i)
	fmt.Println(i, "号释放", i, "个许可", ",当前剩余的许可:", s.AvailablePermits())
	g.Done()
}
