package fight

import (
	"time"

	. "github.com/gfandada/gserver/gameutil/entity"
	. "github.com/gfandada/gserver/goroutine"
)

func startFighPost(fightid FightId) error {
	_, err := Start(&fightPost{
		id: fightid,
	})
	return err
}

func stopFighPost(fightid FightId) error {
	return Stop(NewPostAlias(fightid))
}

// 同步调用
// default timeout 1s
// @params  fightid:战斗id  msg:消息类型  args:自定义参数
func CallFighPost(fightid FightId, msg string, args []interface{}) ([]interface{}, error) {
	return Call(NewPostAlias(fightid), msg, args, 1)
}

// 异步调用
// @params  fightid:战斗id  msg:消息类型  args:自定义参数
func CastFighPost(fightid FightId, msg string, args []interface{}) {
	Cast(NewPostAlias(fightid), msg, args)
}

type fightPost struct {
	id    FightId
	ships map[*Entity]struct{}
}

func (f *fightPost) Name() string {
	return NewPostAlias(f.id)
}

func (f *fightPost) SetTimer() time.Duration {
	return time.Millisecond * 0
}

func (f *fightPost) Init() {
	f.ships = make(map[*Entity]struct{})
	if handler := GetHandler(INIT_POST); handler != nil {
		handler(nil, []interface{}{})
	}
}

func (f *fightPost) Close() {
	f.ships = nil
	if handler := GetHandler(CLOSE_POST); handler != nil {
		handler(nil, []interface{}{})
	}
}

func (f *fightPost) TimerWork() {
	if handler := GetHandler(TIMER_POST); handler != nil {
		handler(nil, []interface{}{})
	}
}
func (f *fightPost) Handler(msg string, args []interface{}, ret chan []interface{}) {
	if handler := GetHandler(msg); handler != nil {
		rets := handler(nil, args)
		// when rets are nil, should be return instead of timeout
		if ret != nil {
			ret <- rets
		}
	}
}
