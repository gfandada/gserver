package entity

import (
	"fmt"
	"strconv"
	"testing"
)

// 模拟一个场景
// 模拟加入1000个entity
// 1000个entity并发的移动
// 通过日志记录每次移动时的邻居列表以及enter和leave列表
func Test_test(t *testing.T) {
	// 创建一个场景
	RegisterSpace(new(Space))
	// 创建1000个entity
	ids := make([]int, 1000)
	for id, _ := range ids {
		idstr := strconv.Itoa(id)
		if ret := RegisterEntity("player"+idstr, new(Entity), true, true); ret == nil {
			t.Error("desc error")
		}
	}
}
