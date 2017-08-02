package cluster

import (
	"io"
	"net"

	"github.com/gfandada/gserver/cluster/pb"
	"github.com/gfandada/gserver/logger"
	"github.com/gfandada/gserver/network/protobuff"
	"google.golang.org/grpc"
)

type Service struct {
	Name    string
	Manager *protobuff.MsgManager
}

type ServiceAck struct {
	Stream   pb.ClusterService_RouterServer
	Die      chan struct{}
	UserData interface{}
	Service  *Service
	SendCh   chan protobuff.RawMessage
}

func Start(name string, add string, code *protobuff.MsgManager) {
	listen, err := net.Listen("tcp", add)
	if err != nil {
		logger.Error("service {%s} listen error {%v}", name, err)
	}
	s := new(Service)
	s.Name = name
	s.Manager = code
	serve := grpc.NewServer()
	pb.RegisterClusterServiceServer(serve, s)
	logger.Info("service {%s:%s} run", name, add)
	go serve.Serve(listen)
}

/************************实现ClusterServiceClient接口*************************/

func (s *Service) Router(stream pb.ClusterService_RouterServer) error {
	die := make(chan struct{})
	defer func() {
		die <- struct{}{}
		close(die)
	}()
	recvChan := s.recv(stream, die)
	ack := &ServiceAck{
		Stream:  stream,
		Die:     die,
		Service: s,
	}
	for {
		select {
		case data, ok := <-recvChan:
			if !ok {
				return nil
			}
			msg, err := s.Manager.Deserialize(data.Data)
			if err == nil {
				logger.Debug("cluster service {%s} recv {%d:%v}", s.Name,
					msg.MsgId, msg.MsgData)
				// 异步路由
				s.Manager.Router(msg, ack)
			}
		// 收到关闭信号
		case <-die:
		}
	}
	return nil
}

func (s *Service) send(stream pb.ClusterService_RouterServer, data protobuff.RawMessage) error {
	dataNew, err := s.Manager.Serialize(data)
	if err != nil {
		logger.Error("cluster service {%s} Serialize error : %v", s.Name, err)
		return err
	}
	sendM := &pb.Message{
		Data: dataNew,
		Id:   uint32(data.MsgId),
	}
	logger.Debug("cluster service {%s} ack {%d:%v}", s.Name,
		sendM.Id, sendM.Data)
	if err := stream.Send(sendM); err != nil {
		logger.Error("cluster service {%s} send error: %v", s.Name, err)
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
				logger.Error("cluster service {%s} recv error : %v", s.Name, err)
				return
			}
			select {
			case recvChan <- in:
				logger.Debug("cluster service {%s} recv {%d:%v}", s.Name,
					in.Id, in.Data)
			case <-die:
			}
		}
	}()
	return recvChan
}

func (s *Service) ack(id uint64, sack *ServiceAck) {
	sack.SendCh = make(chan protobuff.RawMessage, 512)
	Register(id, sack)
	go func() {
		defer close(sack.SendCh)
		for {
			select {
			case data := <-sack.SendCh:
				if err := s.send(sack.Stream, data); err != nil {
					logger.Error("cluster service {%s} ack error {%v}", s.Name, err)
					return
				}
			case <-sack.Die:
			}
		}
	}()
}

/***********************************实现Iack接口************************************/

func (s *ServiceAck) Ack(data []interface{}) {
	switch len(data) {
	// 同步ack
	case 1:
		s.Service.send(s.Stream, data[0].(protobuff.RawMessage))
	// 同步ack，更新session --> [userid]chan
	case 2:
		s.Service.ack(data[1].(uint64), s)
		s.UserData = data[1].(uint64)
		s.Service.send(s.Stream, data[0].(protobuff.RawMessage))
	}
}
