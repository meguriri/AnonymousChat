package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meguriri/AnonymousChat/config"
	"github.com/meguriri/AnonymousChat/dao"
	"github.com/meguriri/AnonymousChat/redis"
)

func main() {
	err := config.Configinit()
	if err != nil {
		fmt.Println("viper err:", err.Error())
	}
	if err := redis.InitClient(); err != nil {
		fmt.Println(err.Error())
	}
	r := gin.Default()
	r.Static("/static", "./static/")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
		sid, err := c.Cookie("login")
		if err != nil {
			fmt.Println("no cookie")
		} else {
			//dao.CheckSession(s)
			var s dao.Session
			sname, err := s.Get(sid)
			if err != nil {
				panic(err)
			}
			fmt.Println(sname)
		}
	})

	r.POST("/login", func(c *gin.Context) {
		name := c.PostForm("name")
		fmt.Println(name)
		var s dao.Session
		sid, err := s.Set(name, 1e11)
		if err != nil {
			panic(err)
		}
		c.SetCookie("login", sid, int(s.MaxLifetime/1e9), "/", "localhost", false, false)
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
	})
	r.Run(":5050")

	// var s dao.Session
	// sid, err1 := s.Set("fuck", 1e10)
	// if err1 != nil {
	// 	panic(err1)
	// }
	// sname, _ := s.Get(sid)
	// fmt.Println(sname)
	// //s.Del(sid)
	// sname, err = s.Get(sid)
	// fmt.Println(sname)

}
