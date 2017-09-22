// 内部错误封装
// base on protobuff
package services

import (
	"encoding/binary"

	"github.com/gfandada/gserver/network"
	"github.com/golang/protobuf/proto"
)

// 构建一个gataway通用内部错误
// @params err:错误描述
func NewInError(err error) []byte {
	return newError(0, err.Error())
}

// 构建一个gataway通用业务错误
// @params id:错误码
func NewLogicError(id int) []byte {
	return newError(id, "")
}

// 构建一个service通用内部错误(错误码为1000)
// @params err:错误描述
func NewSInError(err error) *network.Data_Frame {
	data := newError(0, err.Error())
	return &network.Data_Frame{
		Type:    network.Data_Message,
		Message: data,
	}
}

// 构建一个service通用业务错误
// @params id:错误码
func NewSLogicError(id int) *network.Data_Frame {
	data := newError(id, "")
	return &network.Data_Frame{
		Type:    network.Data_Message,
		Message: data,
	}
}

func newError(id int, str string) []byte {
	rawId := make([]byte, 2)
	binary.BigEndian.PutUint16(rawId, 2)
	data, err := proto.Marshal(&ErrorAck{
		Errid:  proto.Int32(int32(id)),
		Errstr: proto.String(str),
	})
	if err != nil {
		return nil
	}
	c := make([]byte, 2+len(data))
	copy(c, rawId)
	copy(c[len(rawId):], data)
	return c
}

func newServiceError(id int, str string) []byte {
	rawId := make([]byte, 2)
	binary.BigEndian.PutUint16(rawId, 2)
	data, err := proto.Marshal(&ErrorAck{
		Errid:  proto.Int32(int32(id)),
		Errstr: proto.String(str),
	})
	if err != nil {
		return nil
	}
	c := make([]byte, 2+len(data))
	copy(c, rawId)
	copy(c[len(rawId):], data)
	return c
}
