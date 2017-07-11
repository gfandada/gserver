package db

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

/******
type Pool struct {
    // 用来创建redis连接的方法
    Dial func() (Conn, error)
    // 如果设置了给func,那么每次p.Get()的时候都会调用该方法来验证连接的可用性
    TestOnBorrow func(c Conn, t time.Time) error
    // 定义连接池中最大连接数（超过这个数会关闭老的链接，总会保持这个数）
    MaxIdle int
    // 当前连接池中可用的链接数
    MaxActive int
    // 定义链接的超时时间，每次p.Get()的时候会检测这个连接是否超时（超时会关闭，并释放可用连接数）
    IdleTimeout time.Duration
    // 当可用连接数为0是，那么当wait=true,那么当调用p.Get()时，会阻塞等待，否则，返回nil.
    Wait bool
}
******/

const (
	IMAGEKEY = "GSERVERKEY"
)

const (
	INFINITE = (1 << 32) * time.Second
)

type Cache struct {
	pool *redis.Pool // 连接池
	key  string      // 用于记录redis中所有的key
}

type Redis struct {
	MaxIdle            int
	MaxActive          int
	IdleTimeout        int
	RedisServer        string
	DialConnectTimeout int
	DialReadTimeout    int
	DialWriteTimeout   int
	Auth               string
	DbNum              int
}

// 新建redis-pool
func NewRedis(redisCfg Redis) *Cache {
	cache := &Cache{
		key: IMAGEKEY,
	}
	cache.initRedis(redisCfg)
	conn := cache.pool.Get()
	defer conn.Close()
	return cache
}

func (cache *Cache) initRedis(redisCfg Redis) {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", redisCfg.RedisServer)
		if err != nil {
			return nil, err
		}
		if redisCfg.Auth != "" {
			if _, err := c.Do("AUTH", redisCfg.Auth); err != nil {
				c.Close()
				return nil, err
			}
		}
		_, selecterr := c.Do("SELECT", redisCfg.DbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	var maxIdle, maxActive int
	var idleTimeout time.Duration
	if redisCfg.MaxIdle <= 0 {
		maxIdle = 3
	} else {
		maxIdle = redisCfg.MaxIdle
	}
	if redisCfg.MaxActive <= 0 {
		maxActive = 32
	} else {
		maxActive = redisCfg.MaxActive
	}
	if redisCfg.IdleTimeout <= 0 {
		idleTimeout = time.Duration(180) * time.Second
	} else {
		idleTimeout = time.Duration(redisCfg.IdleTimeout) * time.Second
	}
	cache.pool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
		Dial:        dialFunc,
	}
}

/*******************************封装调用接口*******************************/

// 执行redis命令
func (cache *Cache) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := cache.pool.Get()
	defer conn.Close()
	return conn.Do(commandName, args...)
}

// 获取指定key
func (cache *Cache) Get(key string) interface{} {
	if v, err := cache.do("GET", key); err == nil {
		return v
	}
	return nil
}

// 获取多个key
func (cache *Cache) GetMulti(keys []string) []interface{} {
	size := len(keys)
	var rv []interface{}
	conn := cache.pool.Get()
	defer conn.Close()
	var err error
	for _, key := range keys {
		err = conn.Send("GET", key)
		if err != nil {
			goto ERROR
		}
	}
	if err = conn.Flush(); err != nil {
		goto ERROR
	}
	for i := 0; i < size; i++ {
		if v, err := conn.Receive(); err == nil {
			rv = append(rv, v.([]byte))
		} else {
			rv = append(rv, err)
		}
	}
	return rv
ERROR:
	rv = rv[0:0]
	for i := 0; i < size; i++ {
		rv = append(rv, nil)
	}
	return rv
}

// 存储一对k-v
func (cache *Cache) Put(key string, val interface{}, timeout time.Duration) error {
	var err error
	if _, err = cache.do("SETEX", key, int64(timeout/time.Second), val); err != nil {
		return err
	}
	if _, err = cache.do("HSET", cache.key, key, true); err != nil {
		return err
	}
	return err
}

// 删除指定key
func (cache *Cache) Delete(key string) error {
	var err error
	if _, err = cache.do("DEL", key); err != nil {
		return err
	}
	_, err = cache.do("HDEL", cache.key, key)
	return err
}

// 检查指定key是否存在
func (cache *Cache) IsExist(key string) bool {
	v, err := redis.Bool(cache.do("EXISTS", key))
	if err != nil {
		return false
	}
	if !v {
		if _, err = cache.do("HDEL", cache.key, key); err != nil {
			return false
		}
	}
	return v
}

// 自增指定key
func (cache *Cache) Incr(key string) error {
	_, err := redis.Bool(cache.do("INCRBY", key, 1))
	return err
}

// 自减指定key
func (cache *Cache) Decr(key string) error {
	_, err := redis.Bool(cache.do("INCRBY", key, -1))
	return err
}

// 清理所有缓存
func (cache *Cache) ClearAll() error {
	cachedKeys, err := redis.Strings(cache.do("HKEYS", cache.key))
	if err != nil {
		return err
	}
	for _, str := range cachedKeys {
		if _, err = cache.do("DEL", str); err != nil {
			return err
		}
	}
	_, err = cache.do("DEL", cache.key)
	return err
}
