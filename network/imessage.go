// 定义了消息的接口
package network

type Imessage interface {
	NewIMessage() Imessage                        // 拷贝构造器
	Register(msg *RawMessage) error               // 消息注册
	UnRegister(msg *RawMessage)                   // 消息反注册
	Serialize(msg RawMessage) ([]byte, error)     // 序列化消息
	Deserialize(date []byte) (*RawMessage, error) // 反序列化消息
}
