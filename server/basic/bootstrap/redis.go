package bootstrap

import (
	"fmt"
	"seckil/server/basic/config"

	"github.com/go-redis/redis/v8"
)

func InitRedis() {
	configRedis := config.GlobalConfig.Redis
	Addr := fmt.Sprintf("%s:%d", configRedis.Host, configRedis.Port)
	config.Rdb = redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: configRedis.Password, // no password set
		DB:       configRedis.Database, // use default DB
	})

	err := config.Rdb.Ping(config.Ctx).Err()
	if err != nil {
		panic("redis连接失败" + err.Error())
	}
	fmt.Println("redis连接成功")

}
