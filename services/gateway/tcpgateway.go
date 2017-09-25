package gateway

import (
	Loader "github.com/gfandada/gserver/loader"
	"github.com/gfandada/gserver/network"
)

/****************************实现imodule接口******************************/

type TcpGateway struct {
	Config     string
	Coder      network.Imessage
	configdata *network.Config
	tcpServer  *network.TcpServer
}

func (tcp *TcpGateway) OnInit() {
	config := new(network.Config)
	Loader.LoadJson(tcp.Config, config)
	tcp.configdata = config
	tcp.configdata.MsgParser = tcp.Coder
	tcp.configdata.Parser = network.NewMessageParser()
	tcp.configdata.Parser.SetMsgLen(uint16(config.MaxMsgLen), uint16(config.MinMsgLen))
	tcp.configdata.Gate = &Agent{configdata: tcp.configdata}
}

func (tcp *TcpGateway) OnDestroy() {
	tcp.tcpServer.Close()
}

func (tcp *TcpGateway) Run(ChClose chan bool) {
	tcp.tcpServer = network.StartTcp(tcp.configdata)
}
