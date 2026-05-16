package router

import (
	"gin-admin-api/internal/api/account"
	"gin-admin-api/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func InitAccountRouter(Router *gin.RouterGroup, redis *redis.Client, newAccount account.IAccount, signingKey []byte) {
	registerRouter := Router.Group("account")
	auth := middleware.AuthMiddleWare(signingKey)
	registerRouter.POST("register", auth, middleware.OperatorMiddleware(), newAccount.CreateAccountApi)                  // 创建账号
	registerRouter.POST("login", newAccount.LoginAccountApi)                                                           // 登录
	registerRouter.DELETE("/:id", auth, newAccount.DeleteAccountByIdApi)                                               // 根据id删除
	registerRouter.PUT("/modifyPassword/:id", auth, newAccount.ModifyPasswordByIdApi)                                  // 根据id修改密码
	registerRouter.PATCH("/modifyPassword/:id", auth, newAccount.ModifyPasswordByIdApi)                                // 根据id修改密码
	registerRouter.PATCH("/modifyCurrentPassword", auth, newAccount.UpdateCurrentAccountPasswordApi)                   // 修改当前账号密码
	registerRouter.PUT("/modifyCurrentPassword", auth, newAccount.UpdateCurrentAccountPasswordApi)                     // 修改当前账号密码
	registerRouter.PATCH("/status/:id", auth, newAccount.UpdateStatusByIdApi)                                          // 根据id修改状态
	registerRouter.PUT("/status/:id", auth, middleware.OperatorMiddleware(), newAccount.UpdateStatusByIdApi)           // 根据id修改状态
	registerRouter.GET("/:id", auth, middleware.CacheMiddleWare(redis, "repository"), newAccount.GetAccountByIdApi)    // 根据id获取数据
	registerRouter.GET("", auth, middleware.CacheMiddleWare(redis, "repository"), newAccount.GetAccountPageApi)        // 分页获取数据
}
