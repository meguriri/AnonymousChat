package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/meguriri/AnonymousChat/dao"
	"github.com/meguriri/AnonymousChat/logic"
	"github.com/meguriri/AnonymousChat/redis"
	"log"
	"net/http"
)

// 获取index.html
func GetIndex() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取login cookie中的sid
		sid, err := c.Cookie("login")
		//无cookie
		if err != nil {
			log.Printf("[handler.GetIndex error] %v no cookie err: %v\n", c.ClientIP(), err.Error())
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

// 获取在线人数
func GetOnlineUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		//返回在线人数
		c.JSON(http.StatusOK, gin.H{
			"onlineusers": dao.MyManager.GetUserNumber(),
		})
	}
}

// 注册用户
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		nowUsers := dao.MyManager.GetUserNumber()
		if nowUsers == uint(dao.MaxUser) {
			c.JSON(http.StatusOK, gin.H{
				"msg": "当前聊天室人数已满，请稍后重试!",
			})
			return
		}
		//从context中获取user内容，空接口
		ss, ok := c.Get("user")
		if ok {
			user := ss.(dao.User)

			//验证重名
			if logic.CheckDup(user.Nickname) {
				//生成session
				s, _ := json.Marshal(user)
				session := dao.Session{Message: string(s), MaxLifetime: dao.MaxLifetime}

				//生成session id
				sid, err := session.Set()
				if err != nil {
					log.Printf("[handler.Login error] %v session set err: %v\n", c.ClientIP(), err.Error())
				}

				//注册客户端
				dao.MyManager.ClientRegister(sid, &user)

				//生成cookie
				c.SetCookie("login", sid, int(session.MaxLifetime/1e9), "/", "", false, false)
				c.JSON(http.StatusOK, gin.H{
					"msg": "ok",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"msg": "该名称已被使用，请重新输入昵称!",
				})
			}

		}
	}
}

// 获取chatroom.html
func GetChatroom() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "chatroom.html", nil)
	}
}

// 获取用户列表
func GetUserList() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[handler.GetUserList] %v start user list websocket\n", c.ClientIP())
		//获取客户端
		sid, _ := c.Cookie("login")
		client := dao.MyManager.GetClient(sid)

		fmt.Printf("[handler.GetUserList] %v sid: %v\n", client, sid)

		//http升级为websocket协议
		var upgrader = websocket.Upgrader{}
		//绑定客户端LConn
		var connErr error
		client.LConn, connErr = upgrader.Upgrade(c.Writer, c.Request, nil)
		if connErr != nil {
			log.Printf("[handler.GetUserList error] %v upgrade err: %v\n", c.ClientIP(), connErr.Error())
			return
		}
		log.Printf("[handler.GetUserList] %v upgrade is ok\n", c.ClientIP())
		//发送获取初始用户列表
		err := client.LConn.WriteJSON(dao.MyManager.GetUserList())
		if err != nil {
			log.Printf("[handler.GetUserList error] %v origin list err: %v\n", c.ClientIP(), err.Error())
			return
		}
	}
}

// 用户下线
func Offline() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取sid
		sid, _ := c.Cookie("login")

		//注销客户端
		dao.MyManager.ClientUnRegister(sid)

		//redis中删除session
		if err := dao.Del(sid); err != nil {
			log.Printf("[handler.Offline] %v delete session err: %v\n", c.ClientIP(), err.Error())
		}

		//删除cookie
		c.SetCookie("login", "", -1, "/", "localhost", false, false)
		//返回信息
		c.JSON(http.StatusOK, gin.H{
			"msg": "ok",
		})
	}
}

func GetChatsSave() gin.HandlerFunc {
	return func(c *gin.Context) {
		messgaes, err := dao.GetSaveMessage(redis.Rdb)
		if err != nil {
			log.Printf("[handler.GetChatSave error] get save messages err: %v\n", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg":     "error",
				"message": nil,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg":     "ok",
			"message": messgaes,
		})
	}
}
