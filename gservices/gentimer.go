package gservices

import (
	"math"
	"sync/atomic"
	"time"

	"github.com/HuKeping/rbtree"
	"github.com/gfandada/gserver/logger"
)

// 定时器初始值
const DEFAULT = time.Duration(math.MaxInt64)

// 自定义信号
var singal = struct{}{}

type LocalTimerServer struct {
	Id         uint64         // 自增的序列号
	JobList    *rbtree.Rbtree // 基于内存的事件容器
	Count      uint64         // 统计已执行的任务次数
	PauseChan  chan struct{}  // 暂停
	ResumeChan chan struct{}  // 重置
	ExitChan   chan struct{}  // 退出
}

type Ijob interface {
	Notify() <-chan Job // 获取一个chan:当job被执行时，可以从chan中获取消息
	GetCount() uint64   // 获取已经执行的次数
	GetTimes() uint64   // 获取允许执行的最大次数
}

type Job struct {
	Id           uint64        // 唯一id，由timerserver生成，用来区分同一时刻的不同事件
	Times        uint64        // 允许执行的最大次数:0表示无数次
	Count        uint64        // 表示已执行的次数
	IntervalTime time.Duration // 间隔时间：支持ns
	CreateTime   time.Time     // 创建时间
	ActionTime   time.Time     // FIXME 计算得出的本次执行时间点，会有误差
	JobHandler   func()        // 事件函数
	MsgChan      chan Job      // 消息通道，执行时，控制器通过该通道向外部传递消息
}

func NewLocalTimerServer() *LocalTimerServer {
	server := &LocalTimerServer{
		JobList:    rbtree.New(),
		PauseChan:  make(chan struct{}, 0),
		ResumeChan: make(chan struct{}, 0),
		ExitChan:   make(chan struct{}, 0),
	}
	server.start()
	logger.Info("gentimer run")
	return server
}

func (server *LocalTimerServer) start() {
	defalutJob := Job{
		CreateTime:   time.Now(),
		IntervalTime: time.Duration(math.MaxInt64),
		JobHandler:   func() {},
	}
	server.addJob(time.Now(), defalutJob.IntervalTime, 1, defalutJob.JobHandler)
	go server.schedule()
	server.resume()
}

func (server *LocalTimerServer) schedule() {
	timer := time.NewTimer(DEFAULT)
	defer timer.Stop()
PAUSE:
	<-server.ResumeChan
	for {
		value := server.JobList.Min()
		job, _ := value.(*Job)
		timeout := job.ActionTime.Sub(time.Now())
		timer.Reset(timeout)
		select {
		case <-timer.C:
			server.Count++
			job.ExecWithGo(true)
			if job.Times == 0 || job.Times > job.Count {
				server.JobList.Delete(job)
				job.ActionTime = job.ActionTime.Add(job.IntervalTime)
				server.JobList.Insert(job)
			} else {
				server.removeJob(job)
			}
		case <-server.PauseChan:
			goto PAUSE
		case <-server.ExitChan:
			goto EXIT
		}
	}
EXIT:
}

func (server *LocalTimerServer) resume() {
	server.ResumeChan <- singal
}

func (server *LocalTimerServer) exit() {
	server.ExitChan <- singal
}

func (server *LocalTimerServer) pause() {
	server.PauseChan <- singal
}

// 重设指定任务的超时时间
func (server *LocalTimerServer) UpdateJobTimeout(job Ijob, timeout time.Duration) bool {
	// 1s = 1e9 ns
	if timeout.Nanoseconds() <= 0 {
		return false
	}
	now := time.Now()
	server.pause()
	defer server.resume()
	item, ok := job.(*Job)
	if !ok {
		return false
	}
	server.JobList.Delete(item)
	item.ActionTime = now.Add(timeout)
	server.JobList.Insert(item)
	return true
}

// 添加单次任务，适用于单次定时业务逻辑
func (server *LocalTimerServer) AddJobWithInterval(timeout time.Duration, jobFunc func()) (Ijob, bool) {
	if timeout.Nanoseconds() <= 0 {
		return nil, false
	}
	server.pause()
	defer server.resume()
	return server.addJob(time.Now(), timeout, 1, jobFunc)
}

// 添加单次任务，适用于特定时期的活动等类型的任务
func (server *LocalTimerServer) AddJobWithDeadtime(deadtime time.Time, jobFunc func()) (Ijob, bool) {
	now := time.Now()
	timeout := deadtime.Sub(now)
	if timeout.Nanoseconds() <= 0 {
		return nil, false
	}
	server.pause()
	defer server.resume()
	return server.addJob(now, timeout, 1, jobFunc)
}

// 添加重复任务
func (server *LocalTimerServer) AddJobRepeat(jobInterval time.Duration, times uint64, jobFunc func()) (Ijob, bool) {
	if jobInterval.Nanoseconds() <= 0 {
		return nil, false
	}
	server.pause()
	defer server.resume()
	return server.addJob(time.Now(), jobInterval, times, jobFunc)
}

// 移除指定的单项任务
func (server *LocalTimerServer) DelJob(job Ijob) bool {
	if job == nil {
		return false
	}
	server.pause()
	defer server.resume()
	item, ok := job.(*Job)
	if !ok {
		return false
	}
	server.removeJob(item)
	return true
}

// 移除指定的多项任务
func (server *LocalTimerServer) DelJobs(jobs []Ijob) {
	server.pause()
	defer server.resume()
	for _, job := range jobs {
		item, ok := job.(*Job)
		if !ok {
			continue
		}
		server.removeJob(item)
	}
}

// 获取server当前执行次数
func (server *LocalTimerServer) GetCount() uint64 {
	return atomic.LoadUint64(&server.Count)
}

// 重置server的内部状态
func (server *LocalTimerServer) Reset() *LocalTimerServer {
	server.exit()
	server.Count = 0
	server.cleanJobs()
	server.start()
	logger.Info("gentimer Reset")
	return server
}

func (server *LocalTimerServer) cleanJobs() {
	item := server.JobList.Min()
	for item != nil {
		job, ok := item.(*Job)
		if ok {
			server.removeJob(job)
		}
		item = server.JobList.Min()
	}
}

// 获取pending的任务数量
func (server *LocalTimerServer) WaitJobs() uint {
	leng := server.JobList.Len() - 1
	if leng > 0 {
		return leng
	}
	return 0
}

func (server *LocalTimerServer) addJob(createTime time.Time, intervalTime time.Duration, jobTimes uint64, jobFunc func()) (job *Job, inserted bool) {
	inserted = true
	server.Id++
	job = &Job{
		Id:           server.Id,
		Times:        jobTimes,
		CreateTime:   createTime,
		ActionTime:   createTime.Add(intervalTime),
		IntervalTime: intervalTime,
		MsgChan:      make(chan Job, 10),
		JobHandler:   jobFunc,
	}
	server.JobList.Insert(job)
	return
}

func (server *LocalTimerServer) removeJob(job *Job) {
	server.JobList.Delete(job)
	close(job.MsgChan)
}

// 优雅的关闭
// 会依次执行一次（不按照actiontime）
func (server *LocalTimerServer) StopByGrace() {
	server.exit()
	server.immediate()
	logger.Warning("gentimer StopByGrace")
}

// 强制关闭
func (server *LocalTimerServer) StopByForce() {
	server.exit()
	server.cleanJobs()
	logger.Warning("gentimer StopByForce")
}

func (server *LocalTimerServer) immediate() {
	item := server.JobList.Min()
	for item != nil {
		job, ok := item.(*Job)
		if ok {
			atomic.AddUint64(&server.Count, 1)
			job.ExecWithGo(false)
			server.removeJob(job)
		}
		item = server.JobList.Min()
	}
}

/*************************************Job相关************************************/

/****************************实现了rbtree.Item接口****************************/

func (job Job) Less(another rbtree.Item) bool {
	item, ok := another.(*Job)
	if !ok {
		return false
	}
	if !job.ActionTime.Equal(item.ActionTime) {
		return job.ActionTime.Before(item.ActionTime)
	}
	return job.Id < item.Id
}

/********************************实现了Job接口***********************************/

func (job Job) Notify() <-chan Job {
	return job.MsgChan
}

func (job Job) GetCount() uint64 {
	return job.Count
}

func (job Job) GetTimes() uint64 {
	return job.Times
}

func (job *Job) action() {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	job.JobHandler()
}

func (job *Job) ExecWithGo(isGo bool) {
	job.Count++
	if job.JobHandler == nil {
		return
	}
	if isGo {
		go job.action()
	} else {
		job.action()
	}
}
