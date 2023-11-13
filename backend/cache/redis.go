package cache

import (
	conf "cms/config"
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() error {
	redisConf := conf.InitConf().Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         redisConf.Local + ":" + strconv.Itoa(redisConf.Port),
		Password:     redisConf.Password,     // 无密码则为空
		DB:           redisConf.Db,           // 默认数据库
		PoolSize:     redisConf.Poolsize,     // 连接池大小
		MinIdleConns: redisConf.Minidleconns, // 维持空闲连接数
		MaxRetries:   redisConf.Maxretries,   // 最大重试次数
	})

	// 测试连接
	_, err := RedisClient.Ping(Ctx).Result()
	return err
}
