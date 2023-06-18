package global

import (
	"gin-admin-api/config"
	"gin-admin-api/dao"
	"go.uber.org/zap"
)

var (
	// ServerConfig 配置项
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	// Logger 日志
	Logger *zap.Logger
	// DB 数据库
	DB dao.Query
)
