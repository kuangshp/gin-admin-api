package middleware

import "github.com/gin-gonic/gin"

func CacheMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
