package config

import (
	"github.com/meguriri/AnonymousChat/dao"
	"github.com/meguriri/AnonymousChat/redis"
	"github.com/meguriri/AnonymousChat/router"
	"github.com/spf13/viper"
)

func Configinit() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	redis.Address = viper.GetString("redis.address")
	redis.Password = viper.GetString("redis.password")
	redis.DB = viper.GetInt("DB")
	router.HostAddress = viper.GetString("server.router.hostaddress")
	router.Port = viper.GetString("server.router.port")
	dao.MaxUser = viper.GetInt("server.maxuser")
	dao.MaxLifetime = viper.GetInt64("session.maxlifetime")
	return nil
}
