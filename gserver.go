package gserver

import (
	"github.com/gfandada/gserver/cluster"
	"github.com/gfandada/gserver/logger"
	"github.com/gfandada/gserver/module"
)

// 插件化运行模块
func Run(mods ...module.Imodule) {
	module.Run(mods...)
}

// 启动日志
// path: 日志配置文件
func RunLogger(path string) {
	logger.Start(path)
}

// 初始化集群
func RunCluster(path string) {
	cluster.Init(path)
}
