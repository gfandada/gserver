package gservices

import (
	"fmt"
	"testing"
	"time"

	"github.com/gfandada/gserver/logger"
)

type test struct {
	name string
	age  int
}

func Test_timer(t *testing.T) {
	logger.Start("./test.xml")
	server := NewLocalTimerServer()
	// 添加重复定时任务
	job1, ok := server.AddJobRepeat(time.Duration(time.Millisecond*20), 5, func(a []interface{}) {
		fmt.Println("我是重复定时任务", a[0].(*test).name)
	}, []interface{}{&test{name: "hello world", age: 20}})
	if !ok {
		t.Errorf("重复定时任务添加失败")
	}
	// 添加一次定时任务（适用于单次定时任务）
	job, ok := server.AddJobWithInterval(time.Duration(time.Second)*4, func(a []interface{}) {
		fmt.Println("我是单次定时任务")
	}, []interface{}{&test{name: "hello world", age: 20}})
	if !ok {
		t.Errorf("单次定时任务添加失败")
	}
	//	// 添加一次定时任务（适用于活动等有具体时间限制的任务）
	//	_, ok = server.AddJobWithDeadtime(time.Date(2017, 7, 4, 14, 55, 30, 0, time.Local), func() {
	//		fmt.Println("我是活动任务")
	//	})
	//	if !ok {
	//		t.Errorf("活动任务添加失败")
	//	}
	// 修改指定任务，job原来是4s后执行，现在改成3s
	ok = server.UpdateJobTimeout(job, time.Duration(time.Second)*3)
	if !ok {
		t.Errorf("修改指定任务失败")
	}
	time.Sleep(20e6)
	// 取消重复任务
	ok = server.DelJob(job1)
	if !ok {
		t.Errorf("取消重复任务失败")
	}
	// 添加一个过期定时任务
	_, ok = server.AddJobWithDeadtime(time.Date(2016, 7, 4, 14, 16, 30, 0, time.Local), func(a []interface{}) {
		fmt.Println("我是过期任务")
	}, []interface{}{&test{name: "hello world", age: 20}})
	if ok {
		t.Errorf("过期任务不应该添加成功")
	}
	time.Sleep(5e9)
}
