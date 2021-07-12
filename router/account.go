package router

import (
	"gin_admin_api/api"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(Router *gin.RouterGroup) {
	//accountRouter := Router.Group("account", middleware.AuthMiddleWare())
	accountRouter := Router.Group("account")
	{
		accountRouter.GET("", api.AccountList)
		accountRouter.GET("/:id", api.AccountById)
	}
}
