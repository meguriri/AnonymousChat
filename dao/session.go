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

//session定义
type Session struct {
	Message     string //信息
	MaxLifetime int64  //生存时间
}

//生成随机sid
func generateRandomString() string {
	b := make([]byte, 24)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

//生成session
func (T *Session) Set() (string, error) {
	//生成随机sid
	sid := generateRandomString()
	fmt.Println("set sid: ", sid)
	//写入redis
	err := redis.Rdb.Set(sid, T.Message, time.Duration(T.MaxLifetime)).Err()
	if err != nil {
		return "", err
	}
	//返回生成的sid
	return sid, nil
}

//通过sid获取session
func Get(sid string) ([]byte, error) {
	//从redis中根据sid获取session
	session, err := redis.Rdb.Get(sid).Result()
	if err != nil {
		return nil, err
	}
	//返回session的byte数组
	return []byte(session), err
}

//删除session
func Del(sid string) error {
	err := redis.Rdb.Del(sid).Err()
	if err != nil {
		return err
	}
	return nil
}

//判断session是否存在
func Exist(sid string) bool {
	fmt.Println("checking ", sid)
	//从redis中根据sid判断session是否存在
	ok, _ := redis.Rdb.Exists(sid).Result()
	if ok == 1 {
		return true
	}
	return false
}

//func GetAllSession()(int,error){
//
//}
