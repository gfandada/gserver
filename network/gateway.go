// 封装了网关的操作
package network

import (
	"io"
	"net"
	"time"

	"github.com/gfandada/gserver/cluster"
	"github.com/gfandada/gserver/cluster/pb"
	"github.com/gfandada/gserver/logger"
	"github.com/gfandada/gserver/network/protobuff"
	"github.com/gfandada/gserver/util"
)

// 定义网关结构体
type Gate struct {
	MaxConnNum       int      // 允许的最大的连接数
	PendingNum       int      // 最大发送队列长度（server -> client）
	MaxMsgLen        int      // 允许的服务器接收的最大的消息长度
	MessageProcessor Imessage // 用于消息体的处理
	ReadTimeout      int      // 读超时

	//	ServerAddress string // tcp服务地址
	//	LenMsgLen     int    // tcp消息长度

	WsServerAddress string        // websocket服务地址
	HTTPTimeout     time.Duration // http超时时间
	CertFile        string        // wws
	KeyFile         string        // wws
}

type Agent struct {
	Conn      Iconn                                     // 套接字套作接口
	Gate      *Gate                                     // 网关配置数据
	Stream    map[string]pb.ClusterService_RouterClient // 接受集群数据流
	UserData  interface{}                               // 用户数据
	HeartBeat int                                       // 心跳包次数
}

// tcp服务器
// var tcpServer *TcpServer

// ws服务器
var wsServer *WsServer

/***************************实现imodule接口*******************************/

// 启动网关
func (gate *Gate) Run(chClose chan bool) {
	switch {
	//	case tcpServer != nil:
	//		tcpServer.Start()
	//		<-chClose
	//		tcpServer.Close()
	case wsServer != nil:
		wsServer.Start()
		<-chClose
		wsServer.Close()
	}
}

// 初始化
func (gate *Gate) OnInit() {
	if gate == nil {
		logger.Error("tcp_gateway run failed, because gate is nil")
		return
	}
	switch {
	//	case gate.ServerAddress != "":
	//		tcpServer = new(TcpServer)
	//		tcpServer.ServerAddress = gate.ServerAddress
	//		tcpServer.MaxConnNum = gate.MaxConnNum
	//		tcpServer.PendingNum = gate.PendingNum
	//		tcpServer.Agent = func(conn *Conn) Iagent {
	//			arg := &Agent{Conn: conn, Gate: gate}
	//			return arg
	//		}
	case gate.WsServerAddress != "":
		wsServer = new(WsServer)
		wsServer.ServerAddress = gate.WsServerAddress
		wsServer.MaxConnNum = gate.MaxConnNum
		wsServer.MaxMsgLen = gate.MaxMsgLen
		wsServer.PendingNum = gate.PendingNum
		wsServer.HTTPTimeout = gate.HTTPTimeout
		wsServer.CertFile = gate.CertFile
		wsServer.KeyFile = gate.KeyFile
		wsServer.Agent = func(conn *WsConn) Iagent {
			arg := &Agent{Conn: conn, Gate: gate}
			return arg
		}
		wsServer.ReadTimeout = gate.ReadTimeout
	}
}

// 资源回收
func (gate *Gate) OnDestroy() {
	logger.Debug("gateway %d OnDestroy %v", util.GetPid(), gate)
}

/****************************实现了Iagent接口**********************************/

func (agent *Agent) Run() {
	if agent.Gate == nil {
		logger.Error("agent Run params is nil, %v", agent)
		return
	}
	die := make(chan struct{}, 1)
	agent.recv(die)
	for {
		msg, err := agent.Conn.ReadMsg()
		if err != nil {
			logger.Debug("agent run err:%v", err)
			break
		}
		if agent.Gate.MessageProcessor != nil {
			realMsg, errs := agent.Gate.MessageProcessor.Deserialize(msg)
			if errs != nil {
				logger.Error("deserialize err:%v", errs)
				break
			}
			logger.Debug("agent %v read msg %v", agent, realMsg)
			if err := agent.Gate.MessageProcessor.Router(realMsg, agent); err != nil {
				if err := agent.Stream["game"].Send(&pb.Message{
					Data: realMsg.MsgRaw,
				}); err != nil {
					logger.Error("msg route err:%v", errs)
					break
				}
			}
		}
	}
	if agent.Stream["game"] != nil {
		agent.Stream["game"].Send(&pb.Message{
			Id: cluster.CLOSEF,
		})
	}
	die <- struct{}{}
	logger.Debug("ws agent %d stop %v", util.GetPid(), agent)
}

func (agent *Agent) OnClose() {
	logger.Debug("agent OnClose:%v userdata:%v", agent, agent.UserData)
	if agent.UserData != nil {
		DeleteSessionConn(agent.UserData.(uint64))
	}
}

// 接收集群数据
func (agent *Agent) recv(die chan struct{}) {
	// 初始化集群流
	streams := cluster.GetRouterStreams()
	agent.Stream = streams
	for key := range agent.Stream {
		router := func() {
			for {
				data, err := agent.Stream[key].Recv()
				logger.Debug("recv cluster service {%s} stream {%v}", key, data)
				// 流关闭
				if err == io.EOF {
					logger.Info("recv cluster service {%s} stream closed", key)
					return
				}
				if err != nil {
					logger.Info("recv cluster service {%s} stream error %v", key, err)
					return
				}
				agent.WriteMsg(protobuff.RawMessage{
					MsgId:  uint16(data.Id), // 仅用于认证
					MsgRaw: data.Data,       // data
				})
				logger.Debug("cluster service recver {%s} ack client {%v}", key, agent.RemoteAddr())
				select {
				case <-die:
					return
				default:
				}
			}
		}
		go router()
		logger.Info("run cluster service {%s} stream recver", key)
	}
}

/****************************实现了Igateway接口**********************************/

func (agent *Agent) WriteMsg(msg protobuff.RawMessage) {
	if agent.Gate.MessageProcessor != nil {
		data, err := agent.Gate.MessageProcessor.Serialize(msg)
		if err != nil {
			logger.Error("Serialize message %v error: %v", msg, err)
			return
		}
		err = agent.Conn.WriteMsg(data)
		if err != nil {
			logger.Error("write message %v error: %v", msg, err)
		}
	}
}

func (agent *Agent) LocalAddr() net.Addr {
	if agent.Conn != nil {
		return agent.Conn.LocalAddr()
	}
	return nil
}

func (agent *Agent) RemoteAddr() net.Addr {
	if agent.Conn != nil {
		return agent.Conn.RemoteAddr()
	}
	return nil
}

func (agent *Agent) Close() {
	if agent.Conn != nil {
		agent.Conn.Close()
	}
}

func (agent *Agent) Destroy() {
	if agent.Conn != nil {
		agent.Conn.Destroy()
	}
}

func (agent *Agent) GetUserData() interface{} {
	return agent.UserData
}

func (agent *Agent) SetUserData(data interface{}) {
	agent.UserData = data
}

/******************************实现了Iack接口*********************************/

// TODO 需要优化
func (agent *Agent) Ack(data []interface{}) {
	if data == nil {
		return
	}
	switch len(data) {
	case 1:
		// ack自己
		agent.WriteMsg(data[0].(protobuff.RawMessage))
		logger.Debug("agent ack %v message %v", agent, data[0])
		return
	case 3:
		// ack自己
		agent.WriteMsg(data[0].(protobuff.RawMessage))
		// 更新session
		AddSessionConn(data[1].(uint64), data[2].(*Agent))
		agent.SetUserData(data[1].(uint64))
		logger.Debug("agent ack %v message %v", agent, data[0])
		return
	}
}
