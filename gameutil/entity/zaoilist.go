package entity

type zAOIList struct {
	head *aoi
	tail *aoi
}

func newZAOIList() *zAOIList {
	return &zAOIList{}
}

func (sl *zAOIList) Insert(aoi *aoi) {
	insertCoord := aoi.pos.Z
	if sl.head != nil {
		p := sl.head
		for p != nil && p.pos.Z < insertCoord {
			p = p.zNext
		}
		if p == nil {
			tail := sl.tail
			tail.zNext = aoi
			aoi.zPrev = tail
			sl.tail = aoi
		} else {
			prev := p.zPrev
			aoi.zNext = p
			p.zPrev = aoi
			aoi.zPrev = prev

			if prev != nil {
				prev.zNext = aoi
			} else {
				sl.head = aoi
			}
		}
	} else {
		sl.head = aoi
		sl.tail = aoi
	}
}

func (sl *zAOIList) Remove(aoi *aoi) {
	prev := aoi.zPrev
	next := aoi.zNext
	if prev != nil {
		prev.zNext = next
		aoi.zPrev = nil
	} else {
		sl.head = next
	}
	if next != nil {
		next.zPrev = prev
		aoi.zNext = nil
	} else {
		sl.tail = prev
	}
}

func (sl *zAOIList) Move(aoi *aoi, oldCoord Coord) {
	coord := aoi.pos.Z
	// 后移
	if coord > oldCoord {
		next := aoi.zNext
		if next == nil || next.pos.Z >= coord {
			return
		}
		prev := aoi.zPrev
		if prev != nil {
			prev.zNext = next
		} else {
			sl.head = next
		}
		next.zPrev = prev
		prev, next = next, next.zNext
		for next != nil && next.pos.Z < coord {
			prev, next = next, next.zNext
		}
		prev.zNext = aoi
		aoi.zPrev = prev
		if next != nil {
			next.zPrev = aoi
		} else {
			sl.tail = aoi
		}
		aoi.zNext = next
	} else {
		// 前移
		prev := aoi.zPrev
		if prev == nil || prev.pos.Z <= coord {
			return
		}
		next := aoi.zNext
		if next != nil {
			next.zPrev = prev
		} else {
			sl.tail = prev
		}
		prev.zNext = next
		next, prev = prev, prev.zPrev
		for prev != nil && prev.pos.Z > coord {
			next, prev = prev, prev.zPrev
		}
		next.zPrev = aoi
		aoi.zNext = next
		if prev != nil {
			prev.zNext = aoi
		} else {
			sl.head = aoi
		}
		aoi.zPrev = prev
	}
}

func (sl *zAOIList) Mark(aoi *aoi) {
	prev := aoi.zPrev
	coord := aoi.pos.Z
	minCoord := coord - _DEFAULT_AOI_DISTANCE
	for prev != nil && prev.pos.Z >= minCoord {
		prev.markVal += 1
		prev = prev.zPrev
	}
	next := aoi.zNext
	maxCoord := coord + _DEFAULT_AOI_DISTANCE
	for next != nil && next.pos.Z <= maxCoord {
		next.markVal += 1
		next = next.zNext
	}
}

// 清理z轴上不符合条件的mark标记
func (sl *zAOIList) ClearMark(aoi *aoi) {
	prev := aoi.zPrev
	coord := aoi.pos.Z
	minCoord := coord - _DEFAULT_AOI_DISTANCE
	for prev != nil && prev.pos.Z >= minCoord {
		prev.markVal = 0
		prev = prev.zPrev
	}
	next := aoi.zNext
	maxCoord := coord + _DEFAULT_AOI_DISTANCE
	for next != nil && next.pos.Z <= maxCoord {
		next.markVal = 0
		next = next.zNext
	}
}
