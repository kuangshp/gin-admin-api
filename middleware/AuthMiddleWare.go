package middleware

import (
	"fmt"
	"gin_admin_api/common"
	"gin_admin_api/model"
	"gin_admin_api/response"
	"github.com/gin-gonic/gin"
)

// 中间件校验token登录
func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取token
		tokeString := c.GetHeader("token")
		fmt.Println(tokeString, "当前token")
		if tokeString == "" {
			response.Fail(c, "必须传递token")
			c.Abort()
			return
		}
		// 从token中解析出数据
		token, claims, err := common.ParseToken(tokeString)
		if err != nil || !token.Valid {
			response.Fail(c, "token解析错误")
			c.Abort()
			return
		}
		// 可以进一步查询数据库,是否有当前的用户id
		if first := common.DB.Where("id=?", claims.UserId).First(&model.Account{}); first.Error != nil {
			response.Fail(c, "当前token非法")
			c.Abort()
			return
		}
		// 从token中解析出来的数据挂载到上下文上,方便后面的控制器使用
		c.Set("userId", claims.UserId)
		c.Set("userName", claims.Username)
		c.Next()
	}
}
