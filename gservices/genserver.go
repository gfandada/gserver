// 封装了本地genserver，支持异步同步调用
package gservices

import (
	"errors"
	"fmt"
	"lib/logger"
	"sync"
	"time"
)

// 客户端发送的message
type InputMessage struct {
	Msg        interface{}         // FIXME 消息标识（暂时不预留内部消息）
	F          interface{}         // 消息handler
	CB         Iack                // 消息回调
	Args       []interface{}       // 函数调用参数
	OutputChan chan *OutputMessage // 接收返回值的队列（长度为1）
}

// 服务器返回的message
type OutputMessage struct {
	Err error         // 错误描述
	Ret []interface{} // 返回值
}

// server是一个容器，同时也负责执行handle，将结果发送给client
type LocalServer struct {
	mutex          sync.Mutex
	Functions      map[interface{}]interface{} // map[msg]msg_handler主要用于一些检查:server不需要处理一些无效的msg
	MessageBoxChan chan *InputMessage          // 消息队列
	Pending        int                         // 用于记录当前排队的消息数量
}

//
type LocalClient struct {
	Server *LocalServer // rpcserver
	//	ChanOutputSync chan *Output // 同步执行的结果队列
	//	ChanOutputAsyn chan *Output // 异步执行的结果队列
	//	PendingSync    int          // 同步调用时的队列长度
}

type MessageHandler1 func([]interface{})
type MessageHandler2 func([]interface{}) interface{}
type MessageHandler3 func([]interface{}) []interface{}

type MessageHandlerRet1 []interface{}
type MessageHandlerRet2 interface{}

// 新建一个rpcserver
func NewLocalServer(length int) *LocalServer {
	server := new(LocalServer)
	server.Functions = make(map[interface{}]interface{})
	server.MessageBoxChan = make(chan *InputMessage, length)
	return server
}

// 运行server
func (server *LocalServer) Start() {
	go func() {
		for {
			select {
			case inputMessage := <-server.MessageBoxChan:
				if inputMessage == nil {
					server.Pending--
					break
				}
				// 未注册的消息直接丢弃
				if server.Check(inputMessage) {
					server.Exec(inputMessage)
				}
				server.Pending--
			}
		}
	}()
}

// 强制关闭server
// 不再处理剩余的所有消息
func (server *LocalServer) CloseByForce() {
	// 先关闭，防止新写入
	close(server.MessageBoxChan)
	for inputMessage := range server.MessageBoxChan {
		server.ret(inputMessage, &OutputMessage{
			Err: errors.New("chanrpc server closed"),
		})
	}
}

// 优雅的关闭server
// 会处理完之前剩余的消息
func (server *LocalServer) CloseByGrace() {
	// 先关闭，防止新写入
	close(server.MessageBoxChan)
	for inputMessage := range server.MessageBoxChan {
		server.Exec(inputMessage)
	}
}

// 向server中注册一对值
func (server *LocalServer) Register(msg interface{}, msgHandler interface{}) error {
	server.mutex.Lock()
	defer server.mutex.Unlock()
	if _, ok := server.Functions[msg]; ok {
		return errors.New("multiple registration")
	}
	server.Functions[msg] = msgHandler
	server.Pending++
	return nil
}

// 执行方法体
func (server *LocalServer) Exec(input *InputMessage) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("genserver Exec panic error，%v:", r))
			server.ret(input, &OutputMessage{Err: fmt.Errorf("%v", r)})
		}
	}()
	switch input.F.(type) {
	case MessageHandler1:
		input.F.(MessageHandler1)(input.Args)
		server.ret(input, &OutputMessage{})
		return
	case MessageHandler2:
		ret := input.F.(MessageHandler2)(input.Args)
		server.ret(input, &OutputMessage{Ret: []interface{}{ret}})
		return
	case MessageHandler3:
		ret := input.F.(MessageHandler3)(input.Args)
		server.ret(input, &OutputMessage{Ret: ret})
		return
	}
	panic("error call function")
}

// 执行cb
func execCB(input *InputMessage, output *OutputMessage) {
	if output.Ret == nil {
		return
	}
	input.CB.Ack(output.Ret)
}

// 检查消息是否被注册
func (server *LocalServer) Check(input *InputMessage) bool {
	if _, ok := server.Functions[input.Msg]; !ok {
		return false
	}
	return true
}

// 将结果写进chan
func (server *LocalServer) ret(input *InputMessage, output *OutputMessage) {
	if input.OutputChan != nil {
		input.OutputChan <- output
		return
	}
	if input.CB != nil {
		execCB(input, output)
	}
}

// 新建一个rpcclient
func (server *LocalServer) NewLocalClient() *LocalClient {
	client := new(LocalClient)
	client.Server = server
	return client
}

// 同步请求rpcserver
// input:请求参数  timeout:超时时间（2*time.Second）
// FIXME 调用前需要先注册
func (client *LocalClient) Call(input *InputMessage, timeout int) (*OutputMessage, error) {
	client.Server.MessageBoxChan <- input
	select {
	case ret := <-input.OutputChan:
		return ret, nil
	case <-time.After(time.Duration(timeout) * time.Second):
		return nil, errors.New("time out")
	}
}

// 异步请求rpcserver
// FIXME 调用前需要先注册
func (client *LocalClient) Cast(input *InputMessage) {
	client.Server.MessageBoxChan <- input
}
