package network

import (
	"testing"
	"time"

	"github.com/gfandada/gserver/db"
	"github.com/gfandada/gserver/util"
)

// test for session_map.go

//func Test_OptSession(t *testing.T) {
//	NewSessionMap()
//	ret := OptSession(uint64(1), Find, []SessionData{SessionData{
//		Key: "key",
//	}})
//	if ret != nil {
//		t.Error("user 1 should be nil")
//	}
//	// add user 1
//	OptSession(uint64(1), Update, nil)
//	// add user 1
//	OptSession(uint64(1), Update, []SessionData{SessionData{
//		Key:   "key1",
//		Value: 9257,
//	}})
//	// add user 1
//	OptSession(uint64(1), Update, []SessionData{SessionData{
//		Key:   "key2",
//		Value: "hello",
//	}, SessionData{
//		Key:   "key3",
//		Value: 128.123123,
//	}})
//	// find
//	ret1 := OptSession(uint64(1), Find, nil)
//	if ret1 != nil {
//		t.Error("find error")
//	}
//	ret2 := OptSession(uint64(1), Find, []SessionData{SessionData{
//		Key: "key2",
//	}})
//	if ret2 == nil || len(ret2) != 1 || ret2[0].Value != "hello" {
//		t.Error("find error")
//	}
//	ret3 := OptSession(uint64(1), Find, []SessionData{SessionData{
//		Key: "key2",
//	}, SessionData{
//		Key: "key3",
//	}, SessionData{
//		Key: "key4",
//	}})
//	if ret3 == nil || len(ret3) != 3 {
//		t.Error("find error")
//	}
//	if ret3[0].Value != "hello" || ret3[1].Value != 128.123123 || ret3[2].Value != nil {
//		t.Error("find error")
//	}
//	// update key3 key4
//	OptSession(uint64(1), Update, []SessionData{SessionData{
//		Key:   "key3",
//		Value: 12222, // 128.123123
//	}, SessionData{
//		Key:   "key4",
//		Value: "xixixixiixiix",
//	}})
//	ret4 := OptSession(uint64(1), Find, []SessionData{SessionData{
//		Key: "key3",
//	}, SessionData{
//		Key: "key4",
//	}})
//	if ret4 == nil || len(ret4) != 2 {
//		t.Error("find error")
//	}
//	if ret4[0].Value != 12222 || ret4[1].Value != "xixixixiixiix" {
//		t.Error("find error")
//	}
//	// delete ke3
//	OptSession(uint64(1), Delete, []SessionData{SessionData{
//		Key: "key3",
//	}})
//	ret5 := OptSession(uint64(1), Find, []SessionData{SessionData{
//		Key: "key3",
//	}, SessionData{
//		Key: "key4",
//	}})
//	if ret5 == nil || len(ret5) != 2 {
//		t.Error("find error")
//	}
//	if ret5[0].Value != nil || ret5[1].Value != "xixixixiixiix" {
//		t.Error("find error")
//	}
//}

// test for session_redis.go

func Test_session_redis(t *testing.T) {
	db.NewDbPool(db.Redis{
		MaxIdle:            8,
		MaxActive:          64,
		IdleTimeout:        300,
		RedisServer:        "192.168.78.130:6379",
		DialConnectTimeout: 3,
		DialReadTimeout:    3,
		DialWriteTimeout:   3,
		Auth:               "",
	})
	// get
	_, err := GetSession("123123123")
	if err == nil {
		t.Error("GetSession error")
	}
	id := util.NewV4().String()
	// set
	SetSession(&Session{
		SessionId:  id,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	})
	// get
	ret, _ := GetSession(id)
	// update
	SetSession(&Session{
		SessionId:  id,
		UpdateTime: time.Now(),
	})
	ret1, _ := GetSession(id)
	if !ret1.CreateTime.Equal(ret.CreateTime) || ret1.UpdateTime.Before(ret.UpdateTime) {
		t.Error("update error")
	}
}
