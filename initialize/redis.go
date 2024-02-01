package initialize

import (
	"gin-admin-api/global"
	"github.com/go-redis/redis/v8"
)

func InitRedisDb() {
	redisDb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	global.RedisDb = redisDb
}
