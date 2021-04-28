package route

import (
	"gin_admin_api/controller/account"
	"github.com/gin-gonic/gin"
)

func CollectRoute(router *gin.Engine) {
	// 创建账号路由分组
	accountGroup := router.Group("/account")
	account.AccountRouter(accountGroup)
}
