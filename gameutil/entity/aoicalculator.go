package entity

type XZListAOICalculator struct {
	xSweepList *xAOIList
	zSweepList *zAOIList
}

func newXZListAOICalculator() *XZListAOICalculator {
	return &XZListAOICalculator{
		xSweepList: newXAOIList(),
		zSweepList: newZAOIList(),
	}
}

func (cal *XZListAOICalculator) Enter(aoi *aoi, pos Vector3) {
	aoi.pos = pos
	cal.xSweepList.Insert(aoi)
	cal.zSweepList.Insert(aoi)
}

func (cal *XZListAOICalculator) Leave(aoi *aoi) {
	cal.xSweepList.Remove(aoi)
	cal.zSweepList.Remove(aoi)
}

func (cal *XZListAOICalculator) Move(aoi *aoi, pos Vector3) {
	oldPos := aoi.pos
	aoi.pos = pos
	if oldPos.X != pos.X {
		cal.xSweepList.Move(aoi, oldPos.X)
	}
	if oldPos.Z != pos.Z {
		cal.zSweepList.Move(aoi, oldPos.Z)
	}
}

func (cal *XZListAOICalculator) Adjust(aoi *aoi) (enter []*aoi, leave []*aoi) {
	// 邻居判断
	cal.xSweepList.Mark(aoi)
	cal.zSweepList.Mark(aoi)
	for neighbor := range aoi.neighbors {
		naoi := &neighbor.aoi
		if naoi.markVal == 2 {
			naoi.markVal = -2
		} else {
			// 失效的邻居
			leave = append(leave, naoi)
		}
	}
	enter = cal.xSweepList.GetClearMarkedNeighbors(aoi)
	cal.zSweepList.ClearMark(aoi)
	return
}
