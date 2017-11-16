package entity

// 矩形判断
// return true:相交 false:不相交
func collByCircular(entity1, entity2 *Entity) bool {
	pos1 := entity1.GetPosition()
	pos2 := entity2.GetPosition()
	if pos1.X-pos2.X >= 0.1 && pos1.X-(pos2.X+pos2.W/2+pos1.W/2) >= 0.1 {
		return false
	} else if pos1.X-pos2.X <= 0.1 && pos2.X-(pos1.X+pos2.W/2+pos1.W/2) >= 0.1 {
		return false
	} else if pos1.Z-pos2.Z >= 0.1 && pos1.Z-(pos2.Z+pos2.H/2+pos1.H/2) >= 0.1 {
		return false
	} else if pos1.Z-pos2.Z <= 0.1 && pos2.Z-(pos1.Z+pos2.H/2+pos1.H/2) >= 0.1 {
		return false
	}
	return true
}

// 圆形判断
// return true:相交 false:不相交
func collByRectangle(entity1, entity2 *Entity) bool {
	instance := entity1.DistanceTo(entity2)
	return instance-entity1.GetPosition().W/2+entity2.GetPosition().W/2 < 0.1
}
