package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type IRedisUtils interface {
	SetRedisValue(ctx context.Context, key string, value interface{}, expiration int64) error    // 简单的值类型
	SetRedisMapValue(ctx context.Context, key string, value interface{}, expiration int64) error // 复杂的值类型
	GetRedisValue(ctx context.Context, key string) (string, error)                               // GetRedisValue (简单的value) 根据redis的key获取数据
	GetRedisMapValue(ctx context.Context, key string) (map[string]interface{}, error)            // value是一个map的时候
	DelRedisKey(ctx context.Context, key string)                                                 // 根据key来删除
	IncrByKey(ctx context.Context, key string) int                                               // 自增
	DeleteKeyPrefix(ctx context.Context, prefix string)                                          // 根据key前缀来删除
}

type RedisUtils struct {
	redisDb *redis.Client
}

func (r RedisUtils) SetRedisValue(ctx context.Context, key string, value interface{}, expiration int64) error {
	// redis这里要以秒为单位
	err := r.redisDb.Set(ctx, key, value, time.Duration(expiration)*time.Second).Err()
	if err != nil {
		return errors.New("存储redis数据错误" + err.Error())
	}
	return nil
}

func (r RedisUtils) SetRedisMapValue(ctx context.Context, key string, value interface{}, expiration int64) error {
	valueStr, err := json.Marshal(value)
	if err != nil {
		return errors.New("序列化数据失败")
	}
	return r.SetRedisValue(ctx, key, valueStr, expiration)
}

func (r RedisUtils) GetRedisValue(ctx context.Context, key string) (string, error) {
	redisVal, err := r.redisDb.Get(ctx, key).Result()
	if err != nil {
		return "", errors.New(err.Error())
	}
	return redisVal, nil
}

func (r RedisUtils) GetRedisMapValue(ctx context.Context, key string) (map[string]interface{}, error) {
	redisVal, err := r.GetRedisValue(ctx, key)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	var result map[string]interface{}

	if err := json.Unmarshal([]byte(redisVal), &result); err != nil {
		return nil, errors.New("序列化redis数据数据表")
	} else {
		return result, nil
	}
}

func (r RedisUtils) DelRedisKey(ctx context.Context, key string) {
	del := r.redisDb.Del(ctx, key)
	fmt.Println(del, "删除结果")
}

func (r RedisUtils) IncrByKey(ctx context.Context, key string) int {
	incrNum := r.redisDb.Incr(ctx, key)
	u, _ := incrNum.Uint64()
	incrNumStr := strconv.FormatUint(u, 10) // 转换成字符串
	incrNumInt, _ := strconv.Atoi(incrNumStr)
	return incrNumInt
}

func (r RedisUtils) DeleteKeyPrefix(ctx context.Context, prefix string) {
	keys, _ := r.redisDb.Keys(ctx, prefix+"*").Result()
	if len(keys) > 0 {
		err := r.redisDb.Del(ctx, keys...)
		if err != nil {
			panic(err)
		}
		fmt.Printf("成功删除 %d 个键\n", len(keys))
	} else {
		fmt.Println("没有匹配的键需要删除")
	}
}
func NewRedisUtils(redisDb *redis.Client) IRedisUtils {
	return RedisUtils{
		redisDb: redisDb,
	}
}
