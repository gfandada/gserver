// 定义了消息的接口
package network

type Imessage interface {
	SetMaxLen(max int)                            // 设置上限
	Register(msg *RawMessage) error               // 消息注册
	UnRegister(msg *RawMessage)                   // 消息反注册
	Serialize(msg RawMessage) ([]byte, error)     // 序列化消息
	Deserialize(date []byte) (*RawMessage, error) // 反序列化消息
}
