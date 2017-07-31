package cluster

import (
	"io"
	"net"

	"github.com/gfandada/gserver/cluster/pb"
	"github.com/gfandada/gserver/gservices"
	"github.com/gfandada/gserver/logger"
	"github.com/gfandada/gserver/network/protobuff"
	//	"github.com/gfandada/gserver/util"
	"google.golang.org/grpc"
)

type Service struct {
	name      string
	manager   *protobuff.MsgManager
	msgServer *gservices.LocalServer
}

// 启动一个远程rpc服务
// TODO 会优化成配置
func Start(name string, add string) *Service {
	listen, err := net.Listen("tcp", add)
	if err != nil {
		logger.Error("rpc server %s listen error", name)
	}
	s := new(Service)
	s.name = name
	s.manager = protobuff.NewMsgManager()
	s.msgServer = gservices.NewLocalServer(1024 * 2)
	serve := grpc.NewServer()
	pb.RegisterClusterServiceServer(serve, s)
	logger.Info("run cluster service {%s:%s}", name, add)
	go serve.Serve(listen)
	return s
}

// 注册msg(非线程安全)
// 先注册再使用
func (s *Service) Register(rawM protobuff.RawMessage, handler gservices.MessageHandler3) error {
	s.manager.RegisterMessage(rawM, handler, nil)
	s.msgServer.Register(rawM.MsgId, handler)
	return nil
}

// exec msg handler
func (s *Service) Exec(stream pb.ClusterService_RouterServer,
	msg *protobuff.RawMessage, die chan struct{}) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("genserver Exec input %v error: %v", msg, r)
		}
	}()
	exec := s.manager.MsgMap[msg.MsgId].MsgHandler
	if exec == nil {
		panic("no call function")
	}
	ret := make(chan *gservices.OutputMessage, 1)
	outMessage, err := s.msgServer.NewLocalClient().Call(&gservices.InputMessage{
		Msg:        msg.MsgId,
		F:          exec,
		Args:       []interface{}{msg.MsgData},
		OutputChan: ret,
	}, 2)
	if err != nil {
		logger.Error("cluster service {%s} exec {%d} error {%v}", s.name, msg.MsgId, err)
		return
	}
	s.ret(stream, outMessage.Ret, die)
}

// 处理返回值
func (s *Service) ret(stream pb.ClusterService_RouterServer, retD []interface{}, die chan struct{}) {
	switch len(retD) {
	// 同步ack
	case 1:
		s.send(stream, retD[0].(protobuff.RawMessage))
	// 同步ack，更新session --> [userid]chan
	case 2:
		s.ack(retD[1].(uint64), stream, die)
		s.send(stream, retD[0].(protobuff.RawMessage))
	}
}

/************************实现ClusterServiceClient接口*************************/

func (s *Service) Router(stream pb.ClusterService_RouterServer) error {
	die := make(chan struct{})
	defer func() {
		die <- struct{}{}
		close(die)
	}()
	recvChan := s.recv(stream, die)
	for {
		select {
		case data, ok := <-recvChan:
			if !ok {
				return nil
			}
			msg, err := s.manager.Deserialize(data.Data)
			if err == nil {
				logger.Debug("cluster service {%s} recv {%d:%v}", s.name,
					msg.MsgId, msg.MsgData)
				s.Exec(stream, msg, die)
			}
		// 收到关闭信号
		case <-die:
		}
	}
	return nil
}

func (s *Service) send(stream pb.ClusterService_RouterServer, data protobuff.RawMessage) error {
	dataNew, err := s.manager.Serialize(data)
	if err != nil {
		logger.Error("cluster service {%s} Serialize error : %v", s.name, err)
		return err
	}
	sendM := &pb.Message{
		Data: dataNew,
		Id:   uint32(data.MsgId),
	}
	logger.Debug("cluster service {%s} ack {%d:%v}", s.name,
		sendM.Id, sendM.Data)
	if err := stream.Send(sendM); err != nil {
		logger.Error("cluster service {%s} send error: %v", s.name, err)
		return err
	}
	return nil
}

func (s *Service) recv(stream pb.ClusterService_RouterServer, die chan struct{}) chan *pb.Message {
	recvChan := make(chan *pb.Message, 1)
	go func() {
		defer close(recvChan)
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				logger.Error("cluster service {%s} recv error : %v", s.name, err)
				return
			}
			select {
			case recvChan <- in:
				logger.Debug("cluster service {%s} recv {%d:%v} : %v", s.name,
					in.Id, in.Data)
			case <-die:
			}
		}
	}()
	return recvChan
}

func (s *Service) ack(id uint64, stream pb.ClusterService_RouterServer, die chan struct{}) {
	sendch := make(chan protobuff.RawMessage, 512)
	Register(id, sendch)
	go func() {
		defer close(sendch)
		for {
			select {
			case data := <-sendch:
				if err := s.send(stream, data); err != nil {
					logger.Error("cluster service {%s} ack error {%v}", s.name, err)
					return
				}
			case <-die:
			}
		}
	}()
}
