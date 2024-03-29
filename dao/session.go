package dao

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"time"

	"github.com/meguriri/AnonymousChat/redis"
)

var MaxLifetime int64

type Session struct {
	Message     string //信息
	MaxLifetime int64  //生存时间
}

func generateRandomString() string { //生成随机sid
	b := make([]byte, 24)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (T *Session) Set() (string, error) { //生成session
	//生成随机sid
	sid := generateRandomString()
	log.Printf("[session.Set] set sid: %v\n", sid)
	//写入redis
	err := redis.Rdb.Set(sid, T.Message, 0).Err()
	if err != nil {
		return "", err
	}
	if ok, err := redis.Rdb.Expire(sid, time.Second*(time.Duration(T.MaxLifetime)/1e9)).Result(); !ok {
		log.Printf("[session.Set] set expire time err: %v\n", err)
		return "", err
	}
	//返回生成的sid
	return sid, nil
}

func ExpireTime(sid string, newTime int64) (bool, error) {
	var err error
	var ok bool = false
	if ok, err = redis.Rdb.Expire(sid, time.Second*time.Duration(newTime)).Result(); ok {
		return ok, nil
	}
	return ok, err
}

func Get(sid string) ([]byte, error) { //通过sid获取session
	//从redis中根据sid获取session
	session, err := redis.Rdb.Get(sid).Result()
	if err != nil {
		return nil, err
	}
	//返回session的byte数组
	return []byte(session), err
}

func Del(sid string) error { //删除session
	err := redis.Rdb.Del(sid).Err()
	if err != nil {
		return err
	}
	return nil
}

func Exist(sid string) bool { //判断session是否存在
	log.Printf("[session.Exist] checking %v\n", sid)
	//从redis中根据sid判断session是否存在
	ok, _ := redis.Rdb.Exists(sid).Result()
	return ok == 1
}
