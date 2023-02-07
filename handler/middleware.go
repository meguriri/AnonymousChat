package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/meguriri/AnonymousChat/dao"
	"log"
	"net/http"
)

func AuthMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user dao.User
		//获取前端用户信息
		c.ShouldBindJSON(&user)
		if user.Nickname == "" {
			c.Abort()
			log.Printf("[middleware.AuthMiddle error] %v user nickname err\n", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "err"})
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
