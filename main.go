package main

import (
	"fmt"
	"gin_admin_api/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"

	// 初始化数据库连接及日志文件
	_ "gin_admin_api/common"
	"gin_admin_api/global"
	"gin_admin_api/initialize"
	// 数据模型中init方法的执行
	_ "gin_admin_api/model"
)


func main() {
	// 1.初始化配置
	initialize.InitConfig()
	// 配置中获取端口号
	PORT := strconv.Itoa(global.ServerConfig.Port)
	global.Logger.Info("端口号", zap.String("PORT", PORT))
	global.Logger.Sugar().Infof("打印端口号:%s", PORT)
	router := gin.New()
	router.GET("", func(ctx *gin.Context) {
		utils.Success(ctx, gin.H{
			"data": "成功",
		})
	})
	if err := router.Run(fmt.Sprintf(":%s", PORT)); err != nil {
		fmt.Println("服务启动失败" + err.Error())
		global.Logger.Error("服务启动失败:", zap.String("message", err.Error()))
	}

	// 2.初始化路由
	//Router := initialize.Routers()
	//// 获取端口号
	//PORT := strconv.Itoa(global.ServerConfig.Port)
	//url := ginSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", PORT))
	//Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	//global.Logger.Sugar().Infof("服务已经启动:localhost:%s", PORT)
	//if err := Router.Run(fmt.Sprintf(":%s", PORT)); err != nil {
	//	fmt.Println("服务启动失败" + err.Error())
	//	global.Logger.Error("服务启动失败:", zap.String("message", err.Error()))
	//}
}
