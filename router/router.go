package router

import (
	"github.com/gin-gonic/gin"
	h "github.com/meguriri/AnonymousChat/handler"
)

var (
	HostAddress string
	Port  string
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Static("/static", "./static/")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", h.GetIndex)
	r.GET("/onlineusers", h.GetOnlineUsers)
	r.POST("/login", h.Login)
	chatroom:=r.Group("/chatroom")
	{
		chatroom.GET("/",h.GetChatroom)
	}
	return r
}
