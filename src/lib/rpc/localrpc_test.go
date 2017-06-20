package rpc

import (
	"testing"
)

type AddParam struct {
	a int
	b int
}

func Test_localrpc(t *testing.T) {
	// new rpc server
	localServer := NewLocalServer(10)
	localServer.Start()
	// new rpc client
	client := localServer.NewLocalClient()
	ret := make(chan *OutputMessage, 1)
	// 注册两组消息，都是用add作为handler
	client.Server.Register(1, add)
	client.Server.Register(2, add)
	args := make([]interface{}, 1)
	args[0] = &AddParam{
		a: 10,
		b: 20,
	}
	// 同步调用消息1
	outMessage, err := client.Call(&InputMessage{
		Msg:        1,
		F:          add,
		Args:       args,
		OutputChan: ret,
	}, 2)
	if err != nil {
		t.Error(err)
	} else if outMessage.Ret != 30 {
		t.Error("should 30")
	}
	args[0] = &AddParam{
		a: 10,
		b: 80,
	}
	// 同步调用消息2
	outMessage1, err1 := client.Call(&InputMessage{
		Msg:        2,
		F:          add,
		Args:       args,
		OutputChan: ret,
	}, 2)
	if err1 != nil {
		t.Error(err1)
	} else if outMessage1.Ret != 70 {
		t.Error("should 90")
	}
}

func add(addParam []interface{}) interface{} {
	param := addParam[0].(*AddParam)
	ret := param.a + param.b
	return ret
}
