package entity

// aoi计算器
type Iaoicalculator interface {
	Enter(aoi *aoi, pos Vector3)   // 进入（Entity进入Space时调用）
	Leave(aoi *aoi)                // 离开（Entity离开Space时调用
	Move(aoi *aoi, newPos Vector3) // 移动（Entity在Space里移动时调用）
	// 调整关注区域列表（即调整邻居的过程，可能有新邻居，可能有失效的邻居，
	// 那么也就有不同的操作，对新邻居是进入，对失效邻居是离开）
	// @return enter:新邻居 leave:失效的邻居
	Adjust(aoi *aoi) (enter []*aoi, leave []*aoi)
}
