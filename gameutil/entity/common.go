package entity

import (
	"reflect"
)

// Entityid(也是spaceid)
type EntityID string

// Entity容器
type EntitySet map[*Entity]struct{}

// Entity定义
type Entity struct {
	ID                   EntityID        // 唯一id
	IV                   reflect.Value   // Ientity实际类型
	TypeName             string          // 名称
	I                    Ientity         // 行为装载器
	Space                *Space          // space
	destroyed            bool            // 销毁标记:true-已销毁
	typeDesc             *EntityTypeDesc // Entity描述
	aoi                  aoi             // aoi
	declaredServices     StringSet       // 声明的服务的容器
	client               *GameClient     // entity的网络对象
	enteringSpaceRequest struct {
		SpaceID     EntityID
		EnterPos    Vector3
		RequestTime int64 // 请求时间:单位纳秒
	}
}

type StringSet map[string]struct{}

const (
	_SPACE_ENTITY_TYPE   = "__space__"
	_SPACE_KIND_ATTR_KEY = "_K"

	_DEFAULT_AOI_DISTANCE = 100
)
