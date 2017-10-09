package entity

type xAOIList struct {
	head *aoi
	tail *aoi
}

func newXAOIList() *xAOIList {
	return &xAOIList{}
}

func (sl *xAOIList) Insert(aoi *aoi) {
	insertCoord := aoi.pos.X
	if sl.head != nil {
		p := sl.head
		for p != nil && p.pos.X < insertCoord {
			p = p.xNext
		}
		if p == nil {
			tail := sl.tail
			tail.xNext = aoi
			aoi.xPrev = tail
			sl.tail = aoi
		} else {
			prev := p.xPrev
			aoi.xNext = p
			p.xPrev = aoi
			aoi.xPrev = prev
			if prev != nil {
				prev.xNext = aoi
			} else {
				sl.head = aoi
			}
		}
	} else {
		sl.head = aoi
		sl.tail = aoi
	}
}

func (sl *xAOIList) Remove(aoi *aoi) {
	prev := aoi.xPrev
	next := aoi.xNext
	if prev != nil {
		prev.xNext = next
		aoi.xPrev = nil
	} else {
		sl.head = next
	}
	if next != nil {
		next.xPrev = prev
		aoi.xNext = nil
	} else {
		sl.tail = prev
	}
}

func (sl *xAOIList) Move(aoi *aoi, oldCoord Coord) {
	coord := aoi.pos.X
	// 后移
	if coord > oldCoord {
		next := aoi.xNext
		if next == nil || next.pos.X >= coord {
			return
		}
		prev := aoi.xPrev
		if prev != nil {
			prev.xNext = next
		} else {
			sl.head = next
		}
		next.xPrev = prev
		prev, next = next, next.xNext
		for next != nil && next.pos.X < coord {
			prev, next = next, next.xNext
		}
		prev.xNext = aoi
		aoi.xPrev = prev
		if next != nil {
			next.xPrev = aoi
		} else {
			sl.tail = aoi
		}
		aoi.xNext = next
	} else {
		// 前移
		prev := aoi.xPrev
		if prev == nil || prev.pos.X <= coord {
			return
		}
		next := aoi.xNext
		if next != nil {
			next.xPrev = prev
		} else {
			sl.tail = prev
		}
		prev.xNext = next
		next, prev = prev, prev.xPrev
		for prev != nil && prev.pos.X > coord {
			next, prev = prev, prev.xPrev
		}
		next.xPrev = aoi
		aoi.xNext = next
		if prev != nil {
			prev.xNext = aoi
		} else {
			sl.head = aoi
		}
		aoi.xPrev = prev
	}
}

// 双向遍历用于标记邻居状态
// X/Z轴上同时标记就是邻居，即mark值=2
func (sl *xAOIList) Mark(aoi *aoi) {
	prev := aoi.xPrev
	coord := aoi.pos.X
	minCoord := coord - _DEFAULT_AOI_DISTANCE
	for prev != nil && prev.pos.X >= minCoord {
		prev.markVal += 1
		prev = prev.xPrev
	}
	next := aoi.xNext
	maxCoord := coord + _DEFAULT_AOI_DISTANCE
	for next != nil && next.pos.X <= maxCoord {
		next.markVal += 1
		next = next.xNext
	}
}

// 获取aoi的邻居
func (sl *xAOIList) GetClearMarkedNeighbors(aoi *aoi) (enter []*aoi) {
	prev := aoi.xPrev
	coord := aoi.pos.X
	minCoord := coord - _DEFAULT_AOI_DISTANCE
	for prev != nil && prev.pos.X >= minCoord {
		if prev.markVal == 2 {
			enter = append(enter, prev)
		}
		prev.markVal = 0
		prev = prev.xPrev
	}
	next := aoi.xNext
	maxCoord := coord + _DEFAULT_AOI_DISTANCE
	for next != nil && next.pos.X <= maxCoord {
		if next.markVal == 2 {
			enter = append(enter, next)
		}
		next.markVal = 0
		next = next.xNext
	}
	return
}
