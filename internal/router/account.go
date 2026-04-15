package router

import (
	"gin-admin-api/internal/api/account"
	"gin-admin-api/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func InitAccountRouter(Router *gin.RouterGroup, redis *redis.Client, newAccount account.IAccount) {
	registerRouter := Router.Group("repository")
	registerRouter.POST("register", newAccount.CreateAccountApi)                                                                           // 创建账号
	registerRouter.POST("login", newAccount.LoginAccountApi)                                                                               // 登录
	registerRouter.DELETE("/:id", middleware.AuthMiddleWare(), newAccount.DeleteAccountByIdApi)                                            // 根据id删除
	registerRouter.PUT("/modifyPassword/:id", middleware.AuthMiddleWare(), newAccount.ModifyPasswordByIdApi)                               // 根据id修改密码
	registerRouter.PATCH("/modifyPassword/:id", middleware.AuthMiddleWare(), newAccount.ModifyPasswordByIdApi)                             // 根据id修改密码
	registerRouter.PATCH("/modifyCurrentPassword", middleware.AuthMiddleWare(), newAccount.UpdateCurrentAccountPasswordApi)                // 修改当前账号密码
	registerRouter.PUT("/modifyCurrentPassword", middleware.AuthMiddleWare(), newAccount.UpdateCurrentAccountPasswordApi)                  // 修改当前账号密码
	registerRouter.PATCH("/status/:id", middleware.AuthMiddleWare(), newAccount.UpdateStatusByIdApi)                                       // 根据id修改状态
	registerRouter.PUT("/status/:id", middleware.AuthMiddleWare(), newAccount.UpdateStatusByIdApi)                                         // 根据id修改状态
	registerRouter.GET("/:id", middleware.AuthMiddleWare(), middleware.CacheMiddleWare(redis, "repository"), newAccount.GetAccountByIdApi) // 根据id获取数据
	registerRouter.GET("", middleware.AuthMiddleWare(), middleware.CacheMiddleWare(redis, "repository"), newAccount.GetAccountPageApi)     // 分页获取数据
}
