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
	s, _ := json.Marshal(user)
	session := dao.Session{Message: string(s), MaxLifetime: dao.MaxLifetime}
	sid, err := session.Set()

	//添加用户至dao.Users
	dao.Users = make(map[string]*(dao.User))
	dao.Users[sid] = &user

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
		for _, user := range dao.Users {
			dao.UserList = append(dao.UserList, *user)
		}
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
		delete(dao.Users, sid)
		dao.UserList = make([]dao.User, 0, len(dao.Users))
		for _, user := range dao.Users {
			dao.UserList = append(dao.UserList, *user)
		}
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
