package router

import (
	"gin_admin_api/api"
	"github.com/gin-gonic/gin"
)

func InitLoginRouter(Router *gin.RouterGroup) {
	loginRouter := Router.Group("login")
	loginRouter.POST("", api.Login)
}
