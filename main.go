package main

import (
	"fmt"
	"github.com/meguriri/AnonymousChat/config"
	"github.com/meguriri/AnonymousChat/redis"
	"github.com/meguriri/AnonymousChat/router"
)

func main() {

	err := config.Configinit()
	if err != nil {
		fmt.Println("viper err:", err.Error())
	}
	if err := redis.InitClient(); err != nil {
		fmt.Println(err.Error())
	}
	r := router.InitRouter()
	r.Run(router.HostAddress+":"+router.Port)
}
