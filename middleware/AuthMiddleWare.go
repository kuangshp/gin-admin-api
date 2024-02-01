package middleware

import (
	"fmt"
	"gin-admin-api/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleWare 中间件校验token登录
func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头中获取token
		tokeString := ctx.GetHeader("token")
		fmt.Println(tokeString, "当前token")
		if tokeString == "" {
			utils.Fail(ctx, "必须传递token")
			ctx.Abort()
			return
		}
		// 从token中解析出数据
		token, claims, err := utils.ParseToken(tokeString)
		if err != nil || !token.Valid {
			utils.Fail(ctx, "token解析错误")
			ctx.Abort()
			return
		}
		fmt.Println(utils.MapToJson(claims), "解析出来")
		// 从token中解析出来的数据挂载到上下文上,方便后面的控制器使用
		ctx.Set("accountId", claims.AccountId)
		ctx.Set("userName", claims.Username)
		ctx.Next()
	}
}
