package timertask

import (
	"fmt"
	"testing"
	"time"
)

func Test_timerTask(t *testing.T) {
	dm := NewTimerTaskManager()
	dm.AddTask(time.Now().Add(time.Second*5), "test1", func(args ...interface{}) {
		fmt.Println(args...)
	}, []interface{}{1, 2, 3})
	dm.AddTask(time.Now().Add(time.Second*5), "test2", func(args ...interface{}) {
		fmt.Println(args...)
	}, []interface{}{4, 5, 6})
	dm.AddTask(time.Now().Add(time.Second*10), "test3", func(args ...interface{}) {
		fmt.Println(args...)
	}, []interface{}{"hello", "world", "test"})
	dm.AddTask(time.Now().Add(time.Second*15), "test4", func(args ...interface{}) {
		sum := 0
		for arg := range args {
			sum += arg
		}
		fmt.Printf("sum=%v\n", sum)
	}, []interface{}{1, 2, 3})

	// 20秒后关闭
	time.AfterFunc(time.Second*20, func() {
		dm.Close()
	})
	dm.Start()
}
