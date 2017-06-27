package main

import (
	"fmt"
	"gateway"
	"lib/gservices"
	"lib/module"
	"lib/network/protobuff"
	"lib/util"
	"protomsg"
	"time"

	"github.com/golang/protobuf/proto"
)

func loginReqHandler(data []interface{}) {
	fmt.Println("------------------处理登录逻辑-------------------------", util.GetPid())
	fmt.Println("用户发送的数据：", data[0].(*protomsg.LoginReq).GetId(), data[0].(*protomsg.LoginReq).GetStr())
	// 处理完了回复客户端
	msg := &protomsg.LoginAck{
		Str: proto.String("hello i am websocket server,hahhahaahahahahhah"),
	}
	ackMsg := protobuff.RawMessage{
		MsgId:   8888,
		MsgData: msg,
	}
	fmt.Println("------------------处理完登录逻辑，回复客户端-------------------------", data[1].(*gateway.Agent).Conn.LocalAddr(),
		data[1].(*gateway.Agent).RemoteAddr())
	data[1].(*gateway.Agent).WriteMsg(ackMsg)
}

func createTCP() *gateway.Gate {
	gate := new(gateway.Gate)
	gate.MaxConnNum = 1000
	gate.PendingNum = 100
	gate.MaxMsgLen = 1024
	gate.ServerAddress = "localhost:9527"
	gate.LenMsgLen = 1024
	return gate
}

func createWS() *gateway.Gate {
	gate := new(gateway.Gate)
	gate.MaxConnNum = 1000
	gate.PendingNum = 100
	gate.MaxMsgLen = 1024
	gate.WsServerAddress = "localhost:9527"
	gate.HTTPTimeout = 5 * time.Second
	return gate
}

func registerLoginMsg(code *protobuff.MsgManager, loginServer *gservices.LocalServer) {
	// 注册loginreq消息
	loginServer.Register(uint16(6780), loginReqHandler)
	code.RegisterMessage(protobuff.RawMessage{
		MsgId:   6780,
		MsgData: &protomsg.LoginReq{},
	}, loginReqHandler, loginServer)
	// 注册loginack消息（消息管理器只处理被注册的消息，本次注册只用于被检查）
	code.RegisterMessage(protobuff.RawMessage{
		MsgId:   8888,
		MsgData: &protomsg.LoginAck{},
	}, loginReqHandler, loginServer)
}

func main() {
	// 新建一个简易服务器
	//gate := createWS()
	gate := createTCP()
	// 采用pb作为数据序列化方法
	code := protobuff.NewMsgManager()
	// 注册一个登录服务器，允许同时1000个登录
	loginServer := gservices.NewLocalServer(1000)
	// 注册消息
	registerLoginMsg(code, loginServer)
	// 启动登录服务器
	loginServer.Start()
	gate.MessageProcessor = code
	fmt.Println("启动网关：", gate)
	// 启动网关模块
	module.Run(gate)
}
