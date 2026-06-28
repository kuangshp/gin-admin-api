package router

import (
	"gin-admin-api/internal/api/menu"
	"gin-admin-api/internal/middleware"
	"github.com/casbin/casbin/v2"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func InitMenuRouter(Router *gin.RouterGroup, redis *redis.Client, enforcer *casbin.Enforcer, db *gorm.DB, menuHandler menu.IMenu) {
	authRouter := Router.Group(
		"menu",
		middleware.AuthMiddleWare(redis),
		middleware.OperatorMiddleware(),
		middleware.DataPermissionMiddleware(db),
	)
	authRouter.GET("", menuHandler.GetMenusApi) // 获取菜单列表
}
