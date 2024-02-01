package router

import (
	"gin-admin-api/api/account"
	"gin-admin-api/middleware"
	"github.com/gin-gonic/gin"
	"time"
)

func InitAccountRouter(Router *gin.RouterGroup) {
	registerRouter := Router.Group("account")
	newAccount := account.NewAccount()
	registerRouter.POST("register", newAccount.CreateAccountApi)                                                                                                                                // 创建账号
	registerRouter.POST("login", newAccount.LoginAccountApi)                                                                                                                                    // 登录
	registerRouter.DELETE("/:id", middleware.AuthMiddleWare(), newAccount.DeleteAccountByIdApi)                                                                                                 // 根据id删除
	registerRouter.PUT("/modifyPassword/:id", middleware.AuthMiddleWare(), newAccount.ModifyPasswordByIdApi)                                                                                    // 根据id修改密码
	registerRouter.PATCH("/modifyPassword/:id", middleware.AuthMiddleWare(), newAccount.ModifyPasswordByIdApi)                                                                                  // 根据id修改密码
	registerRouter.PATCH("/modifyCurrentPassword", middleware.AuthMiddleWare(), newAccount.UpdateCurrentAccountPasswordApi)                                                                     // 修改当前账号密码
	registerRouter.PUT("/modifyCurrentPassword", middleware.AuthMiddleWare(), newAccount.UpdateCurrentAccountPasswordApi)                                                                       // 修改当前账号密码
	registerRouter.PATCH("/status/:id", middleware.AuthMiddleWare(), newAccount.UpdateStatusByIdApi)                                                                                            // 根据id修改状态
	registerRouter.PUT("/status/:id", middleware.AuthMiddleWare(), newAccount.UpdateStatusByIdApi)                                                                                              // 根据id修改状态
	registerRouter.GET("/:id", middleware.AuthMiddleWare(), middleware.CacheMiddleWare("account"), newAccount.GetAccountByIdApi)                                                                // 根据id获取数据
	registerRouter.GET("", middleware.AuthMiddleWare(), middleware.RedisRequestLockRequestLock("account", 10*time.Second), middleware.CacheMiddleWare("account"), newAccount.GetAccountPageApi) // 分页获取数据
}
