package main

import (
	"fmt"
	"gin_admin_api/common"
	"gin_admin_api/global"
	"gin_admin_api/initialize"
	"go.uber.org/zap"
	"strconv"
)

func init() {
	//initialize.InitDataSource()
}

func main() {
	// 1.初始化配置
	initialize.InitConfig()
	// 2.初始化路由
	Router := initialize.Routers()
	// 3.数据库相关的
	common.InitDB()
	initialize.InitDataSource()
	//// 4.全局使用中间件
	//// 全局使用中间件
	//Router.Use(middleware.CorsMiddleWare(), middleware.LoggerMiddleWare(), middleware.RecoverMiddleWare())

	// 获取端口号
	PORT := strconv.Itoa(global.ServerConfig.Port)
	fmt.Println(fmt.Sprintf("服务已经启动:localhost:%s", PORT))
	if err := Router.Run(fmt.Sprintf(":%s", PORT)); err != nil {
		fmt.Println("服务启动失败" + err.Error())
		global.Logger.Error("服务启动失败:", zap.String("message", err.Error()))
	}
}
