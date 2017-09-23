package network

type Iagent interface {
	NewIagent() Iagent            // 拷贝构造器
	Close()                       // 关闭
	Start(Iconn)                  // 启动
	GetUserData() interface{}     // 获取用户数据
	SetUserData(data interface{}) // 设置用户数据
}
