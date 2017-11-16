package entity

import (
	"fmt"
	"math"
	"unsafe"
)

// 单位
type Coord float32

// 位置
type Vector3 struct {
	X    Coord // X轴
	Y    Coord // Y轴
	Z    Coord // Z轴
	VX   Coord // X轴上的速度
	VZ   Coord // Z轴上的速度
	W    Coord // 度宽，for体积
	H    Coord // 高宽，for体积
	TIME int64 // 时间戳：ms
}

// AOI
type aoi struct {
	pos       Vector3   // 当前位置
	neighbors EntitySet // 邻居列表
	xNext     *aoi      // x轴后指针
	xPrev     *aoi      // x轴前指针
	zNext     *aoi      // z轴后指针
	zPrev     *aoi      // z轴前指针
	markVal   int       // FIXME 用来标记邻居
}

func (p Vector3) String() string {
	return fmt.Sprintf("(%.2f, %.2f, %.2f)", p.X, p.Y, p.Z)
}

// 计算p、o两个位置间的距离
func (p Vector3) DistanceTo(o Vector3) Coord {
	dx := p.X - o.X
	dy := p.Y - o.Y
	dz := p.Z - o.Z
	return Coord(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}

// p-o
func (p Vector3) Sub(o Vector3) Vector3 {
	return Vector3{p.X - o.X, p.Y - o.Y, p.Z - o.Z, p.VX - o.VX, p.VZ - o.VZ,
		0, 0, 0}
}

// p+o
func (p Vector3) Add(o Vector3) Vector3 {
	return Vector3{p.X + o.X, p.Y + o.Y, p.Z + o.Z, p.VX + o.VX, p.VZ + o.VZ,
		0, 0, 0}
}

// p*m
func (p Vector3) Mul(m Coord) Vector3 {
	return Vector3{p.X * m, p.Y * m, p.Z * m, p.VX * m, p.VZ * m,
		0, 0, 0}
}

func (p *Vector3) Normalize() {
	d := Coord(math.Sqrt(float64(p.X*p.X + p.Y + p.Y + p.Z*p.Z)))
	if d == 0 {
		return
	}
	p.X /= d
	p.Y /= d
	p.Z /= d
}

func (p Vector3) Normalized() Vector3 {
	p.Normalize()
	return p
}

func initAOI(aoi *aoi) {
	aoi.neighbors = EntitySet{}
}

// 获取指定aoi自己的entity
var aoiFieldOffset uintptr

func init() {
	dummyEntity := (*Entity)(unsafe.Pointer(&aoiFieldOffset))
	aoiFieldOffset = uintptr(unsafe.Pointer(&dummyEntity.aoi)) - uintptr(unsafe.Pointer(dummyEntity))
}

// 获取aoi自己的Entity
func (aoi *aoi) getEntity() *Entity {
	return (*Entity)(unsafe.Pointer((uintptr)(unsafe.Pointer(aoi)) - aoiFieldOffset))
}

// 添加一个关注区域
func (aoi *aoi) interest(other *Entity) {
	aoi.neighbors.Add(other)
}

// 移除一个关注区域
func (aoi *aoi) uninterest(other *Entity) {
	aoi.neighbors.Del(other)
}

type aoiSet map[*aoi]struct{}

func (aoiset aoiSet) Add(aoi *aoi) {
	aoiset[aoi] = struct{}{}
}

func (aoiset aoiSet) Del(aoi *aoi) {
	delete(aoiset, aoi)
}

func (aoiset aoiSet) Contains(aoi *aoi) bool {
	_, ok := aoiset[aoi]
	return ok
}

func (aoiset aoiSet) Join(other aoiSet) aoiSet {
	join := aoiSet{}
	for aoi := range aoiset {
		if other.Contains(aoi) {
			join.Add(aoi)
		}
	}
	return join
}
