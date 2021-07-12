package main

import (
	"fmt"
	"gin_admin_api/global"
	"gin_admin_api/initialize"
	"gin_admin_api/middleware"
	"gin_admin_api/model"
	"go.uber.org/zap"
	"strconv"
)

func init() {
	global.DB.AutoMigrate(&model.AccountEntity{})
}

func main() {
	// 1.初始化路由
	Router := initialize.Routers()
	// 2.全局使用中间件
	// 全局使用中间件
	Router.Use(middleware.CorsMiddleWare(), middleware.LoggerMiddleWare(), middleware.RecoverMiddleWare())

	// 获取端口号
	PORT := strconv.Itoa(global.ServerConfig.Port)
	fmt.Println(fmt.Sprintf("服务已经启动:localhost:%s", PORT))
	if err := Router.Run(fmt.Sprintf(":%s", PORT)); err != nil {
		fmt.Println("服务启动失败" + err.Error())
		global.Logger.Error("服务启动失败:", zap.String("message", err.Error()))
	}
}
