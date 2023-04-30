package router

import (
	"gin-admin-api/api/account"
	"gin-admin-api/global"
	"gin-admin-api/middleware"
	"github.com/gin-gonic/gin"
)

func InitAccountRouter(Router *gin.RouterGroup) {
	accountRouter := Router.Group("account")
	newAccount := account.NewAccount(global.DB)
	accountRouter.POST("register", newAccount.Register)                                                                 // 注册
	accountRouter.POST("login", newAccount.Login)                                                                       // 登录
	accountRouter.DELETE("/:id", middleware.AuthMiddleWare(), newAccount.DeleteAccountById)                             // 根据id删除
	accountRouter.PUT("/modifyPassword/:id", middleware.AuthMiddleWare(), newAccount.ModifyPasswordById)                // 根据id修改密码
	accountRouter.PATCH("/modifyPassword/:id", middleware.AuthMiddleWare(), newAccount.ModifyPasswordById)              // 根据id修改密码
	accountRouter.PATCH("/modifyCurrentPassword", middleware.AuthMiddleWare(), newAccount.UpdateCurrentAccountPassword) // 修改当前密码
	accountRouter.PUT("/modifyCurrentPassword", middleware.AuthMiddleWare(), newAccount.UpdateCurrentAccountPassword)   // 修改当前密码
	accountRouter.PATCH("/status/:id", middleware.AuthMiddleWare(), newAccount.UpdateStatusById)                        // 修改状态
	accountRouter.PUT("/status/:id", middleware.AuthMiddleWare(), newAccount.UpdateStatusById)                          // 修改状态
	accountRouter.GET("/:id", middleware.AuthMiddleWare(), newAccount.GetAccountById)                                   // 根据id获取数据
	accountRouter.GET("", middleware.AuthMiddleWare(), newAccount.GetAccountPage)                                       // 分页获取数据
}
