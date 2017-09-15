package main

import (
	Module "./module"
	Network "./network"
	Services "./services"
	Discovery "./services/discovery"
	GateService "./services/gateway"
	Service "./services/service"
)

// 运行模块
// @params 模块转载器...
func Run(mods ...Module.Imodule) {
	Module.Run(mods...)
}

// 运行websocket网关
// @params path:配置文件 discpath:服务发现配置
func RunWSGateway(path string, discpath string) {
	gate := new(GateService.WsGateway)
	gate.Config = path
	Discovery.Init(discpath)
	Run(gate)
}

// 运行service
// @params path:配置 discpath:服务发现配置
func RunCluster(path string, discpath string) {
	service := new(Service.Service)
	service.Config = path
	Run(service)
}

// 注册消息handler
// @params list:消息列表
func RegisterHandler(list []*Services.MsgHandler) {
	for _, v := range list {
		Services.Register(v.MsgId, v.MsgHandler)
	}
}

// 绑定消息编码器
// @params coder:消息编码器
func BindCoder(coder Network.Imessage) {
}

// 注册消息至编码器中
// @params coder:消息编码器
func RegisterCoder(list []*Network.RawMessage, coder Network.Imessage) {
	for _, v := range list {
		coder.Register(v)
	}
}
