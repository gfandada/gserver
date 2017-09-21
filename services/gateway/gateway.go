// 通用websocket网关
package gateway

import (
	"net"
	"time"

	Loader "github.com/gfandada/gserver/loader"
	"github.com/gfandada/gserver/network"
	"github.com/gorilla/websocket"
)

type WsGateway struct {
	Config     string            // 配置文件
	Coder      network.Imessage  // 编码器
	configdata *network.Config   // 配置数据
	wsServer   *network.WsServer // 服务器
	conn       *websocket.Conn   // 当前会话
}

func (wg *WsGateway) Start(conn *websocket.Conn) {
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

func (wg *WsGateway) Close() {

}

func (wg *WsGateway) WriteMsg(msg interface{}) {

}

func (wg *WsGateway) LocalAddr() net.Addr {
	return wg.conn.LocalAddr()
}

func (wg *WsGateway) RemoteAddr() net.Addr {
	return wg.conn.RemoteAddr()
}

func (wg *WsGateway) GetUserData() interface{} {
	return nil
}

func (wg *WsGateway) SetUserData(data interface{}) {

}

/****************************实现imodule接口******************************/

func (wg *WsGateway) OnInit() {
	config := new(network.Config)
	Loader.LoadJson(wg.Config, config)
	wg.configdata = config
	wg.configdata.Gate = new(WsGateway)
	wg.configdata.MsgParser = network.NewMsgManager()
}

func (wg *WsGateway) OnDestroy() {
	wg.wsServer.Close()
}

func (wg *WsGateway) Run(ChClose chan bool) {
	wg.wsServer = network.Start(wg.configdata)
}
