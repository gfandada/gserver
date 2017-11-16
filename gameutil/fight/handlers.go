// FIXME 非线程安全
package fight

type handler func(inner, args []interface{}) []interface{}

var (
	_handlers map[string]handler
)

func init() {
	_handlers = make(map[string]handler)
}

func RegisterHandler(msg string, hand handler) {
	_handlers[msg] = hand
}

func UnRegisterHandler(msg string) {
	delete(_handlers, msg)
}

func GetHandler(msg string) handler {
	return _handlers[msg]
}
