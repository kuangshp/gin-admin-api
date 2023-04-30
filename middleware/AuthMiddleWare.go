package middleware

import (
	"gin-admin-api/global"
	"gin-admin-api/model"
	"gin-admin-api/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

// AuthMiddleWare 中间件校验token登录
func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取token
		tokeString := c.GetHeader("token")
		fmt.Println(tokeString, "当前token")
		if tokeString == "" {
			utils.Fail(c, "必须传递token")
			c.Abort()
			return
		}
		// 从token中解析出数据
		token, claims, err := utils.ParseToken(tokeString)
		if err != nil || !token.Valid {
			utils.Fail(c, "token解析错误")
			c.Abort()
			return
		}
		// 可以进一步查询数据库,是否有当前的用户id
		if first := global.DB.Where("id=?", claims.UserId).First(&model.AccountEntity{}); first.Error != nil {
			utils.Fail(c, "当前token非法")
			c.Abort()
			return
		}
		// 从token中解析出来的数据挂载到上下文上,方便后面的控制器使用
		c.Set("accountId", claims.UserId)
		c.Set("userName", claims.Username)
		c.Next()
	}
}
