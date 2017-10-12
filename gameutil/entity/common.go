package entity

type EntityId string

type SpaceId string

type EntitySet map[*Entity]struct{}

const (
	_DEFAULT_AOI_DISTANCE = 2
)

const (
	NIL_SPACE           = 0 // 空场景
	DEFAULT_FIGHT_SPACE = 1 // 默认的战斗场景
)

const (
	NEUTRAL_FLAG        = 0 // 中立方
	ENEMY_SHIP_FLAG     = 1 // 敌方战船
	FRIEND_SHIP_FLAG    = 2 // 我方战船
	ENEMY_SOLDIER_FLAG  = 3 // 敌方小兵
	FRIEND_SOLDIER_FLAG = 4 // 我方小兵
	ENEMY_TOWER_FLAG    = 5 // 敌方防御塔
	FRIEND_TOWER_FLAG   = 6 // 我方防御塔
)
