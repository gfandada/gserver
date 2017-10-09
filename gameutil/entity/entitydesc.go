package entity

import (
	"reflect"
)

// TODO 目前只支持本地rpc
type rpcDesc struct {
	Func       reflect.Value // cb
	Flags      uint          // TODO 标记：用于区分fb执行粒度
	MethodType reflect.Type  // 主要是用于区分本地执行or远程执行...
	NumArgs    int           // 参数数量
}

type rpcDescMap map[string]*rpcDesc // key:方法名称

type EntityTypeDesc struct {
	isPersistent    bool         // 是否持久化
	useAOI          bool         // 是否具有aoi行为
	entityType      reflect.Type // entity对象类型
	rpcDescs        rpcDescMap   // handler容器
	allClientAttrs  StringSet
	clientAttrs     StringSet
	persistentAttrs StringSet
}

// TODO flag
func (rdm rpcDescMap) visit(method reflect.Method) {
	methodType := method.Type
	rdm[method.Name] = &rpcDesc{
		Func:       method.Func,
		Flags:      0,
		MethodType: methodType,
		NumArgs:    methodType.NumIn() - 1,
	}
}
