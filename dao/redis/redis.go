package redis

import (
	"fmt"
	"go_web_app/settings"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

func Init(conf *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
		PoolSize: conf.PoolSize,
		//Password: viper.GetString("redis.password"),
		//DB:       viper.GetInt("redis.db"),
		//PoolSize: viper.GetInt("redis.poolsize"),
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
