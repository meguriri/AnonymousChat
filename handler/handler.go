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

var (
	NowUser int
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
		"onlineusers": NowUser,
	})
}

//注册用户
func Login(c *gin.Context) {
	//定义user实例
	var user dao.User
	//从context中获取user内容，空接口
	ss, ok := c.Get("user")
	if ok {
		fmt.Println("login: ", ss)
		//空接口转换user
		fmt.Println("Login: ", user)
		dao.UserList = append(dao.UserList, user)
		//
		s, _ := json.Marshal(user)
		session := dao.Session{Message: string(s), MaxLifetime: dao.MaxLifetime}
		sid, err := session.Set()
		if err != nil {
			panic(err)
		}
		//
		NowUser++
		//
		c.SetCookie("login", sid, int(session.MaxLifetime/1e9), "/", "localhost", false, false)
		c.JSON(http.StatusOK, gin.H{
			"msg": "ok",
		})
	}
}

func GetChatroom(c *gin.Context) {
	c.HTML(http.StatusOK, "chatroom.html", nil)
}

func GetUserList(c *gin.Context) {
	var upgrader = websocket.Upgrader{}
	var conn, _ = upgrader.Upgrade(c.Writer, c.Request, nil)
	go func(conn *websocket.Conn) {
		for {
			conn.WriteJSON(dao.UserList)
			time.Sleep(time.Second * 15)
		}
	}(conn)
}

func Offline() gin.HandlerFunc {
	return func(c *gin.Context) {
		NowUser--
		sid, _ := c.Cookie("login")
		fmt.Println("cookie: ", sid)
		if err := dao.Del(sid); err != nil {
			fmt.Println(err)
		}
		c.SetCookie("login", "", -1, "/", "localhost", false, false)
		c.JSON(http.StatusOK, gin.H{
			"msg": "ok",
		})
	}
}
