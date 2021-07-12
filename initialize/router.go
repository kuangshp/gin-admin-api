package initialize

import (
	"gin_admin_api/middleware"
	"gin_admin_api/router"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middleware.CorsMiddleWare(), middleware.LoggerMiddleWare(), middleware.RecoverMiddleWare())
	ApiGroup := Router.Group("/api/v1")
	router.InitUserRouter(ApiGroup)
	router.InitLoginRouter(ApiGroup)
	router.InitRegisterRouter(ApiGroup)
	return Router
}
