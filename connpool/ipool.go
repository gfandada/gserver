package connpool

type Ipool interface {
	Get() (interface{}, error) // 获取一个连接
	Put(interface{}) error     // 归还一个连接
	Close(interface{}) error   // 关闭一个连接
	Release()                  // 释放连接池
	Len() int                  // 获取连接池的容量
}
