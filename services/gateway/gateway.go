// 通用websocket网关
package gateway

import (
	"net"
	"time"

	Services "../"
	"../../network"
	"github.com/gorilla/websocket"
)

type WsGateway struct {
	Config string // 配置文件
}

func (wg *WsGateway) Start(conn *websocket.Conn, config *network.Config) {
	defer conn.Close()
	in := make(chan []byte)
	defer close(in)
	var sess Session
	sess.Die = make(chan struct{})
	if sender := startSender(conn, &sess, in, config); sender == nil {
		close(sess.Die)
		return
	}
	for {
		conn.SetReadDeadline(time.Now().Add(config.ReadDeadline))
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

func (wg *WsGateway) WriteMsg(msg interface{}) {}

func (wg *WsGateway) LocalAddr() net.Addr {
	return nil
}

func (wg *WsGateway) RemoteAddr() net.Addr {
	return nil
}

func (wg *WsGateway) GetUserData() interface{} {
	return nil
}

func (wg *WsGateway) SetUserData(data interface{}) {}

/****************************实现imodule接口******************************/

func (wg *WsGateway) OnInit() {

}

func (wg *WsGateway) OnDestroy() {

}

func (wg *WsGateway) Run(ChClose chan bool) {
	ws := new(network.WsServer)
	ws.ServerAddress = "localhost:9527"
	ws.MaxConnNum = 100
	ws.Gate = new(WsGateway)
	ws.MsgParser = network.NewMsgManager()
	// 注册消息
	ws.MsgParser.Register(&network.RawMessage{
		MsgId:   uint16(1002),
		MsgData: &Services.TouristsLoginReq{},
	})
	ws.MsgParser.Register(&network.RawMessage{
		MsgId:   uint16(1003),
		MsgData: &Services.TouristsLoginAck{},
	})
	ws.Start()
}
