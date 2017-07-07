package network

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gfandada/gserver/db"
)

type Session struct {
	SessionId  string    // sessionId
	CreateTime time.Time // session创建时间
	UpdateTime time.Time // session最近一次更新时间
	ExpireTime time.Time // FIXME session过期时间：暂未使用
	EndTime    time.Time // 离线时间
	Data       Isession  // 支持存放部分业务数据
}

func SetSession(session *Session) error {
	value, err := GetSession(session.SessionId)
	var newSession *Session
	if value == nil {
		newSession = session
	} else {
		newSession = merge(session, value)
	}
	ret, err1 := json.Marshal(newSession)
	if err1 != nil {
		return err1
	}
	_, err = db.Exec("SET", session.SessionId, ret)
	return err
}

func merge(newSession *Session, oldSession *Session) *Session {
	session := new(Session)
	if !newSession.CreateTime.Equal(time.Time{}) {
		session.CreateTime = newSession.CreateTime
	} else {
		session.CreateTime = oldSession.CreateTime
	}
	if !newSession.UpdateTime.Equal(time.Time{}) {
		session.UpdateTime = newSession.UpdateTime
	} else {
		session.UpdateTime = oldSession.UpdateTime
	}
	if !newSession.ExpireTime.Equal(time.Time{}) {
		session.ExpireTime = newSession.ExpireTime
	} else {
		session.ExpireTime = oldSession.ExpireTime
	}
	session.SessionId = oldSession.SessionId
	session.Data = newSession.Data
	return session
}

func DelSession(sessionID string) error {
	_, err := db.Exec("DEL", sessionID)
	return err
}

func GetSession(sessionID string) (*Session, error) {
	value, err := redis.Bytes(db.Exec("GET", sessionID))
	if err != nil {
		return nil, errors.New("session not exists, sessionID: " + sessionID)
	}
	session := &Session{}
	err = json.Unmarshal(value, session)
	if err != nil {
		return nil, err
	}
	return session, nil
}
