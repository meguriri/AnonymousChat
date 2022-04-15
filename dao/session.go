package dao

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/meguriri/AnonymousChat/redis"
	"io"
	"time"
)

var MaxLifetime int64

type Session struct {
	Message     string
	MaxLifetime int64
}

func generateRandomString() string {
	b := make([]byte, 24)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (T *Session) Set() (string, error) {
	sid := generateRandomString()
	fmt.Println("set sid: ", sid)
	err := redis.Rdb.Set(sid, T.Message, time.Duration(T.MaxLifetime)).Err()
	if err != nil {
		return "", err
	}
	return sid, nil
}

func Get(sid string) (string, error) {
	session, err := redis.Rdb.Get(sid).Result()
	if err != nil {
		return "", err
	}
	return session, err
}

func Del(sid string) error {
	err := redis.Rdb.Del(sid).Err()
	if err != nil {
		return err
	}
	return nil
}

func Exist(sid string) bool {
	fmt.Println("checking ", sid)
	ok, _ := redis.Rdb.Exists(sid).Result()
	if ok == 1 {
		return true
	}
	return false
}

//func GetAllSession()(int,error){
//
//}
