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
	NEUTRAL_FLAG      = 0 // 中立方
	BLUE_SHIP_FLAG    = 1 // 蓝色方战船
	RED_SHIP_FLAG     = 2 // 红色方战船
	BLUE_SOLDIER_FLAG = 3 // 蓝色方小兵
	RED_SOLDIER_FLAG  = 4 // 红色方小兵
	BLUE_TOWER_FLAG   = 5 // 蓝色方防御塔
	RED_TOWER_FLAG    = 6 // 红色方防御塔
)
