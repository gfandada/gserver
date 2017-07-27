// 线程安全且管理高效
// 未添加连接心跳
package connpool

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gfandada/gserver/logger"
)

// 配置
type PoolConfig struct {
	// 连接池中拥有的最小连接数
	MinCap int
	// 连接池中拥有的最大的连接数
	MaxCap int
	// 生成连接的方法
	Factory func() (interface{}, error)
	// 关闭连接的方法
	Close func(interface{}) error
	// 连接最大空闲时间，超时将失效
	IdleTimeout time.Duration
}

// 连接池内部结构
type channelPool struct {
	mu          sync.Mutex
	conns       chan *idleConn
	factory     func() (interface{}, error)
	close       func(interface{}) error
	idleTimeout time.Duration
}

// 连接内部结构
type idleConn struct {
	conn interface{}
	t    time.Time
}

var (
	CLOSEPOOL = errors.New("pool closed")
)

// 新建
func NewChannelPool(poolConfig *PoolConfig) (Ipool, error) {
	if poolConfig.MinCap < 0 || poolConfig.MaxCap <= 0 || poolConfig.MinCap > poolConfig.MaxCap {
		logger.Error("poolConfig error %v", poolConfig)
		return nil, errors.New("poolConfig error")
	}
	c := &channelPool{
		conns:       make(chan *idleConn, poolConfig.MaxCap),
		factory:     poolConfig.Factory,
		close:       poolConfig.Close,
		idleTimeout: poolConfig.IdleTimeout,
	}
	for i := 0; i < poolConfig.MinCap; i++ {
		conn, err := c.factory()
		if err != nil {
			c.Release()
			return nil, fmt.Errorf("factory error: %s", err)
		}
		c.conns <- &idleConn{conn: conn, t: time.Now()}
	}
	return c, nil
}

// 获取所有连接
func (c *channelPool) getConns() chan *idleConn {
	c.mu.Lock()
	conns := c.conns
	c.mu.Unlock()
	return conns
}

// 从pool中取一个连接
func (c *channelPool) Get() (interface{}, error) {
	conns := c.getConns()
	if conns == nil {
		return nil, CLOSEPOOL
	}
	for {
		select {
		case conn := <-conns:
			if conn == nil {
				return nil, CLOSEPOOL
			}
			if timeout := c.idleTimeout; timeout > 0 {
				if conn.t.Add(timeout).Before(time.Now()) {
					c.Close(conn.conn)
					continue
				}
			}
			return conn.conn, nil
		default:
			conn, err := c.factory()
			if err != nil {
				return nil, err
			}
			return conn, nil
		}
	}
}

// 回收连接
func (c *channelPool) Put(conn interface{}) error {
	if conn == nil {
		return errors.New("conn is nil")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conns == nil {
		return c.Close(conn)
	}
	select {
	case c.conns <- &idleConn{conn: conn, t: time.Now()}:
		return nil
	// 已满
	default:
		return c.Close(conn)
	}
}

// 关闭指定连接
func (c *channelPool) Close(conn interface{}) error {
	return c.close(conn)
}

// 释放连接池中所有链接
func (c *channelPool) Release() {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	closeFun := c.close
	c.close = nil
	c.mu.Unlock()
	if conns == nil {
		return
	}
	close(conns)
	for conn := range conns {
		closeFun(conn.conn)
	}
}

// 获取连接池中已有的连接
func (c *channelPool) Len() int {
	return len(c.getConns())
}
