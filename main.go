package main

import (
	"fmt"
	"gin_admin_api/global"
	"gin_admin_api/initialize"
	"gin_admin_api/model"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"strconv"

	// 初始化数据库连接及日志文件
	_ "gin_admin_api/common"
	// 数据模型中init方法的执行
	_ "gin_admin_api/model"
	// 文档
	_ "gin_admin_api/docs"
)


// @title 权限系统API文档
// @version 1.0
// @description 使用gin+mysql实现权限系统的api接口
// @host 127.0.0.1:9090/api/v1
// @BasePath
func main() {
	// 1.初始化配置
	initialize.InitConfig()
	// 初始化自定义校验器
	initialize.InitValidate()

	// 初始化数据库(如果是手动创建数据表的时候不需要这个)
	global.DB.AutoMigrate(&model.AccountEntity{})

	//2.初始化路由
	router := initialize.Routers()
	// 获取端口号
	PORT := strconv.Itoa(global.ServerConfig.Port)
	fmt.Println(PORT + "当前端口")

	url := ginSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", PORT))
	// swagger访问地址:localhost:9000/swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	global.Logger.Sugar().Infof("服务已经启动:localhost:%s", PORT)

	// 启动服务
	if err := router.Run(fmt.Sprintf(":%s", PORT)); err != nil {
		global.Logger.Sugar().Panic("服务启动失败:%s", err.Error())
	}
}
