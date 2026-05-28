package main

import (
	"flag"
	"fmt"
	"gin-admin-api/initialize"
	"log"
)

var envString string

func init() {
	flag.StringVar(&envString, "envString", "dev", "环境配置：dev | prod")
}

// @title gin-admin-api
// @version 1.0
// @description gin-admin-api 后台管理接口文档
// @BasePath /
func main() {
	flag.Parse()
	fmt.Printf("▶ 启动环境: %s\n", envString)

	if err := initialize.GenerateSwaggerDocs(); err != nil {
		log.Fatalf("生成 Swagger 文档失败: %v", err)
	}

	app, err := InitApp(envString)
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	// 初始化管理员账号数据
	if err = app.InitSqlData(); err != nil {
		log.Fatalf("初始化管理员账号失败: %v", err)
	}

	if err = app.Run(); err != nil {
		log.Fatalf("服务异常退出: %v", err)
	}
}
