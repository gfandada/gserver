// 通用websocket网关
package gateway

import (
	"net"
	"time"

	Loader "github.com/gfandada/gserver/loader"
	"github.com/gfandada/gserver/network"
	"github.com/gorilla/websocket"
)

type Agent struct {
	configdata *network.Config // 配置数据
	conn       *websocket.Conn // 当前会话
}

func (wg *Agent) Start(conn *websocket.Conn) {
	defer conn.Close()
	wg.conn = conn
	config := wg.configdata
	in := make(chan []byte)
	defer close(in)
	var sess Session
	sess.Die = make(chan struct{})
	if sender := startSender(conn, &sess, in, config); sender == nil {
		close(sess.Die)
		return
	}
	for {
		conn.SetReadDeadline(time.Now().Add(time.Duration(config.ReadDeadline) * time.Second))
		conn.SetWriteDeadline(time.Now().Add(time.Duration(config.WriteDeadline) * time.Second))
		payload, err := config.Parser.ReadBody(conn)
		if err != nil {
			return
		}
		select {
		case in <- payload:
		case <-sess.Die:
			return
		}
	}
	wg.Close()
}

func (wg *Agent) NewIagent() network.Iagent {
	return &Agent{configdata: wg.configdata}
}

func (wg *Agent) Close() {

}

func (wg *Agent) WriteMsg(msg interface{}) {

}

func (wg *Agent) LocalAddr() net.Addr {
	return wg.conn.LocalAddr()
}

func (wg *Agent) RemoteAddr() net.Addr {
	return wg.conn.RemoteAddr()
}

func (wg *Agent) GetUserData() interface{} {
	return nil
}

func (wg *Agent) SetUserData(data interface{}) {

}

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
	wg.wsServer = network.Start(wg.configdata)
}
