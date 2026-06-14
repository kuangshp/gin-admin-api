package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	gormplus "github.com/kuangshp/gorm-plus"
)

func OperatorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 JWT claims 或 session 中拿用户名
		accountIDValue, exists := c.Get("accountId")
		if !exists {
			c.Next()
			return
		}
		accountId, ok := accountIDValue.(int64)
		if !ok {
			c.Next()
			return
		}
		ctx := context.WithValue(c.Request.Context(), gormplus.CtxContextKey1, accountId)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
