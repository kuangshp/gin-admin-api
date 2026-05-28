package router

import (
	"gin-admin-api/internal/api/auth"
	"gin-admin-api/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func InitAuthRouter(Router *gin.RouterGroup, redis *redis.Client, authHandler auth.IAuth) {
	authRouter := Router.Group("auth")
	authRouter.POST("login", authHandler.AccountLoginApi)                              // 登录
	authRouter.POST("logout", middleware.AuthMiddleWare(redis), authHandler.LogOutApi) // 退出登录
	authRouter.GET("captcha", authHandler.GetCaptchaApi)                               // 获取图形验证码
	authRouter.POST("captcha/verify", authHandler.VerifyCaptchaApi)                    // 验证图形验证码
}
