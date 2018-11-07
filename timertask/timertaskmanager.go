/**
通用的全局定时任务管理器
未实现任务持久化
*/
package timertask

import (
	"errors"
	"fmt"
	"time"
)

// 创建一个定时任务管理器
func NewTimerTaskManager() *TimerTaskManager {
	manager := &TimerTaskManager{
		curIndex:  0,
		closed:    make(chan struct{}),
		taskClose: make(chan struct{}),
		timeClose: make(chan struct{}),
		startTime: time.Now(),
	}
	for i := 0; i < 3600; i++ {
		manager.slots[i] = make(map[string]*TimerTask)
	}
	return manager
}

// 启动
func (manager *TimerTaskManager) Start() {
	go manager.taskLoop()
	go manager.timeLoop()
	select {
	case <-manager.closed:
		manager.taskClose <- struct{}{}
		manager.timeClose <- struct{}{}
	}
}

// 停止
func (manager *TimerTaskManager) Close() {
	manager.closed <- struct{}{}
}

// 循环执行当前slot中cycleNum=0的任务
func (manager *TimerTaskManager) taskLoop() {
	defer func() {
		fmt.Println("taskLoop exit")
	}()
	fmt.Println("taskLoop start")
	for {
		select {
		case <-manager.taskClose:
			return
		default:
			// 取出当前槽的任务
			tasks := manager.slots[manager.curIndex]
			if len(tasks) > 0 {
				// 遍历任务，判断任务循环次数等于0，则执行任务
				// 否则任务循环次数减1
				for k, v := range tasks {
					if v.cycleNum == 0 {
						// FIXME 限流请使用semaphore
						go v.exec(v.params...)
						delete(tasks, k)
					} else {
						v.cycleNum--
					}
				}
			}
		}
	}
}

// 定时去更新curIndex
func (manager *TimerTaskManager) timeLoop() {
	defer func() {
		fmt.Println("timeLoop exit")
	}()
	fmt.Println("timeLoop start")
	tick := time.NewTicker(time.Second)
	for {
		select {
		case <-manager.timeClose:
			return
		case <-tick.C:
			fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " timeLoop loop")
			// 判断当前下标，如果等于3599则重置为0，否则加1
			if manager.curIndex == 3599 {
				manager.curIndex = 0
			} else {
				manager.curIndex++
			}
		}
	}
}

// 添加任务
func (manager *TimerTaskManager) AddTask(t time.Time, key string, exec TimerTaskFunc, params []interface{}) error {
	if manager.startTime.After(t) {
		return errors.New("manager.startTime is after t")
	}
	// 当前时间与指定时间相差秒数
	subSecond := t.Unix() - manager.startTime.Unix()
	// 计算循环次数
	cycleNum := int(subSecond / 3600)
	// 计算任务所在的slots的下标
	ix := subSecond % 3600
	tasks := manager.slots[ix]
	if _, ok := tasks[key]; ok {
		return errors.New("slots has existed task " + key)
	}
	tasks[key] = &TimerTask{
		cycleNum: cycleNum,
		exec:     exec,
		params:   params,
	}
	return nil
}
