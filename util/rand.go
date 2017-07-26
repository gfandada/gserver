package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 随机从[b1,b2]中获取一个数值
func RandInterval(b1, b2 int32) int32 {
	if b1 == b2 {
		return b1
	}
	min, max := int64(b1), int64(b2)
	if min > max {
		min, max = max, min
	}
	return int32(rand.Int63n(max-min+1) + min)
}

// 随机从[b1,b2]中获取n个数值
func RandIntervalN(b1, b2 int32, n uint32) []int32 {
	if b1 == b2 {
		return []int32{b1}
	}
	min, max := int64(b1), int64(b2)
	if min > max {
		min, max = max, min
	}
	l := max - min + 1
	if int64(n) > l {
		n = uint32(l)
	}
	r := make([]int32, n)
	m := make(map[int32]int32)
	for i := uint32(0); i < n; i++ {
		v := int32(rand.Int63n(l) + min)
		if mv, ok := m[v]; ok {
			r[i] = mv
		} else {
			r[i] = v
		}
		lv := int32(l - 1 + min)
		if v != lv {
			if mv, ok := m[lv]; ok {
				m[v] = mv
			} else {
				m[v] = lv
			}
		}
		l--
	}
	return r
}

// 随机计算
// 方法调用者随机权重值为b1，权重池为b2
// 返回调用者的命中情况：true命中 false未命中
// 此方法不适用于有自定义的起始区域的数值命中
func RandHit(b1, b2 int32) bool {
	if b1 == b2 {
		return true
	}
	if b1 <= 0 {
		return false
	}
	b3 := RandInterval(1, b2)
	if b3 <= b1 {
		return true
	}
	return false
}
