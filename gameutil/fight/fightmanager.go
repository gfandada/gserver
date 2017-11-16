package fight

import (
	. "github.com/gfandada/gserver/gameutil/entity"
)

// 新建一个战斗调度器进程别名
// @params fightid:战斗id
func NewFightSchedulerAlias(fightid FightId) string {
	return string(fightid) + "scheduler"
}

// 新建一个战斗伤害计算进程别名
// @params fightid:战斗id
func NewDamageCalcAlias(fightid FightId) string {
	return string(fightid) + "damagecalc"
}

// 新建一个战斗奖励计算进程别名
// @params fightid:战斗id
func NewAwardAlias(fightid FightId) string {
	return string(fightid) + "award"
}

// 新建一个post进程别名
// @params fightid:战斗id
func NewPostAlias(fightid FightId) string {
	return string(fightid) + "post"
}

// 创建一个默认的场景
func NewDefaultSpace() *Space {
	space := NewSpace(DEFAULT_FIGHT_SPACE, new(Space))
	RegisterSpace(space)
	return space
}

// 创建一场战斗
// @params fightid:战斗id flag:战斗标识(3v3 or 5v5)
func CreateFight(fightid FightId, flag int) error {
	return startFightScheduler(fightid, flag)
}

// 销毁一场战斗
// @params fightid:战斗id
func DestroyFight(fightid FightId) error {
	return stopFightScheduler(fightid)
}

// for default init
func InitFight() {
	RegisterHandler("", nil)
}
