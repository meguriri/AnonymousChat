package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/meguriri/AnonymousChat/dao"
	"net/http"
)

var (
	Manager dao.Manager
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
		Manager.BroadCastChan <- message
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
			Socket:      conn,
			MessageChan: make(chan dao.Message, 1024),
		}
		fmt.Println("re client:", client)
		Manager.Register <- client
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

func WebsocketManange() {
	fmt.Println("start manager!!")
	for {
		fmt.Println("choose")
		select {
		case client := <-Manager.Register:
			fmt.Println("<-ma")
			//注册客户端
			Manager.Lock.Lock()
			Manager.Group[client.Id] = client
			Manager.ClientCount += 1
			fmt.Printf("客户端注册: 客户端id为%s", client.Id)
			Manager.Lock.Unlock()
		case client := <-Manager.UnRegister:
			//注销客户端
			Manager.Lock.Lock()
			if _, ok := Manager.Group[client.Id]; ok {
				//关闭消息通道
				close(client.MessageChan)
				//删除分组中客户
				delete(Manager.Group, client.Id)
				//客户端数量减1
				Manager.ClientCount -= 1
				fmt.Printf("客户端注销: 客户端id为%s", client.Id)
			}
			Manager.Lock.Unlock()
		case data := <-Manager.BroadCastChan:
			//将数据广播给所有客户端
			for _, conn := range Manager.Group {
				fmt.Println(conn.Id, data)
				conn.MessageChan <- data
			}

		}
	}
}
