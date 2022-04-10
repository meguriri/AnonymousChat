package redis

import "github.com/go-redis/redis"

//声明一个全局的rdb变量
var (
	Rdb      *redis.Client
	Address  string
	Password string
	DB       int
)

//初始化连接
func InitClient() error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     Address,
		Password: Password,
		DB:       DB,
	})
	_, err := Rdb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
