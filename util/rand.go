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
