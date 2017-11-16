package fight

import (
	"time"

	. "github.com/gfandada/gserver/gameutil/entity"
	. "github.com/gfandada/gserver/goroutine"
)

func startFightScheduler(fightid FightId, flag int) error {
	_, err := Start(&fightScheduler{
		id:   fightid,
		flag: flag,
	})
	return err
}

func stopFightScheduler(fightid FightId) error {
	return StopByName(NewFightSchedulerAlias(fightid))
}

// 同步调用
// default timeout 1s
// @params  fightid:战斗id  msg:消息类型  args:自定义参数
func CallFightScheduler(fightid FightId, msg string, args []interface{}) ([]interface{}, error) {
	return CallByName(NewFightSchedulerAlias(fightid), msg, args, 1)
}

// 异步调用
// @params  fightid:战斗id  msg:消息类型  args:自定义参数
func CastFightScheduler(fightid FightId, msg string, args []interface{}) {
	CastByName(NewFightSchedulerAlias(fightid), msg, args)
}

// 解析
func ParseSchedulerInner(inner []interface{}) (int, FightId, *Space,
	map[*Entity]struct{}, map[*Entity]struct{}, map[*Entity]struct{},
	*FightTimer) {
	return inner[0].(int),
		inner[1].(FightId),
		inner[2].(*Space),
		inner[3].(map[*Entity]struct{}),
		inner[4].(map[*Entity]struct{}),
		inner[5].(map[*Entity]struct{}),
		inner[6].(*FightTimer)
}

type fightScheduler struct {
	id           FightId              // 战斗id
	flag         int                  // 3v3 or 5v5 等
	defalutSpace *Space               // 场景
	ships        map[*Entity]struct{} // 玩家
	soldiers     map[*Entity]struct{} // 小兵
	towers       map[*Entity]struct{} // 防御塔和水晶容器
	timer        *FightTimer          // 定时任务调度器
}

func (f *fightScheduler) Name() string {
	return NewFightSchedulerAlias(f.id)
}

func (f *fightScheduler) Timer() time.Duration {
	return time.Millisecond * 100
}

func (f *fightScheduler) InitGo() {
	f.ships = make(map[*Entity]struct{})
	f.soldiers = make(map[*Entity]struct{})
	f.towers = make(map[*Entity]struct{})
	f.defalutSpace = NewDefaultSpace()
	f.timer = new(FightTimer)
	f.timer.init()
	startFightDamageCalc(f.id)
	startFightAward(f.id)
	startFighPost(f.id)
	if handler := GetHandler(INIT_SCHEDULER); handler != nil {
		handler([]interface{}{f.flag, f.id, f.defalutSpace, f.ships,
			f.soldiers, f.towers, f.timer}, []interface{}{})
	}
}

func (f *fightScheduler) CloseGo() {
	f.timer.stop()
	stopFightDamageCalc(f.id)
	stopFightAward(f.id)
	stopFighPost(f.id)
	for entity := range f.ships {
		entity.LeaveSpace()
		UnRegisterEntity(entity.Id)
	}
	for entity := range f.soldiers {
		entity.LeaveSpace()
		UnRegisterEntity(entity.Id)
	}
	for entity := range f.towers {
		entity.LeaveSpace()
		UnRegisterEntity(entity.Id)
	}
	f.ships = nil
	f.soldiers = nil
	f.towers = nil
	UnRegisterSpace(f.defalutSpace.Id)
	if handler := GetHandler(CLOSE_SCHEDULER); handler != nil {
		handler([]interface{}{f.flag, f.id, f.defalutSpace, f.ships,
			f.soldiers, f.towers, f.timer}, []interface{}{})
	}
}

func (f *fightScheduler) Timer_work() {
	if f.flag*2 != len(f.ships) {
		return
	}
	if handler := GetHandler(TIMER_SCHEDULER); handler != nil {
		handler([]interface{}{f.flag, f.id, f.defalutSpace, f.ships,
			f.soldiers, f.towers, f.timer}, []interface{}{})
	}
}

func (f *fightScheduler) Handler(msg string, args []interface{}, ret chan []interface{}) {
	if handler := GetHandler(msg); handler != nil {
		rets := handler([]interface{}{f.flag, f.id, f.defalutSpace, f.ships,
			f.soldiers, f.towers, f.timer}, args)
		// when rets are nil, should be return instead of timeout
		if ret != nil {
			ret <- rets
		}
	}
}
