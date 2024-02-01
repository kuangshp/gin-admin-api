package global

import (
	"gin-admin-api/config"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	// ServerConfig 配置项
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	// Logger 日志
	Logger *zap.Logger
	// DB 数据库
	DB      *gorm.DB
	RedisDb *redis.Client
)
