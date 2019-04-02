package fight

import (
	"time"

	. "github.com/gfandada/gserver/goroutine"
)

func startFightDamageCalc(fightid FightId) error {
	_, err := Start(&fightDamageCalc{
		id: fightid,
	})
	return err
}

func stopFightDamageCalc(fightid FightId) error {
	return Stop(NewDamageCalcAlias(fightid))
}

// 同步调用
// default timeout 1s
// @params  fightid:战斗id  msg:消息类型  args:自定义参数
func CallFightDamageCalc(fightid FightId, msg string, args []interface{}) ([]interface{}, error) {
	return Call(NewDamageCalcAlias(fightid), msg, args, 1)
}

// 异步调用
// @params  fightid:战斗id  msg:消息类型  args:自定义参数
func CastFightDamageCalc(fightid FightId, msg string, args []interface{}) {
	Cast(NewDamageCalcAlias(fightid), msg, args)
}

// 解析
func ParseDamageCalcInner(inner []interface{}) FightId {
	return inner[0].(FightId)
}

type fightDamageCalc struct {
	id FightId
}

func (f *fightDamageCalc) Name() string {
	return NewDamageCalcAlias(f.id)
}

func (f *fightDamageCalc) SetTimer() time.Duration {
	return time.Millisecond * 0
}

func (f *fightDamageCalc) Init() {
	if handler := GetHandler(INIT_DAMAGE); handler != nil {
		handler(nil, []interface{}{})
	}
}

func (f *fightDamageCalc) Close() {
	if handler := GetHandler(CLOSE_DAMAGE); handler != nil {
		handler(nil, []interface{}{})
	}
}

func (f *fightDamageCalc) TimerWork() {
	if handler := GetHandler(TIMER_DAMAGE); handler != nil {
		handler([]interface{}{f.id}, []interface{}{})
	}
}

func (f *fightDamageCalc) Handler(msg string, args []interface{}, ret chan []interface{}) {
	if handler := GetHandler(msg); handler != nil {
		rets := handler([]interface{}{f.id}, args)
		// when rets are nil, should be return instead of timeout
		if ret != nil {
			ret <- rets
		}
	}
}
