package router

import (
	"gin-admin-api/internal/api/role"
	"gin-admin-api/internal/middleware"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func InitRoleRouter(Router *gin.RouterGroup, redis *redis.Client, enforcer *casbin.Enforcer, db *gorm.DB, roleHandler role.IRole) {
	authRouter := Router.Group(
		"role",
		middleware.AuthMiddleWare(redis),
		middleware.OperatorMiddleware(),
		middleware.DataPermissionMiddleware(db),
		middleware.CasbinMiddleWare(enforcer),
	)
	authRouter.POST("", roleHandler.CreateRoleApi)                 // 创建角色
	authRouter.DELETE("/:id", roleHandler.DeleteRoleByIdApi)       // 根据id删除角色
	authRouter.PUT("/:id", roleHandler.ModifyRoleByIdApi)          // 根据id修改角色
	authRouter.POST("page", roleHandler.GetRolePageApi)            // 分页获取角色
	authRouter.GET("list", roleHandler.GetRoleListApi)             // 获取角色列表
	authRouter.GET("detail/:id", roleHandler.GetRoleDetailByIdApi) // 获取角色详情
}
