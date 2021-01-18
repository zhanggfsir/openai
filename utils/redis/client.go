package redis

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	//"push_serv/config"
	//"push_serv/log"
	"sync"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
)

//var Rdb = redis.NewClient(&redis.Options{
//	Addr:     beego.AppConfig.DefaultString("cache_redis_host", ""),
//	Password: beego.AppConfig.DefaultString("cache_redis_password", ""), // no password set
//	DB:       4,                                                         // use default DB
//})

// C 对外使用的对象
var (
	client *redis.Client

	locker    *redislock.Client
	lockermtx sync.Mutex

	Nil = redis.Nil
)

func init() {
	logs.Info("-------init redis --------------------")
	client = redis.NewClient(&redis.Options{
		Addr:     beego.AppConfig.DefaultString("cache_redis_host", ""),
		Password: beego.AppConfig.DefaultString("cache_redis_password", ""), // no password set
		DB:       4,

		//Addr:     config.Config.Redis.Addr,
		//Password: config.Config.Redis.Password,
		//DB:       config.Config.Redis.LockDB,
	})

	pingcmd := client.Ping(context.Background())
	if err := pingcmd.Err(); nil != err {
		//log.Error.Err(err).Msg("redis conn error")
		logs.Error("redis conn error")
		panic(err)
	}

	locker = redislock.New(client)
}

// C 获取client
func C() *redis.Client {
	return client
}

