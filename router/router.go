package router

import (
	"github.com/gin-gonic/gin"
	h "github.com/meguriri/AnonymousChat/handler"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Static("/static", "./static/")
	r.LoadHTMLGlob("templates/*")
	r.GET("/", h.GetIndex)
	r.GET("/onlineusers", h.GetOnlineUsers)
	r.POST("/login", h.Login)
	return r
}
