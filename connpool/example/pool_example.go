package main

import (
	"fmt"
	"net"
	"time"

	"github.com/gfandada/gserver/connpool"
)

func main() {
	factory := func() (interface{}, error) { return net.Dial("tcp", "localhost:9527") }
	close := func(v interface{}) error { return v.(net.Conn).Close() }
	poolConfig := &connpool.PoolConfig{
		MinCap:      5,
		MaxCap:      30,
		Factory:     factory,
		Close:       close,
		IdleTimeout: 15 * time.Second,
	}
	p, err := connpool.NewChannelPool(poolConfig)
	if err != nil {
		fmt.Println(err)
	}
	//从连接池中取得一个链接
	v, err := p.Get()
	fmt.Println("------------------------papapapa--------------------------")
	//将连接放回连接池中
	p.Put(v)
	//查看当前链接中的数量
	leng := p.Len()
	fmt.Println(leng)
	//释放连接池中的所有链接
	p.Release()
}
