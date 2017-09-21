// 通用服务
package service

import (
	"net"

	Loader "../../loader"
	"../../network"
	"google.golang.org/grpc"
)

func Start(addr string, coder network.Imessage) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	ser := grpc.NewServer()
	ins := new(Agent)
	ins.msgParser = coder
	network.RegisterServiceServer(ser, ins)
	go ser.Serve(lis)
}

/****************************实现imodule接口******************************/

type Service struct {
	Config     string
	Coder      network.Imessage
	configdata *network.Config
}

func (s *Service) OnInit() {
	Loader.LoadJson(s.Config, s.configdata)
}

func (s *Service) OnDestroy() {

}

func (s *Service) Run(ChClose chan bool) {
	Start(s.configdata.ServerAddress, s.Coder)
}
