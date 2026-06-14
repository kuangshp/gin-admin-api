package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	gormplus "github.com/kuangshp/gorm-plus"
)

func OperatorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 JWT claims 或 session 中拿用户名
		accountId, exists := c.Get("accountId")
		fmt.Println("OperatorMiddleware 拿到的 accountId:", accountId, "exists:", exists)
		ctx := context.WithValue(c.Request.Context(), gormplus.CtxContextKey1, accountId)
		ctx = context.WithValue(ctx, gormplus.CtxContextKey1, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
