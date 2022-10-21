package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/meguriri/AnonymousChat/dao"
	"net/http"
)

func AuthMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user dao.User
		//获取前端用户信息
		c.ShouldBindJSON(&user)
		if user.Nickname == "1" {
			c.Abort()
			fmt.Println("err!!!!!")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "err"})
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
