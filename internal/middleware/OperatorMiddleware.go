package middleware

import (
	"context"
	"fmt"
	"gin-admin-api/internal/plugin"
	"github.com/gin-gonic/gin"
)

func OperatorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 JWT claims 或 session 中拿用户名
		accountId, exists := c.Get("accountId")
		fmt.Println("OperatorMiddleware 拿到的 accountId:", accountId, "exists:", exists)
		ctx := context.WithValue(c.Request.Context(), plugin.CtxOperatorKey, accountId)
		ctx = context.WithValue(ctx, plugin.CtxGinContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
