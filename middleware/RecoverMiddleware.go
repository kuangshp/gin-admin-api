package middleware

import (
	"fmt"
	"gin_admin_api/response"
	"github.com/gin-gonic/gin"
)

func RecoverMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				response.Fail(c, fmt.Sprint(err))
				c.Abort()
				return
			}
		}()
	}
}
