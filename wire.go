//go:build wireinject
// +build wireinject

package main

import (
	"gin-admin-api/initialize"
	"gin-admin-api/internal/api/account"
	"gin-admin-api/internal/api/auth"
	"gin-admin-api/internal/api/base"
	"gin-admin-api/internal/router"

	"github.com/google/wire"
)

// InitApp Wire 入口函数，描述完整依赖图
// 执行 make wire 后自动生成 wire_gen.go，本文件不参与编译
func InitApp(envString string) (*initialize.App, error) {
	wire.Build(
		// 基础设施
		initialize.NewConfig,
		initialize.NewLogger,
		initialize.NewDB,
		initialize.NewRedis,
		// 基础控制器
		base.NewBaseApi,
		// 接入层
		auth.NewAuth,
		account.NewSysAccount,
		// 路由 & 服务
		router.NewAdminRouter,
		initialize.NewRouter,
		initialize.NewApp,
	)
	return nil, nil
}
