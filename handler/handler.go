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

func GetIndex(c *gin.Context) {
	sid, err := c.Cookie("login")
	if err != nil {
		fmt.Println("no cookie")
		c.HTML(http.StatusOK, "index.html", nil)
	} else {
		ok := dao.Exist(sid)
		if ok {
			c.Redirect(http.StatusFound, "/chatroom")
		} else {
			c.HTML(http.StatusOK, "index.html", nil)
		}
	}
}

func GetOnlineUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"onlineusers": NowUser,
	})
}

func Login(c *gin.Context) {
	var user dao.User
	c.ShouldBindJSON(&user)
	dao.UserList = append(dao.UserList, user)
	s, _ := json.Marshal(user)
	session := dao.Session{Message: string(s), MaxLifetime: dao.MaxLifetime}
	sid, err := session.Set()
	if err != nil {
		panic(err)
	}
	NowUser++
	c.SetCookie("login", sid, int(session.MaxLifetime/1e9), "/", "localhost", false, false)
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func GetChatroom(c *gin.Context) {
	c.HTML(http.StatusOK, "chatroom.html", nil)
}

func GetUserList(c *gin.Context) {
	var upgrader = websocket.Upgrader{}
	var conn, _ = upgrader.Upgrade(c.Writer, c.Request, nil)
	go func(conn *websocket.Conn) {
		//u := []dao.User{
		//	{Nickname: "wawa", Gender: 0, Color: [3]uint{240, 200, 210}},
		//	{Nickname: "mama", Gender: 1, Color: [3]uint{100, 230, 210}},
		//}
		for {
			conn.WriteJSON(dao.UserList)
			time.Sleep(time.Second * 3)
		}

	}(conn)
}

func Offline() func(*gin.Context) {
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
