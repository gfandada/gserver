package fight

type FightId string

// 内置开放函数标识
const (
	TIMER_SCHEDULER = "timer_scheduler" // scheduler帧函数
	TIMER_AWARD     = "timer_award"     // award帧函数
	TIMER_DAMAGE    = "timer_damage"    // damage帧函数
	TIMER_POST      = "timer_post"      // post帧函数
	INIT_SCHEDULER  = "init_scheduler"  // scheduler初始化函数
	INIT_AWARD      = "init_award"      // award初始化函数
	INIT_DAMAGE     = "init_damage"     // damage初始化函数
	INIT_POST       = "init_post"       // post初始化函数
	CLOSE_SCHEDULER = "close_scheduler" // scheduler清理函数
	CLOSE_AWARD     = "close_award"     // award清理函数
	CLOSE_DAMAGE    = "close_damage"    // damage清理函数
	CLOSE_POST      = "close_post"      // post清理函数
)
