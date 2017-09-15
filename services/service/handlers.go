// handler管理器
// 非线程安全
package service

type handler func([]interface{}) []interface{}

var (
	handlers map[uint16]handler
)

func init() {
	handlers = make(map[uint16]handler)
}

func RegisterHandler(key uint16, value handler) {
	handlers[key] = value
}

func GetHandler(key uint16) handler {
	return handlers[key]
}

func IsRegister(key uint16) bool {
	_, ok := handlers[key]
	return ok
}
