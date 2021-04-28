package main

import (
	"fmt"
	"gin_admin_api/common"
	_ "gin_admin_api/common"
	"gin_admin_api/middleware"
	"gin_admin_api/model"
	"gin_admin_api/route"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {
	common.DB.AutoMigrate(&model.Account{})
}

func main() {
	router := gin.Default()
	// 全局使用中间件
	router.Use(middleware.CorsMiddleware(), middleware.RecoverMiddleware())
	// 注册路由组
	route.CollectRoute(router)

	port := viper.GetString("server.port")
	fmt.Println("当前端口", port)
	if port != "" {
		router.Run(":" + port)
	} else {
		router.Run()
	}
}
