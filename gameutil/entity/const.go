package entity

import (
	"time"
)

const (
	// 帧时间
	GAME_SERVICE_TICK_INTERVAL = time.Millisecond * 10

	// 进入space请求超时时间
	ENTER_SPACE_REQUEST_TIMEOUT = DISPATCHER_MIGRATE_TIMEOUT + time.Minute
	// entity迁移超时时间
	DISPATCHER_MIGRATE_TIMEOUT = time.Minute * 5
)

const (
	PLAYER_TYPE = "player_type_"
)
