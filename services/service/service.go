// 通用服务
package service

import (
	"net"

	"../../network"
	"google.golang.org/grpc"
)

func Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	ins := new(Agent)
	network.RegisterServiceServer(s, ins)
	return s.Serve(lis)
}

/****************************实现imodule接口******************************/

type Service struct {
	Config string // 配置
}

// 服务注册
func (s *Service) OnInit() {

}

// 服务销毁
func (s *Service) OnDestroy() {

}

func (s *Service) Run(ChClose chan bool) {
	// TODO DEMO
	lis, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		return
	}
	ser := grpc.NewServer()
	ins := new(Agent)
	ins.msgParser = network.NewMsgManager()
	network.RegisterServiceServer(ser, ins)
	go ser.Serve(lis)
}
