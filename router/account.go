package router

import (
	"gin-admin-api/api/account"
	"gin-admin-api/global"
	"gin-admin-api/middleware"
	"github.com/gin-gonic/gin"
)

func InitAccountRouter(Router *gin.RouterGroup) {
	registerRouter := Router.Group("account")
	newAccount := account.NewAccount(&global.DB)
	registerRouter.POST("register", newAccount.CreateAccountApi)                                                            // 创建账号
	registerRouter.POST("login", newAccount.LoginAccountApi)                                                                // 登录
	registerRouter.DELETE("/:id", middleware.AuthMiddleWare(), newAccount.DeleteAccountByIdApi)                             // 根据id删除
	registerRouter.PUT("/modifyPassword/:id", middleware.AuthMiddleWare(), newAccount.ModifyPasswordByIdApi)                // 根据id修改密码
	registerRouter.PATCH("/modifyPassword/:id", middleware.AuthMiddleWare(), newAccount.ModifyPasswordByIdApi)              // 根据id修改密码
	registerRouter.PATCH("/modifyCurrentPassword", middleware.AuthMiddleWare(), newAccount.UpdateCurrentAccountPasswordApi) // 修改当前账号密码
	registerRouter.PUT("/modifyCurrentPassword", middleware.AuthMiddleWare(), newAccount.UpdateCurrentAccountPasswordApi)   // 修改当前账号密码
	registerRouter.PATCH("/status/:id", middleware.AuthMiddleWare(), newAccount.UpdateStatusByIdApi)                        // 根据id修改状态
	registerRouter.PUT("/status/:id", middleware.AuthMiddleWare(), newAccount.UpdateStatusByIdApi)                          // 根据id修改状态
	registerRouter.GET("/:id", middleware.AuthMiddleWare(), newAccount.GetAccountByIdApi)                                   // 根据id获取数据
	//accountRouter.GET("", middleware.AuthMiddleWare(), newAccount.GetAccountPage)                                       // 分页获取数据
}
