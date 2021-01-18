package casbinUtil

import (
	"github.com/astaxie/beego"
	"github.com/go-redis/redis/v8"
)

func RedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     beego.AppConfig.DefaultString("cache_redis_host", "172.16.222.253"),
		Password: beego.AppConfig.DefaultString("cache_redis_password", "123456"), // no password set
		DB:       0,  // use default DB
	})
	return rdb
}


