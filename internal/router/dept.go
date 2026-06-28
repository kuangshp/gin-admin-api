package router

import (
	"gin-admin-api/internal/api/dept"
	"gin-admin-api/internal/middleware"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func InitDeptRouter(Router *gin.RouterGroup, redis *redis.Client, enforcer *casbin.Enforcer, db *gorm.DB, deptHandler dept.IDept) {
	authRouter := Router.Group(
		"dept",
		middleware.AuthMiddleWare(redis),
		middleware.OperatorMiddleware(),
		middleware.DataPermissionMiddleware(db),
		middleware.CasbinMiddleWare(enforcer),
	)
	authRouter.POST("", deptHandler.CreateDeptApi)
	authRouter.DELETE("/:id", deptHandler.DeleteDeptByIdApi)
	authRouter.PUT("/:id", deptHandler.ModifyDeptByIdApi)
	authRouter.POST("page", deptHandler.GetDeptPageApi)
	authRouter.GET("list", deptHandler.GetDeptListApi)
	authRouter.GET("detail/:id", deptHandler.GetDeptDetailByIdApi)
}
