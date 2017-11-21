package fight

import (
	"time"

	. "github.com/gfandada/gserver/gameutil/entity"
	. "github.com/gfandada/gserver/goroutine"
)

// for start
func startFightAward(fightid FightId) error {
	_, err := Start(&fightAward{
		id: fightid,
	})
	return err
}

// for stop
func stopFightAward(fightid FightId) error {
	return StopByName(NewAwardAlias(fightid))
}

// 同步调用
// default timeout 1s
// @params  fightid:战斗id  msg:消息类型  args:自定义参数
func CallFightAward(fightid FightId, msg string, args []interface{}) ([]interface{}, error) {
	return CallByName(NewAwardAlias(fightid), msg, args, 1)
}

// 异步调用
// @params  fightid:战斗id  msg:消息类型  args:自定义参数
func CastFightAward(fightid FightId, msg string, args []interface{}) {
	CastByName(NewAwardAlias(fightid), msg, args)
}

// 解析
func ParseAwardInner(inner []interface{}) (FightId, map[EntityId][]EntityId, map[EntityId][]int) {
	return inner[0].(FightId),
		inner[1].(map[EntityId][]EntityId),
		inner[2].(map[EntityId][]int)
}

type fightAward struct {
	id         FightId
	kill       map[EntityId][]EntityId // 击杀过程(含助攻)
	statistics map[EntityId][]int      // 统计数据:[击杀,死亡,助攻]
}

func (f *fightAward) Name() string {
	return NewAwardAlias(f.id)
}

func (f *fightAward) Timer() time.Duration {
	return time.Millisecond * 0
}

func (f *fightAward) InitGo() {
	f.kill = make(map[EntityId][]EntityId)
	f.statistics = make(map[EntityId][]int)
	if handler := GetHandler(INIT_AWARD); handler != nil {
		handler([]interface{}{f.id, f.kill, f.statistics}, []interface{}{})
	}
}

func (f *fightAward) CloseGo() {
	if handler := GetHandler(CLOSE_AWARD); handler != nil {
		handler([]interface{}{f.id, f.kill, f.statistics}, []interface{}{})
	}
	f.kill = nil
	f.statistics = nil
}

func (f *fightAward) Timer_work() {
	if handler := GetHandler(TIMER_AWARD); handler != nil {
		handler([]interface{}{f.id, f.kill, f.statistics}, []interface{}{})
	}
}

func (f *fightAward) Handler(msg string, args []interface{}, ret chan []interface{}) {
	if handler := GetHandler(msg); handler != nil {
		rets := handler([]interface{}{f.id, f.kill, f.statistics}, args)
		// when rets are nil, should be return instead of timeout
		if ret != nil {
			ret <- rets
		}
	}
}
