// gserver - demo
package main

import (
	"fmt"
	"lib/gservices"
	"lib/module"
	"lib/network"
	"lib/network/protobuff"
	"protomsg"
	"time"

	"github.com/golang/protobuf/proto"
)

// 登录请求处理器
func loginReqHandler(data []interface{}) []interface{} {

	/*****************这里是具体的登录逻辑*******************/

	/*
		1 如果处理完登录逻辑后，不需要做任何操作，直接return nil
		2 如果处理完登录逻辑后，需要回复一个自定义的pb数据给当前客户端，可以这样使用：
			msg := &protomsg.LoginAck{
				Str: proto.String("hello i am websocket server,hahhahaahahahahhah"),
			}
			ackMsg := protobuff.RawMessage{
				MsgId:   8888,
				MsgData: msg,
			}
			return []interface{}{ackMsg}
		3 如果处理完登录逻辑后，需要回复一个自定义的pb数据给当前客户端，同时更新session数据：
			msg := &protomsg.LoginAck{
				Str: proto.String("hello i am websocket server,hahhahaahahahahhah"),
			}
			ackMsg := protobuff.RawMessage{
				MsgId:   8888,
				MsgData: msg,
			}
			session := []network.SessionData{
			network.SessionData{
				Key:   "type1",
				Value: "type1_value",
			}, network.SessionData{
				Key:   "type2",
				Value: "type2_value",
			},
			userId := 1
			return []interface{}{ackMsg, userId, network.Update, session}
		}
	*/

	// 处理完了回复客户端
	msg := &protomsg.LoginAck{
		Str: proto.String("hello i am websocket server,hahhahaahahahahhah"),
	}
	ackMsg := protobuff.RawMessage{
		MsgId:   8888,
		MsgData: msg,
	}
	// 更新session数据，记录当前用户id
	session := []network.SessionData{
		network.SessionData{
			Key:   "type1",
			Value: "type1_value",
		}, network.SessionData{
			Key:   "type2",
			Value: "type2_value",
		}, network.SessionData{
			Key:   "conn",
			Value: data[1].(*network.Agent),
		},
	}
	return []interface{}{ackMsg, uint32(0), network.Update, session}
}

func createTCP() *network.Gate {
	gate := new(network.Gate)
	gate.MaxConnNum = 1000
	gate.PendingNum = 100
	gate.MaxMsgLen = 1024
	gate.ServerAddress = "localhost:9527"
	gate.LenMsgLen = 1024
	return gate
}

func createWS() *network.Gate {
	gate := new(network.Gate)
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
	gate := createWS()
	//gate := createTCP()
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
