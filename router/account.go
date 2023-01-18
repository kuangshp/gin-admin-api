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
	accountRouter.POST("register", newAccount.Register)
	accountRouter.POST("login", newAccount.Login)
	accountRouter.DELETE("/:id", middleware.AuthMiddleWare(), newAccount.DeleteAccountById)
	accountRouter.PUT("/:id", middleware.AuthMiddleWare(), newAccount.ModifyPasswordById)
	accountRouter.PATCH("/modifyPassword", middleware.AuthMiddleWare(), newAccount.UpdateCurrentAccountPassword)
	accountRouter.PATCH("/:id", middleware.AuthMiddleWare(), newAccount.UpdateStatusById)
	accountRouter.GET("/test", newAccount.Test)
	accountRouter.GET("/:id", middleware.AuthMiddleWare(), newAccount.GetAccountById)
	accountRouter.GET("", middleware.AuthMiddleWare(), newAccount.GetAccountPage)
}
