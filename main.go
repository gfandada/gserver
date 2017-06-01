package main

import (
	"gserver/src/gateway"
	"gserver/src/lib/module"
)

func main() {
	// TODO 暂时通过代码组装配置
	gate := new(gateway.Gate)
	gate.MaxConnNum = 1000
	gate.PendingNum = 100
	gate.MaxMsgLen = 1024
	gate.ServerAddress = "localhost:9527"
	gate.LenMsgLen = 1024
	module.Run(gate)
}
