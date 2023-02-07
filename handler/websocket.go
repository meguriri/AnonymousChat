package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/meguriri/AnonymousChat/dao"
	"log"
	"net/http"
)

func SendMsg() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取session
		sid, _ := c.Cookie("login")
		s, _ := dao.Get(sid)

		//绑定user
		var user dao.User
		if err := json.Unmarshal(s, &user); err != nil {
			log.Printf("[websocket.SendMsg error] %v json unmarshal err: %v\n", c.ClientIP(), err)
		}

		//定义消息实例
		var message = dao.Message{
			SendUser: user,
			SendTime: c.PostForm("sendtime"),
			Content:  c.PostForm("content"),
		}

		//发送心跳
		client := dao.MyManager.GetClient(sid)
		client.HeartBeat <- struct{}{}

		//将消息放入广播器
		dao.MyManager.BroadCastChan <- message

		//消息存入redis
		//msg,_:=json.Marshal(message)
		//fmt.Println("msg:",string(msg))

		c.JSON(http.StatusOK, gin.H{
			"user": string(s),
			"msg":  "send ok",
		})
	}
}

func ReciveMsg() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[websocket.ReviveMsg] %v start message receive websocket\n", c.ClientIP())

		//获取客户端
		sid, _ := c.Cookie("login")
		client := dao.MyManager.GetClient(sid)

		//http升级为websocket协议
		var upgrader = websocket.Upgrader{}

		//绑定客户端MConn
		var connErr error
		client.MConn, _ = upgrader.Upgrade(c.Writer, c.Request, nil)
		if connErr != nil {
			log.Printf("[websocket.ReciveMsg error] %v upgrade err: %v\n", c.ClientIP(), connErr.Error())
			return
		}
		log.Printf("[websocket.ReciveMsg] %v upgrade is ok\n", c.ClientIP())
	}
}
