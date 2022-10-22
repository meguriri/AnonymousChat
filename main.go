package main

import (
	"fmt"
	"github.com/meguriri/AnonymousChat/config"
	"github.com/meguriri/AnonymousChat/dao"
	"github.com/meguriri/AnonymousChat/redis"
	"github.com/meguriri/AnonymousChat/router"
)

func main() {
	//读取配置文件
	err := config.Configinit()
	if err != nil {
		fmt.Println("viper err:", err.Error())
	}

	//初始化redis
	if err := redis.InitClient(); err != nil {
		fmt.Println(err.Error())
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
