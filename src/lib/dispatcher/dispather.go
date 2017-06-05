// 消息处理分发器
package dispather

import (
	"fmt"
)

var CfgMaps map[int]Dispather

type Dispather struct {
	Layer         string      // 层标识
	NodeWeight    string      // 节点权重算法
	MessageHandle interface{} // 消息回调
	Sync          bool        // 同步or异步：ture同步 false异步
}

/*********************************实现idispather接口**********************************/

// TODO 暂未实现从文件加载
func (dispather *Dispather) Load() {
	CfgMaps = make(map[int]Dispather)
	var dispather Dispather
	for i := 0; i < 10; i++ {
		dispather = new(Dispather)
		dispather.Layer = Layer
		dispather.NodeWeight = cluster.Local
	}
}

func (dispather *Dispather) UnLoad() {

}

/*********************************实现idispather接口**********************************/

func (dispather *Dispather) Load() {

}

func (dispather *Dispather) UnLoad() {

}

/*********************************实现imodule接口**********************************/

// 启动分发器
func (dispather *Dispather) Run(chClose chan bool) {

}

// 初始化
func (dispather *Dispather) OnInit() {

}

// 资源回收
func (dispather *Dispather) OnDestroy() {

}

/*****************************封装了dispatcher内部逻辑*********************************/

// 装载配置文件
func loadCfg() {

}

// 启动
func start() {

}
