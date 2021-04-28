package route

import (
	"gin_admin_api/controller/account"
	"gin_admin_api/controller/login"
	"gin_admin_api/controller/register"
	"gin_admin_api/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRoute(router *gin.Engine) {
	// 创建账号路由分组
	accountGroup := router.Group("/account", middleware.AuthMiddleWare())
	account.AccountRouter(accountGroup)
	// 登录的路由
	loginGroup := router.Group("/login")
	login.LoginRouter(loginGroup)

	registerGroup := router.Group("/register")
	register.RegisterRouter(registerGroup)
}
