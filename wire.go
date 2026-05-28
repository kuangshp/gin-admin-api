//go:build wireinject
// +build wireinject

package main

import (
	"gin-admin-api/initialize"
	"gin-admin-api/internal/api/account"
	accountMapper "gin-admin-api/internal/api/account/mapper"
	"gin-admin-api/internal/api/auth"
	authMapper "gin-admin-api/internal/api/auth/mapper"
	"gin-admin-api/internal/api/base"
	"gin-admin-api/internal/dal/repository"
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
		// 数据访问层
		repository.NewSysAccountRepository,
		repository.NewSysAccountRoleRepository,
		// mapper转换层
		accountMapper.NewSysAccountMapper,
		authMapper.NewAuthMapper,
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
