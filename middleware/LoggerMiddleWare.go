package middleware

import (
	"fmt"
	"gin-admin-api/global"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		//请求方式
		method := c.Request.Method
		//请求路由
		reqUrl := c.Request.RequestURI
		//状态码
		statusCode := c.Writer.Status()
		//请求ip
		clientIP := c.ClientIP()
		// 打印日志
		//loggerMap := map[string]interface{} {
		//	"status_code":statusCode,
		//	"client_ip": clientIP,
		//	"req_method":method,
		//	"req_uri": reqUrl,
		//}
		//marshal, _ := json.Marshal(loggerMap)
		loggerStr := fmt.Sprintf("status_code:%d,client_ip:%s,req_method:%s,req_uri:%s", statusCode, clientIP, method, reqUrl)
		global.Logger.Info("中间件本次请求", zap.String("http", loggerStr))

	}
}
