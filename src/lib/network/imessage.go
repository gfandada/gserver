// 定义了消息的接口
package network

type Imessage interface {
	Serialize(msg interface{}) ([]byte, error)          // 序列化消息
	Deserialize(date []byte) (interface{}, error)       // 反序列化消息
	Router(msg interface{}, userData interface{}) error // 消息路由
}
