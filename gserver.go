package gserver

import (
	"github.com/gfandada/gserver/logger"
	"github.com/gfandada/gserver/module"
)

// 插件化运行模块
func Run(mods ...module.Imodule) {
	logger.Start()
	module.Run(mods...)
}

// 启动日志
func RunLogger(path string) {
	logger.Start(path)
}
