package gateway

import (
	Loader "github.com/gfandada/gserver/loader"
	"github.com/gfandada/gserver/network"
)

/****************************实现imodule接口******************************/

type WsGateway struct {
	Config     string
	Coder      network.Imessage
	configdata *network.Config
	wsServer   *network.WsServer
}

func (wg *WsGateway) OnInit() {
	config := new(network.Config)
	Loader.LoadJson(wg.Config, config)
	wg.configdata = config
	wg.configdata.MsgParser = wg.Coder
	wg.configdata.MsgParser.SetMaxLen(config.MaxMsgLen)
	wg.configdata.Parser = network.NewMessageParser()
	wg.configdata.Parser.SetMsgLen(uint16(config.MaxMsgLen), uint16(config.MinMsgLen))
	wg.configdata.Gate = &Agent{configdata: wg.configdata}
}

func (wg *WsGateway) OnDestroy() {
	wg.wsServer.Close()
}

func (wg *WsGateway) Run(ChClose chan bool) {
	wg.wsServer = network.StartWs(wg.configdata)
}
