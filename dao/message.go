package dao

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
)

//type Message interface {
//	Save(*redis.Client)
//}

//type PictureMessage struct {
//}
//
//type FileMessage struct {
//}
//
//type VoiceMessage struct {
//}

var (
	MaxMessageSave int
)

type Message struct {
	//发送用户
	SendUser User `json:"senduser"`
	//发送时间
	SendTime string `json:"sendtime"`
	//内容
	Content string `json:"content"`
}

func (t *Message) Save(r *redis.Client) error {
	log.Printf("[message.Save] start save\n")
	//获取聊天记录长度
	messageListLen, err := r.LLen("message").Result()
	if err != nil {
		return err
	}

	//获取聊天的记录的json
	message, _ := json.Marshal(*t)
	log.Printf("[message.Save] get message: %v\n", message)
	//大于保存上限
	if messageListLen >= int64(MaxMessageSave) {
		//删除一条记录
		message, err := r.LPop("message").Result()
		if err != nil {
			return err
		}
		log.Printf("[message.Save] pop result: %v\n", message)
	}

	//插入聊天记录
	r.RPush("message", message)

	return nil
}

func GetSaveMessage(r *redis.Client) ([]Message, error) {
	messages := make([]Message, MaxMessageSave)
	strs, err := r.LRange("message", 0, -1).Result()
	if err != nil {
		return nil, err
	}
	for i, v := range strs {
		err := json.Unmarshal([]byte(v), &messages[i])
		if err != nil {
			return nil, err
		}
	}
	log.Printf("[message.GetSaveMessage] Get All Messages ok: %+v\n", messages)
	return messages, nil
}
