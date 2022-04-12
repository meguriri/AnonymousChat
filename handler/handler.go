package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/meguriri/AnonymousChat/dao"
	"net/http"
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
		ok:= dao.Exist(sid)
		if ok{
			c.Redirect(http.StatusFound, "/chatroom")
		}else{
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
	s,_ := json.Marshal(user)
	session:= dao.Session{Message: string(s), MaxLifetime: dao.MaxLifetime}
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

func GetChatroom(c *gin.Context){
	c.HTML(http.StatusOK,"chatroom.html",nil)
}