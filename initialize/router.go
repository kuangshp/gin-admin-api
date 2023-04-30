package initialize

import (
	"gin-admin-api/middleware"
	"gin-admin-api/router"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	// 注册中间件
	Router.Use(
		middleware.CorsMiddleWare(),    // 跨域的
		middleware.LoggerMiddleWare(),  // 日志
		middleware.RecoverMiddleWare(), // 异常的
	)
	// 配置全局路径
	ApiGroup := Router.Group("/api/v1/admin")
	// 注册路由
	router.InitAccountRouter(ApiGroup) // 账号中心
	return Router
}
