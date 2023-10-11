package redis

import (
	"fmt"

	"github.com/go-redis/redis"

	"github.com/spf13/viper"
)

var rdb *redis.Client

func Init() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprint("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.poolsize"),
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		fmt.Printf("connect redis failed, err:%v\n", err)
		return
	}

	return
}

func Close() {
	rdb.Close()
}
