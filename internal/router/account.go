package router

import (
	"gin-admin-api/internal/api/account"
	"gin-admin-api/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func InitAccountRouter(Router *gin.RouterGroup, redis *redis.Client, newAccount account.ISysAccount) {
	registerRouter := Router.Group("account")
	registerRouter.POST("register", middleware.AuthMiddleWare(redis), middleware.OperatorMiddleware(), newAccount.CreateSysAccountApi)               // 创建账号
	registerRouter.DELETE("/:id", middleware.AuthMiddleWare(redis), newAccount.DeleteSysAccountByIdApi)                                              // 根据id删除
	registerRouter.PUT("/modify/:id", middleware.AuthMiddleWare(redis), newAccount.ModifySysAccountByIdApi)                                          // 根据id修改
	registerRouter.PATCH("/modify/:id", middleware.AuthMiddleWare(redis), newAccount.ModifySysAccountByIdApi)                                        // 根据id修改
	registerRouter.PATCH("/modifyPassword/:id", middleware.AuthMiddleWare(redis), newAccount.ResetPasswordByIdApi)                                   // 根据id重置密码
	registerRouter.PUT("/modifyPassword/:id", middleware.AuthMiddleWare(redis), newAccount.ResetPasswordByIdApi)                                     // 根据id重置密码
	registerRouter.PATCH("/modifyCurrentPassword", middleware.AuthMiddleWare(redis), newAccount.ModifyCurrentSysAccountPasswordApi)                  // 修改当前账号密码
	registerRouter.PUT("/modifyCurrentPassword", middleware.AuthMiddleWare(redis), newAccount.ModifyCurrentSysAccountPasswordApi)                    // 修改当前账号密码
	registerRouter.GET("/:id", middleware.AuthMiddleWare(redis), middleware.CacheMiddleWare(redis, "repository"), newAccount.GetSysAccountDetailApi) // 根据id获取数据
	registerRouter.GET("", middleware.AuthMiddleWare(redis), middleware.CacheMiddleWare(redis, "repository"), newAccount.GetSysAccountPageApi)       // 分页获取数据
}
