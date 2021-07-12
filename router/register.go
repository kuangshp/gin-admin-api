package router

import (
	"gin_admin_api/api"
	"github.com/gin-gonic/gin"
)

func InitRegisterRouter(Router *gin.RouterGroup) {
	registerRouter := Router.Group("register")
	registerRouter.POST("", api.Register)
}
