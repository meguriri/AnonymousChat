package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/meguriri/AnonymousChat/dao"
	"net/http"
)

//获取index.html
func GetIndex() gin.HandlerFunc {
	return func(c *gin.Context) {
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
}

//获取在线人数
func GetOnlineUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		//返回在线人数
		c.JSON(http.StatusOK, gin.H{
			"onlineusers": dao.MyManager.GetUserNumber(),
		})
	}
}

//注册用户
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从context中获取user内容，空接口
		ss, ok := c.Get("user")
		if ok {
			user := ss.(dao.User)
			fmt.Println("Login: ", user)

			//生成session
			s, _ := json.Marshal(user)
			session := dao.Session{Message: string(s), MaxLifetime: dao.MaxLifetime}

			//生成session id
			sid, err := session.Set()
			if err != nil {
				fmt.Println("session set error", err)
			}

			//注册客户端
			dao.MyManager.ClientRegist(sid, &user)

			//生成cookie
			c.SetCookie("login", sid, int(session.MaxLifetime/1e9), "/", "localhost", false, false)
			c.JSON(http.StatusOK, gin.H{
				"msg": "ok",
			})
		}
	}
}

//获取chatroom.html
func GetChatroom() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "chatroom.html", nil)
	}
}

//获取用户列表
func GetUserList() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取客户端
		sid, _ := c.Cookie("login")
		client := dao.MyManager.GetClient(sid)

		//http升级为websocket协议
		var upgrader = websocket.Upgrader{}

		//获取一个conn实例
		client.LConn, _ = upgrader.Upgrade(c.Writer, c.Request, nil)

		//发送获取初始用户列表
		err := client.LConn.WriteJSON(dao.MyManager.GetUserList())
		if err != nil {
			fmt.Println("origin list error", err)
		}
		//client.UserListSignal <- struct{}{}

	}
}

func Offline() gin.HandlerFunc {
	return func(c *gin.Context) {

		//获取sid
		sid, _ := c.Cookie("login")
		fmt.Println("cookie: ", sid)

		//注销客户端
		dao.MyManager.ClientUnRegist(sid)

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
