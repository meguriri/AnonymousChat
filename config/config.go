package config

import (
	"fmt"
	"github.com/meguriri/AnonymousChat/redis"
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
	fmt.Println(redis.Address, redis.Password, redis.DB)
	return nil
}
