package entity

import (
	//	"fmt"
	"testing"
)

// for 10*10
// min 1 * 1
func Test_collision(t *testing.T) {
	entity1 := &Entity{
		aoi: aoi{
			pos: Vector3{
				X: Coord(0),
				Z: Coord(0),
				W: Coord(2),
				H: Coord(2),
			},
		},
	}
	count := 0
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			entity2 := &Entity{
				aoi: aoi{
					pos: Vector3{
						X: Coord(i),
						Z: Coord(j),
						W: Coord(2),
						H: Coord(3),
					},
				},
			}
			if collByCircular(entity1, entity2) {
				count++
			}
		}
	}
	// for entity2
	// 7*7=25
	// 6*6=25
	// 5*5=16
	// 4*4=16
	// 3*3=9
	// 2*2=9
	//fmt.Println("=======相交的矩形========", count)
}
