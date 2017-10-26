package entity

import (
	"fmt"
	//	"fmt"
	"strconv"
	"testing"

	"github.com/gfandada/gserver/util"
)

// 一万个可行走点测试aoi
const (
	X = 100
	Y = 0
	Z = 100
)

func Test_space(t *testing.T) {
	// 构建一个场景
	space := NewSpace(1, new(Space))
	RegisterSpace(space)
	for i := 0; i < X; i++ {
		for j := 0; j < Z; j++ {
			entity := NewEntity(&EntityDesc{
				Name:   "entity" + strconv.Itoa(i) + ":" + strconv.Itoa(j),
				UseAOI: true,
			})
			entity.BindIentity(new(Entity))
			RegisterEntity(entity)
			entity.EnterSpace(space.Id, Vector3{
				X: Coord(i),
				Y: Coord(0),
				Z: Coord(j),
			})
		}
	}
	// 主进程调度
	for i := 0; i < 1; i++ {
		for key := range space.entities {
			x := util.RandInterval(0, X-1)
			z := util.RandInterval(0, Z-1)
			//fmt.Println("调度时的随机坐标", x, ":", z)
			key.MoveSpace(Vector3{
				X: Coord(x),
				Y: Coord(0),
				Z: Coord(z),
			})
		}
		for key := range space.entities {
			fmt.Println("我是", key.aoi.pos, key.Desc.Name, "我有", len(key.Neighbors()), "个邻居")
		}
	}
}
