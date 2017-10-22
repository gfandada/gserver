package entity

type EntityId string

type SpaceId string

type EntitySet map[*Entity]struct{}

const (
	_DEFAULT_AOI_DISTANCE = 10
)

const (
	NIL_SPACE           = 0 // 空场景
	DEFAULT_FIGHT_SPACE = 1 // 默认的战斗场景
)
