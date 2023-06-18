package main

import (
	"fmt"
	"gin-admin-api/global"
	"gin-admin-api/initialize"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	// 初始化数据库连接及日志文件
	_ "gin-admin-api/common"
	// 数据模型中init方法的执行
	_ "gin-admin-api/model"
	// 文档
)

func main() {
	// 1.初始化配置
	initialize.InitConfig()
	// 初始化自定义校验器
	initialize.InitValidate()
	//2.初始化路由
	router := initialize.Routers()
	// 获取端口号
	PORT := strconv.Itoa(global.ServerConfig.Port)
	fmt.Println(PORT + "当前端口")
	global.Logger.Sugar().Infof("服务已经启动:localhost:%s", PORT)
	// 优雅退出程序
	go func() {
		// 启动服务
		if err := router.Run(fmt.Sprintf(":%s", PORT)); err != nil {
			global.Logger.Sugar().Panic("服务启动失败:%s", err.Error())
		}
	}()
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
}
