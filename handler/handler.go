package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/meguriri/AnonymousChat/dao"
	"net/http"
	"time"
)

//获取index.html
func GetIndex(c *gin.Context) {
	//获取login cookie中的sid
	sid, err := c.Cookie("login")
	//无cookie
	if err != nil {
		fmt.Println("no cookie")
		//返回index.html
		c.HTML(http.StatusOK, "index.html", nil)
	} else { //有cookie
		//查询sid
		ok := dao.Exist(sid)
		if ok {
			//sid存在，转到聊天室
			c.Redirect(http.StatusFound, "/chatroom")
		} else {
			//返回index.html
			c.HTML(http.StatusOK, "index.html", nil)
		}
	}
}

//获取在线人数
func GetOnlineUsers(c *gin.Context) {
	//返回在线人数
	c.JSON(http.StatusOK, gin.H{
		"onlineusers": Manager.GetUserNumber(),
	})
}

//注册用户
func Login(c *gin.Context) {

	//从context中获取user内容，空接口
	ss, ok := c.Get("user")
	if ok {
		user := ss.(dao.User)
		fmt.Println("Login: ", user)

		//dao.UserList = append(dao.UserList, user)

		//生成session
		s, _ := json.Marshal(user)
		session := dao.Session{Message: string(s), MaxLifetime: dao.MaxLifetime}

		//生成session id
		sid, err := session.Set()
		if err != nil {
			panic(err)
		}

		//注册客户端
		Manager.ClientRegist(sid, &user)

		//dao.NowUser++

		//生成cookie
		c.SetCookie("login", sid, int(session.MaxLifetime/1e9), "/", "localhost", false, false)
		c.JSON(http.StatusOK, gin.H{
			"msg": "ok",
		})
	}
}

//获取chatroom.html
func GetChatroom(c *gin.Context) {
	c.HTML(http.StatusOK, "chatroom.html", nil)
}

//获取用户列表
func GetUserList(c *gin.Context) {
	//http升级为websocket协议
	var upgrader = websocket.Upgrader{}

	//获取一个conn实例
	var conn, _ = upgrader.Upgrade(c.Writer, c.Request, nil)

	go func(conn *websocket.Conn) {
		for {
			conn.WriteJSON(Manager.GetUserList())
			time.Sleep(time.Second * 15)
		}
	}(conn)
}

func Offline() gin.HandlerFunc {
	return func(c *gin.Context) {

		//在线用户数量-1
		//dao.NowUser--

		//获取sid
		sid, _ := c.Cookie("login")
		fmt.Println("cookie: ", sid)

		//注销客户端
		Manager.ClientUnRegist(sid)

		//redis中删除session
		if err := dao.Del(sid); err != nil {
			fmt.Println(err)
		}

		//删除cookie
		c.SetCookie("login", "", -1, "/", "localhost", false, false)
		//返回信息
		c.JSON(http.StatusOK, gin.H{
			"msg": "ok",
		})
	}
}
