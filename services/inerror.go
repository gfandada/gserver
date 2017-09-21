// 内部错误封装
package services

import (
	"github.com/gfandada/gserver/network"
)

// 构建一个gataway通用内部错误(错误码为0)
// @params err:错误描述
func NewInError(err error) []byte {
	return nil
}

// 构建一个gataway通用业务错误
// @params id:错误码
func NewLogicError(id int) []byte {
	return nil
}

// 构建一个service通用内部错误(错误码为1000)
// @params err:错误描述
func NewSInError(err error) *network.Data_Frame {
	return nil
}

// 构建一个service通用业务错误
// @params id:错误码
func NewSLogicError(id int) *network.Data_Frame {
	return nil
}
