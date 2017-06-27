// 定义了消息的接口
package network

import (
	"lib/network/protobuff"
)

type Imessage interface {
	Serialize(msg protobuff.RawMessage) ([][]byte, error)         // 序列化消息
	Deserialize(date []byte) (*protobuff.RawMessage, error)       // 反序列化消息
	Router(msg *protobuff.RawMessage, userData interface{}) error // 消息路由
}
