package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/meguriri/AnonymousChat/dao"
	"net/http"
)

func SendMsg() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取信息：消息内容，发送人
		sid, _ := c.Cookie("login")
		s, _ := dao.Get(sid)
		var user dao.User
		if err := json.Unmarshal(s, &user); err != nil {
			fmt.Println("json err: ", err)
		}
		var message = dao.Message{
			SendUser: user,
			SendTime: c.PostForm("sendtime"),
			Content:  c.PostForm("content"),
		}
		//msg,_:=json.Marshal(message)
		//fmt.Println("msg:",string(msg))
		dao.MyManager.BroadCastChan <- message
		//redis<-msgs
		c.JSON(http.StatusOK, gin.H{
			"msg": "send ok",
		})
	}
}

func ReciveMsg() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("recive!")
		var upgrader = websocket.Upgrader{}
		var conn, _ = upgrader.Upgrade(c.Writer, c.Request, nil)
		sid, _ := c.Cookie("login")
		fmt.Println("re cookie:", sid)
		client := &dao.Client{
			Id:          sid,
			MConn:       conn,
			MessageChan: make(chan dao.Message, 1024),
		}
		fmt.Println("re client:", client)
		dao.MyManager.Register <- client
		go func(conn *websocket.Conn) {
			for {
				select {
				case data := <-client.MessageChan:
					conn.WriteJSON(data)
				}
			}
		}(conn)
	}
}
