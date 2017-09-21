// 通用服务
package service

import (
	"net"

	Loader "github.com/gfandada/gserver/loader"
	"github.com/gfandada/gserver/network"
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
	config := new(network.Config)
	Loader.LoadJson(s.Config, config)
	s.configdata = config
}

func (s *Service) OnDestroy() {

}

func (s *Service) Run(ChClose chan bool) {
	Start(s.configdata.ServerAddress, s.Coder)
}
