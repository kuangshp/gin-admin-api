package main

import (
	"fmt"
	"gin_admin_api/common"
	_ "gin_admin_api/docs"
	"gin_admin_api/global"
	"gin_admin_api/initialize"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"strconv"
)

// @title 权限系统API文档
// @version 1.0
// @description 使用gin+mysql实现权限系统的api接口
// @host 127.0.0.1:9090/api/v1
// @BasePath
func main() {
	// 1.初始化配置
	initialize.InitConfig()
	// 2.初始化路由
	Router := initialize.Routers()
	// 3.数据库相关的
	common.InitDB()
	initialize.InitDataSource()
	// 获取端口号
	PORT := strconv.Itoa(global.ServerConfig.Port)
	url := ginSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", PORT))
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	fmt.Println(fmt.Sprintf("服务已经启动:localhost:%s", PORT))
	if err := Router.Run(fmt.Sprintf(":%s", PORT)); err != nil {
		fmt.Println("服务启动失败" + err.Error())
		global.Logger.Error("服务启动失败:", zap.String("message", err.Error()))
	}
}
