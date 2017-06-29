package gserver

import (
	"github.com/gfandada/gserver/logger"
	"github.com/gfandada/gserver/module"
)

func Run(mods ...module.Imodule) {
	logger.Start()
	module.Run(mods...)
}
