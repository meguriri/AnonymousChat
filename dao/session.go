package dao

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/meguriri/AnonymousChat/redis"
	"io"
	"time"
)

type Manager interface {
	Get(sid string) (string, error)
	Set(sname string) (string, error)
	Del(sid string) error
}

type Session struct {
	MaxLifetime int64
}

func generateRandomString() string {
	b := make([]byte, 24)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (T *Session) Set(sname string, maxlife int64) (string, error) {
	T.MaxLifetime = maxlife
	sid := generateRandomString()
	fmt.Println("set sid: ", sid)
	err := redis.Rdb.Set(sid, sname, time.Duration(T.MaxLifetime)).Err()
	if err != nil {
		return "", err
	}
	return sid, nil
}

func (T *Session) Get(sid string) (string, error) {
	sname, err := redis.Rdb.Get(sid).Result()
	if err != nil {
		return "", err
	}
	return sname, err
}

func (T *Session) Del(sid string) error {
	err := redis.Rdb.Del(sid).Err()
	if err != nil {
		return err
	}
	return nil
}

//func  CheckSession(sid string)error{
//	var newSession Session
//	_,err:=newSession.Get(sid)
//	if err !=nil{
//		newSession.Set()
//		return nil
//	}
//}
