package main

import (
	"github.com/meguriri/AnonymousChat/config"
	"github.com/meguriri/AnonymousChat/dao"
	"github.com/meguriri/AnonymousChat/redis"
	"github.com/meguriri/AnonymousChat/router"
	"log"
)

func main() {
	//读取配置文件
	if err := config.Configinit(); err != nil {
		log.Printf("[config.Configinit error] viper err: %v\n", err.Error())
		return
	}

	//初始化redis
	if err := redis.InitClient(); err != nil {
		log.Printf("[redis.InitClient error] redis connect err: %v\n", err.Error())
		return
	}

	//初始化管理器
	dao.MyManager.InitManager()

	//开启一个协程管理websocket连接
	go dao.MyManager.Managed()

	//初始化路由器
	r := router.InitRouter()

	//运行
	r.Run(router.HostAddress + ":" + router.Port)
}
