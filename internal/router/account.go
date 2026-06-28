package router

import (
	"gin-admin-api/internal/api/account"
	"gin-admin-api/internal/middleware"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func InitAccountRouter(Router *gin.RouterGroup, redis *redis.Client, enforcer *casbin.Enforcer, db *gorm.DB, newAccount account.ISysAccount) {
	registerRouter := Router.Group(
		"account",
		middleware.AuthMiddleWare(redis),
		middleware.OperatorMiddleware(),
		middleware.DataPermissionMiddleware(db),
		middleware.CasbinMiddleWare(enforcer),
	)
	registerRouter.POST("", newAccount.CreateSysAccountApi)                                      // 创建账号
	registerRouter.DELETE("/:id", newAccount.DeleteSysAccountByIdApi)                            // 根据id删除
	registerRouter.PUT("/:id", newAccount.ModifySysAccountByIdApi)                               // 根据id修改
	registerRouter.POST("/resetPassword", newAccount.ResetPasswordByIdApi)                       // 根据id重置密码
	registerRouter.POST("/modifyCurrentPassword", newAccount.ModifyCurrentSysAccountPasswordApi) // 修改当前账号密码
	registerRouter.GET("/info", newAccount.GetCurrentSysAccountInfoApi)                          // 获取当前登录账号信息
	registerRouter.GET("/list", newAccount.GetSysAccountListApi)                                 // 获取数据列表
	registerRouter.GET("/detail/:id", newAccount.GetSysAccountDetailApi)                         // 根据id获取数据
	registerRouter.POST("page", newAccount.GetSysAccountPageApi)                                 // 分页获取数据
}
