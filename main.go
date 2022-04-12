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
	r.Run(":5050")
	// var s dao.Session
	// sid, err1 := s.Set("fuck", 1e10)
	// if err1 != nil {
	// 	panic(err1)
	// }
	// sname, _ := s.Get(sid)
	// fmt.Println(sname)
	// //s.Del(sid)
	// sname, err = s.Get(sid)
	// fmt.Println(sname)

}
