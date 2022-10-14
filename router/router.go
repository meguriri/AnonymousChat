package router

import (
	"github.com/gin-gonic/gin"
	h "github.com/meguriri/AnonymousChat/handler"
)

var (
	HostAddress string
	Port        string
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Static("/static", "./static/")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", h.GetIndex)                    //登录界面
	r.GET("/onlineusers", h.GetOnlineUsers)   //在线人数
	r.POST("/login", h.AuthMiddle(), h.Login) //登录
	chatroom := r.Group("/chatroom")          //聊天室
	{
		chatroom.GET("/", h.GetChatroom)         //聊天室界面
		chatroom.GET("/offline", h.Offline())    //下线
		chatroom.GET("/userlist", h.GetUserList) //在线列表
		chatroom.POST("/send", h.SendMsg())      //发送消息
		chatroom.GET("/recive", h.ReciveMsg())   //接受消息
	}
	return r
}
