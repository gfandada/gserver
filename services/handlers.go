// 简单的handler管理器
// 非线程安全
package services

type handler func([]interface{}) []interface{}

type MsgHandler struct {
	MsgId      uint16
	MsgHandler handler
}

var (
	_handlers map[uint16]handler
)

func init() {
	_handlers = make(map[uint16]handler)
}

func Register(id uint16, hand handler) {
	_handlers[id] = hand
}

func UnRegister(id uint16) {
	delete(_handlers, id)
}

func GetHandler(id uint16) handler {
	return _handlers[id]
}
