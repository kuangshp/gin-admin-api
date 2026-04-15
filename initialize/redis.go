package initialize

import (
	"context"
	"fmt"
	"gin-admin-api/config"
	"github.com/go-redis/redis/v8"
)

// NewRedis Wire Provider：初始化 Redis 客户端，返回 *redis.Client
func NewRedis(cfg *config.ServerConfig) (*redis.Client, error) {
	r := cfg.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.Host, r.Port),
		Password: r.Password,
		DB:       r.DB,
	})

	// 连通性测试
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		// Redis 不强制依赖，连接失败时仅打印警告，不中断启动
		fmt.Printf("⚠ Redis 连接失败（可选依赖）: %v\n", err)
		return client, nil
	}

	fmt.Printf("✔ Redis 连接成功 [%s:%d]\n", r.Host, r.Port)
	return client, nil
}
