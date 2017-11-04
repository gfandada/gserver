package entity

import (
	"fmt"
	"strings"
)

const (
	KindPlain    = iota // 平原
	KindRiver           // 河流
	KindMountain        // 山峰
	KindBlocker         // 阻挡
	KindFrom            // 起点
	KindTo              // 终点
	KindPath            // 路径点
)

// 用于渲染
// FIXME 仅用于测试
var KindRunes = map[int]rune{
	KindPlain:    '.',
	KindRiver:    '~',
	KindMountain: 'M',
	KindBlocker:  'X',
	KindFrom:     'F',
	KindTo:       'T',
	KindPath:     'A',
}

// 用于渲染
// FIXME 仅用于测试
var RuneKinds = map[rune]int{
	'.': KindPlain,
	'~': KindRiver,
	'M': KindMountain,
	'X': KindBlocker,
	'F': KindFrom,
	'T': KindTo,
	'A': KindPath,
}

// 定义消耗
var KindCosts = map[int]float64{
	KindPlain:    1.0,
	KindFrom:     1.0,
	KindTo:       1.0,
	KindRiver:    2.0,
	KindMountain: 3.0,
}

type WayPoint struct {
	Kind int      // 类型
	X, Y int      // 坐标
	Id   EntityId // 实体id
	W    World    // 属主
}

// FIXME 忽略体积碰撞
func (way *WayPoint) PathNeighbors() []Pather {
	neighbors := []Pather{}
	direction := [][]int{
		{-1, 0}, // 左方
		{1, 0},  // 右方
		{0, -1}, // 下方
		{0, 1},  // 上方
	}
	for _, offset := range direction {
		if n := way.W.GetWayPoint(way.X+offset[0], way.Y+offset[1]); n != nil &&
			n.Kind != KindBlocker {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

func (way *WayPoint) PathNeighborCost(to Pather) float64 {
	toT := to.(*WayPoint)
	return KindCosts[toT.Kind]
}

func (way *WayPoint) PathEstimatedCost(to Pather) float64 {
	toT := to.(*WayPoint)
	absX := toT.X - way.X
	if absX < 0 {
		absX = -absX
	}
	absY := toT.Y - way.Y
	if absY < 0 {
		absY = -absY
	}
	return float64(absX + absY)
}

// 路径点构成的世界
// map[WayPoint.X]map[WayPoint.Y]*WayPoint
type World map[int]map[int]*WayPoint

func (w World) GetWayPoint(x, y int) *WayPoint {
	if w[x] == nil {
		return nil
	}
	return w[x][y]
}

func (w World) SetWayPoint(way *WayPoint, x, y int) {
	if w[x] == nil {
		w[x] = map[int]*WayPoint{}
	}
	w[x][y] = way
	way.X = x
	way.Y = y
	way.W = w
}

func (w World) FirstOfKind(kind int) *WayPoint {
	for _, row := range w {
		for _, wayPoint := range row {
			if wayPoint.Kind == kind {
				return wayPoint
			}
		}
	}
	return nil
}

// 获取起点
func (w World) From() *WayPoint {
	return w.FirstOfKind(KindFrom)
}

// 获取终点
func (w World) To() *WayPoint {
	return w.FirstOfKind(KindTo)
}

// 渲染路径
// FIXME 仅用于测试
func (w World) RenderPath(path []Pather) string {
	width := len(w)
	if width == 0 {
		return ""
	}
	height := len(w[0])
	pathLocs := map[string]bool{}
	for _, p := range path {
		pT := p.(*WayPoint)
		pathLocs[fmt.Sprintf("%d,%d", pT.X, pT.Y)] = true
	}
	rows := make([]string, height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			t := w.GetWayPoint(x, y)
			r := ' '
			if pathLocs[fmt.Sprintf("%d,%d", x, y)] {
				r = KindRunes[KindPath]
			} else if t != nil {
				r = KindRunes[t.Kind]
			}
			rows[y] += string(r)
		}
	}
	return strings.Join(rows, "\n")
}

// 渲染世界
// FIXME 仅用于测试
func ParseWorld(input string) World {
	w := World{}
	for y, row := range strings.Split(strings.TrimSpace(input), "\n") {
		for x, raw := range row {
			kind, ok := RuneKinds[raw]
			if !ok {
				kind = KindBlocker
			}
			w.SetWayPoint(&WayPoint{
				Kind: kind,
			}, x, y)
		}
	}
	return w
}
