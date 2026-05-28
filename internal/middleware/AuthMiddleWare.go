package middleware

import (
	"gin-admin-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// AuthMiddleWare 中间件校验token登录
func AuthMiddleWare(redisClients ...*redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头中获取token
		tokenString := ctx.GetHeader("token")
		if tokenString == "" {
			utils.Fail(ctx, "必须传递token")
			ctx.Abort()
			return
		}
		// 从token中解析出数据
		token, claims, err := utils.ParseToken(tokenString)
		if err != nil || !token.Valid {
			utils.Fail(ctx, "token解析错误")
			ctx.Abort()
			return
		}
		if len(redisClients) > 0 && redisClients[0] != nil {
			redisDb := utils.NewRedisUtils(redisClients[0])
			if _, err = redisDb.ExistsKey(ctx, utils.AuthTokenRedisKey(claims.AccountId, tokenString)); err != nil {
				utils.Fail(ctx, "token已失效,请重新登录")
				ctx.Abort()
				return
			}
		}
		// 从token中解析出来的数据挂载到上下文上,方便后面的控制器使用
		ctx.Set("accountId", claims.AccountId)
		ctx.Set("userName", claims.Username)
		ctx.Set("isAdmin", claims.IsAdmin)
		ctx.Next()
	}
}
