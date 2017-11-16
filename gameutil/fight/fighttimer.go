package fight

import (
	"time"

	. "github.com/gfandada/gserver/gameutil/entity"
	. "github.com/gfandada/gserver/gservices"
)

type FightTimer struct {
	server *LocalTimerServer
	jobs   map[EntityId]Ijob
}

func (f *FightTimer) init() {
	f.server = NewLocalTimerServer()
	f.jobs = make(map[EntityId]Ijob)
}

func (f *FightTimer) stop() {
	if f.server != nil {
		f.server.StopByForce()
	}
	f.jobs = nil
}

func (f *FightTimer) AddRepeatJob(entityId EntityId, jobInterval time.Duration, times uint64,
	jobFunc MessageHandler1, args []interface{}) {
	job, err := f.server.AddJobRepeat(jobInterval, times, jobFunc, args)
	if !err {
		// LOG
		return
	}
	f.jobs[entityId] = job
}

func (f *FightTimer) AddOneJob(entity *Entity, jobInterval time.Duration, jobFunc MessageHandler1,
	args []interface{}) {
	job, err := f.server.AddJobWithInterval(jobInterval, jobFunc, args)
	if !err {
		// LOG
		return
	}
	f.jobs[entity.Id] = job
}

func (f *FightTimer) DelAiJob(entityId EntityId) {
	f.server.DelJob(f.jobs[entityId])
	delete(f.jobs, entityId)
}
