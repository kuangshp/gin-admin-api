package utils

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type IRedisUtils interface {
	ExistsKey(ctx context.Context, key string) (int64, error)                                   // 检查key是否存在
	SetNxRedisValue(ctx context.Context, key string, value interface{}, expiration int64) error // 不存的时候才设置
	SetRedisValue(ctx context.Context, key string, value interface{}, expiration int64) error   // 简单的值类型
	GetRedisValue(ctx context.Context, key string) (string, error)                              // GetRedisValue (简单的value) 根据redis的key获取数据
	SetRedisMapValue(ctx context.Context, value ...interface{}) error                           // 复杂的值类型
	GetRedisMapValue(ctx context.Context, keys ...string) ([]interface{}, error)                // value是一个map的时候
	DelRedisMapKey(ctx context.Context, key1, key2 string) (int64, error)                       // 删除map其中一个key
	DelRedisKey(ctx context.Context, key string)                                                // 根据key来删除
	IncrByKey(ctx context.Context, key string) int                                              // 自增
	DeleteKeyPrefix(ctx context.Context, prefix string)                                         // 根据key前缀来删除
}

type RedisUtils struct {
	redisDb *redis.Client
}

func (r RedisUtils) ExistsKey(ctx context.Context, key string) (int64, error) {
	result, _ := r.redisDb.Exists(ctx, key).Result()
	if result > 0 {
		return result, nil
	}
	return 0, errors.New("key不存在")
}

func (r RedisUtils) SetNxRedisValue(ctx context.Context, key string, value interface{}, expiration int64) error {
	// redis这里要以秒为单位,0表示永不过期
	err := r.redisDb.SetNX(ctx, key, value, time.Duration(expiration)*time.Second).Err()
	if err != nil {
		return errors.New("存储redis数据错误" + err.Error())
	}
	return nil
}

func (r RedisUtils) SetRedisValue(ctx context.Context, key string, value interface{}, expiration int64) error {
	// redis这里要以秒为单位,0表示永不过期
	err := r.redisDb.Set(ctx, key, value, time.Duration(expiration)*time.Second).Err()
	if err != nil {
		return errors.New("存储redis数据错误" + err.Error())
	}
	return nil
}

func (r RedisUtils) SetRedisMapValue(ctx context.Context, value ...interface{}) error {
	err := r.redisDb.MSet(ctx, value).Err()
	if err != nil {
		return errors.New("设置值失败")
	}
	return nil
}

func (r RedisUtils) GetRedisValue(ctx context.Context, key string) (string, error) {
	redisVal, err := r.redisDb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", errors.New("传递的key错误")
	} else if err != nil {
		return "", errors.New(err.Error())
	}
	return redisVal, nil
}

func (r RedisUtils) GetRedisMapValue(ctx context.Context, keys ...string) ([]interface{}, error) {
	return r.redisDb.MGet(ctx, keys...).Result()
}

func (r RedisUtils) DelRedisMapKey(ctx context.Context, key1, key2 string) (int64, error) {
	return r.redisDb.Del(ctx, key1, key2).Result()
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
