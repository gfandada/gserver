package chanrpc

import (
	"errors"
	"fmt"
)

// 输出参数
type RetInfo struct {
	Ret interface{} // nil || interface{} || []interface{} 回调需要的args
	Err error       // 回调需要的error
	CB  interface{} // func(err error) || func(ret interface{}, err error) || func(ret []interface{}, err error)
}

// 输入参数
type CallInfo struct {
	F       interface{}   // 调用的方法：这个方法是通过kv的方式注册在Server.Functions中
	Args    []interface{} // 调用方法时传入的参数
	ChanRet chan *RetInfo // 调用后的返回值:其实就是Client中定义的队列
	CB      interface{}   // 需要传递给RetInfo的回调函数
}

// server数据结构
type Server struct {
	Functions map[interface{}]interface{} // 将func注册进server中，key:消息 value:消息处理函数
	Chancall  chan *CallInfo              // func的输入/输出数据
}

// client数据结构
// client的存在就是为了消费server
type Client struct {
	Server          *Server       // 目标server标识
	ChanSyncRet     chan *RetInfo // 同步chan：其实是已经执行完后的一个结果队列
	ChanAsynRet     chan *RetInfo // 异步chan：其实是已经执行完后的一个结果队列
	PendingAsynCall int           // 同步调用时的pending值
}

// 新建一个rpc服务器，需要指定通道容量
func NewServer(length int) *Server {
	server := new(Server)
	server.Functions = make(map[interface{}]interface{})
	server.Chancall = make(chan *CallInfo, length)
	return server
}

// 向rpc服务器注册方法
// FIXME calling Register before Open and Go
func (server *Server) Register(key interface{}, value interface{}) {
	if server == nil || server.Functions == nil {
		return
	}
	if _, ok := server.Functions[key]; ok {
		fmt.Println("已经被注册了:", key, value)
		return
	}
	server.Functions[key] = value
}

// 维护callInfo，RetInfo的关系(输入输出)
func (server *Server) ret(call *CallInfo, ret *RetInfo) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ret error:", r)
		}
	}()
	// 执行完了，再把开始注册进来的回调，放在ret上
	ret.CB = call.CB
	call.ChanRet <- ret
	return nil
}

// 执行方法体
func (server *Server) exec(call *CallInfo) error {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println("ret error:", r)
		}
		err := fmt.Errorf("%v", r)
		server.ret(call, &RetInfo{Err: err})
	}()
	// 执行方法，并维护输入输出的关系
	switch call.F.(type) {
	case func([]interface{}):
		call.F.(func([]interface{}))(call.Args)
		return server.ret(call, &RetInfo{})
	case func([]interface{}) interface{}:
		ret := call.F.(func([]interface{}) interface{})(call.Args)
		return server.ret(call, &RetInfo{Ret: ret})
	case func([]interface{}) []interface{}:
		ret := call.F.(func([]interface{}) []interface{})(call.Args)
		return server.ret(call, &RetInfo{Ret: ret})
	}
	panic("error call function")
}

// 执行方法体
func (server *Server) Exec(callInfo *CallInfo) {
	if err := server.exec(callInfo); err != nil {
		fmt.Println("Exec error:", err)
	}
}

// 根据key,设置对应func的callinfo
func (server *Server) Go(key interface{}, args ...interface{}) {
	value := server.Functions[key]
	if value == nil {
		fmt.Println("server key no register:", key)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("server go recover:", key)
		}
	}()
	// 设置相关数据
	server.Chancall <- &CallInfo{
		F:    value,
		Args: args,
	}
}

// 基于rpc-server,新建一个rpc-client
func (server *Server) Open(length int) *Client {
	client := NewClient(length)
	client.Attach(server)
	return client
}

// 关闭rpc-server
func (server *Server) Close() {
	// 先关闭，防止新写入
	close(server.Chancall)
	for value := range server.Chancall {
		server.ret(value, &RetInfo{
			Err: errors.New("chanrpc server closed"),
		})
	}
}

func (server *Server) call0(key interface{}, args ...interface{}) error {
	return server.Open(0).call0(key, args...)
}

func (server *Server) call1(key interface{}, args ...interface{}) (interface{}, error) {
	return server.Open(0).call1(key, args...)
}

func (server *Server) call2(key interface{}, args ...interface{}) ([]interface{}, error) {
	return server.Open(0).call2(key, args...)
}

/*******************************rpc-client操作*********************************/

// 新建一个rpc客户端
func NewClient(length int) *Client {
	client := new(Client)
	client.ChanSyncRet = make(chan *RetInfo, length)
	client.ChanAsynRet = make(chan *RetInfo, length)
	return client
}

// 关联rpc-server
func (client *Client) Attach(server *Server) {
	client.Server = server
}

// call0函数调用的实体
func (client *Client) call0(key interface{}, args ...interface{}) (err error) {
	// 获取具体的方法
	fun, err := client.function(key, 0)
	if err == nil {
		return
	}
	if err = client.call(&CallInfo{
		F:       fun,
		Args:    args,
		ChanRet: client.ChanSyncRet,
	}, true); err != nil {
		return
	}
	ret := <-client.ChanSyncRet
	return ret.Err
}

// call1函数调用的实体
func (client *Client) call1(key interface{}, args ...interface{}) (fun interface{}, err error) {
	// 获取具体的方法
	fun, err = client.function(key, 0)
	if err == nil {
		return
	}
	if err = client.call(&CallInfo{
		F:       fun,
		Args:    args,
		ChanRet: client.ChanSyncRet,
	}, true); err != nil {
		return
	}
	ret := <-client.ChanSyncRet
	return fun, ret.Err
}

// call2函数调用的实体
func (client *Client) call2(key interface{}, args ...interface{}) ([]interface{}, error) {
	// 获取具体的方法
	fun, err := client.function(key, 2)
	if err == nil {
		return nil, err
	}
	if err = client.call(&CallInfo{
		F:       fun,
		Args:    args,
		ChanRet: client.ChanSyncRet,
	}, true); err != nil {
		return nil, err
	}
	ret := <-client.ChanSyncRet
	return assert(ret.Ret), ret.Err
}

// 本方法作用是通过key从server中获取被注册的真实函数
func (client *Client) function(key interface{}, n int) (fun interface{}, err error) {
	if client.Server == nil {
		err = errors.New("server not attach")
		return
	}
	// 获取已被注册的方法
	fun = client.Server.Functions[key]
	if fun == nil {
		err = fmt.Errorf("function id %v: function not registered", key)
		return
	}
	// 执行函数体
	// TODO
	var ok bool
	switch n {
	case 0:
		_, ok = fun.(func([]interface{}))
	case 1:
		_, ok = fun.(func([]interface{}) interface{})
	case 2:
		_, ok = fun.(func([]interface{}) []interface{})
	default:
		panic("bug")
	}
	if !ok {
		err = fmt.Errorf("function key %v: return type mismatch", key)
	}
	return
}

// 将参数数据写入server的chan中
// block=true同步  block=false异步
func (client *Client) call(callInfo *CallInfo, block bool) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	if block {
		client.Server.Chancall <- callInfo
	} else {
		select {
		case client.Server.Chancall <- callInfo:
		default:
			err = errors.New("chanrpc channel full")
		}
	}
	return
}

// 本接口只是负责往异步队列中注册一个方法
func (client *Client) AsynCall(key interface{}, args ...interface{}) {
	// 异常检查
	if len(args) < 1 {
		panic("callback function not found")
	}
	argsNew := args[:len(args)-1]
	cb := args[len(args)-1]
	var n int
	switch cb.(type) {
	case func(error):
		n = 0
	case func(interface{}, error):
		n = 1
	case func([]interface{}, error):
		n = 2
	default:
		panic("definition of callback function is invalid")
	}
	if client.PendingAsynCall >= cap(client.ChanAsynRet) {
		fmt.Println("too many calls")
		return
	}
	// 写入异步队列中
	client.asynCall(key, argsNew, cb, n)
	// 更新异步队列pending梳理
	client.PendingAsynCall++
}

func (client *Client) asynCall(key interface{}, args []interface{}, cb interface{}, n int) {
	f, err := client.function(key, n)
	if err != nil {
		// 将失败的也记录到ChanAsynRet队列中
		client.ChanAsynRet <- &RetInfo{
			Err: err,
			CB:  cb,
		}
		return
	}
	err = client.call(&CallInfo{
		F:       f,
		Args:    args,
		ChanRet: client.ChanAsynRet,
		CB:      cb,
	}, false)
	if err != nil {
		client.ChanAsynRet <- &RetInfo{
			Err: err,
			CB:  cb,
		}
		return
	}
}

// 执行异步回调函数
func (client *Client) ExecCB(retInfo *RetInfo) {
	client.PendingAsynCall--
	client.execCB(retInfo)
}

func (client *Client) execCB(retInfo *RetInfo) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("execCB error:", r)
		}
	}()
	switch retInfo.CB.(type) {
	case func(error):
		retInfo.CB.(func(error))(retInfo.Err)
	case func(interface{}, error):
		retInfo.CB.(func(interface{}, error))(retInfo.Ret, retInfo.Err)
	case func([]interface{}, error):
		retInfo.CB.(func([]interface{}, error))(assert(retInfo.Ret), retInfo.Err)
	default:
		panic("bug")
	}
	return
}

// 空闲：true  不空闲：false
func (client *Client) Idle() bool {
	return client.PendingAsynCall == 0
}

// 关闭rpc-client
func (client *Client) Close() {
	for client.PendingAsynCall > 0 {
		client.ExecCB(<-client.ChanAsynRet)
	}
}

func assert(i interface{}) []interface{} {
	if i == nil {
		return nil
	} else {
		return i.([]interface{})
	}
}
