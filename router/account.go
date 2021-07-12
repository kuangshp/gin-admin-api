package router

import (
	"gin_admin_api/api"
	"gin_admin_api/middleware"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(Router *gin.RouterGroup) {
	accountRouter := Router.Group("account", middleware.AuthMiddleWare())
	{
		accountRouter.GET("", api.AccountList)
		accountRouter.GET("/:id", api.AccountById)
	}
}
