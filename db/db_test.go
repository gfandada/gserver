package db

import (
	"fmt"
	"testing"

	"github.com/garyburd/redigo/redis"
)

//func Bool(reply interface{}, err error) (bool, error)
//func Bytes(reply interface{}, err error) ([]byte, error)
//func Float64(reply interface{}, err error) (float64, error)
//func Int(reply interface{}, err error) (int, error)
//func Int64(reply interface{}, err error) (int64, error)
//func String(reply interface{}, err error) (string, error)
//func Strings(reply interface{}, err error) ([]string, error)
//func Uint64(reply interface{}, err error) (uint64, error)

func Test_db(t *testing.T) {
	NewDbPool(Redis{
		MaxIdle:            8,
		MaxActive:          64,
		IdleTimeout:        300,
		RedisServer:        "192.168.78.130:6379",
		DialConnectTimeout: 3,
		DialReadTimeout:    3,
		DialWriteTimeout:   3,
		Auth:               "",
	})
	// set单值
	_, err := Exec("SET", "email", "gfandada@gmail.com")
	if err != nil {
		t.Errorf("set error")
	}
	value, err1 := redis.String(Exec("GET", "email"))
	if err1 != nil || value != "gfandada@gmail.com" {
		t.Errorf("get error")
	}
	_, err = Exec("SET", "qq", 1009310068)
	if err != nil {
		t.Errorf("set error")
	}
	value2, err2 := redis.Int(Exec("GET", "qq"))
	if err2 != nil || value2 != 1009310068 {
		t.Errorf("get error")
	}
	_, err = Exec("SET", "github", "https://github.com/gfandada")
	if err != nil {
		t.Errorf("set error")
	}
	// 适用于排行榜的有序set
	// 有序Set，支持每个键值（比如玩家id）拥有一个分数（score），每次往这个set里添加元素，
	// Redis会对其进行排序，修改某一元素的score后，也会更新排序，在获取数据时，可以指定排序范围。
	// 更重要的是，这个排序结果会被保存起来，不用在服务器启动时重新计算。
	Exec("ZADD", "test1", "120", "user1", "122", "user3")
	value3, err3 := redis.Strings(Exec("ZRANGE", "test1", "0", "1111111"))
	if err3 != nil || len(value3) != 2 || value3[0] != "user1" || value3[1] != "user3" {
		t.Errorf("ZRANGE error")
	}
	// 跨服（消息队列）
	// Redis提供的List数据类型，可以用来实现一个消息队列。
	// 由于它是独立于游戏服务器的，所以多个游戏服务器可以通过它来交换数据、发送事件。
	// Redis还提供了发布、订阅的事件模型。
	// 利用这些，我们就不必自己去实现一套服务器间的通信框架，方便地实现服务器组。
	_, err = Exec("LPUSH", "test123", "a", "b", "c")
	if err != nil {
		t.Error(err)
	}
	values, errs := redis.Strings(Exec("BLPOP", "test123", 1000))
	fmt.Println(values, errs)
	if errs != nil {
		t.Error(errs)
	}
	// 自增id
	Exec("SET", "userid", "10000")
	Exec("INCR", "userid")
	Exec("GET", "userid")
}
