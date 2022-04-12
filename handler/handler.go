package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/meguriri/AnonymousChat/dao"
	"net/http"
)

func GetIndex(c *gin.Context) {
	sid, err := c.Cookie("login")
	if err != nil {
		fmt.Println("no cookie")
		c.HTML(http.StatusOK, "index.html", nil)
	} else {
		//dao.CheckSession(s)
		var s dao.Session
		sname, err := s.Get(sid)
		if err != nil {
			panic(err)
		}
		fmt.Println(sname)
		//
	}
}

func GetOnlineUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"onlineusers": 1000,
	})
}

func Login(c *gin.Context) {
	var user dao.User
	c.ShouldBindJSON(&user)
	fmt.Println(user)
	//var s dao.Session
	//sid, err := s.Set(user.Nickname, 1e11)
	//if err != nil {
	//	panic(err)
	//}
	//c.SetCookie("login", sid, int(s.MaxLifetime/1e9), "/", "localhost", false, false)

	//a :=map[string]interface{}{}
	//a["name"]=name
	//a["sex"]="male"
	//err:=redis.Rdb.HMSet(name,a).Err()
	//if err!=nil{
	//	fmt.Println("redis err: ",err.Error() )
	//}
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}
