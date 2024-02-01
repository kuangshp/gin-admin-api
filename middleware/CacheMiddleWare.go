package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gin-admin-api/global"
	"gin-admin-api/utils"
	"github.com/gin-gonic/gin"
	"io"
)

// 自定义一个结构体，实现 gin.ResponseWriter interface
type responseWriter struct {
	gin.ResponseWriter
	b *bytes.Buffer
}

// 重写 Write([]byte) (int, error) 方法
func (w responseWriter) Write(b []byte) (int, error) {
	//向一个bytes.buffer中写一份数据来为获取body使用
	w.b.Write(b)
	//完成gin.Context.Writer.Write()原有功能
	return w.ResponseWriter.Write(b)
}

func CacheMiddleWare(prefix string, expiration ...int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 根据当前请求参数来设置redis的key,注意区分请求方式
		method := ctx.Request.Method
		var redisKey = ""
		if method == "POST" {
			var bodyData map[string]interface{}
			body, err := io.ReadAll(ctx.Request.Body)
			// 等于拷贝一份往下传递,否则下接口的方法中拿不到请求体数据
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			if err != nil {
				fmt.Println(err, "????")
			}
			if err = json.Unmarshal(body, &bodyData); err != nil {
				fmt.Println("序列化失败", err)
			}
			fmt.Println("拼接后字符串:", utils.MapKeySort(bodyData))
			fmt.Println("加密后:", utils.Md5(utils.MapKeySort(bodyData)))
			redisKey = fmt.Sprintf("%s_%s", prefix, utils.Md5(utils.MapKeySort(bodyData)))
		} else if method == "GET" {
			queryParams := ctx.Request.URL.Query()
			fmt.Println("拼接后字符串:", utils.MapKeySort(queryParams))
			fmt.Println("加密后:", utils.Md5(utils.MapKeySort(queryParams)))
			redisKey = fmt.Sprintf("%s_%s", prefix, utils.Md5(utils.MapKeySort(queryParams)))
		}
		redisDb := utils.NewRedisUtils(global.RedisDb)
		redisData, _ := redisDb.GetRedisValue(ctx, redisKey)
		fmt.Println("redis缓存数据1111:", redisData)
		if redisData != "" {
			fmt.Println("进来了1")
			var responseData map[string]interface{}
			err := json.Unmarshal([]byte(redisData), &responseData)
			if err != nil {
				ctx.Next()
				return
			}
			utils.Success(ctx, responseData)
			ctx.Abort()
		} else {
			fmt.Println("进来了2")
			writer := responseWriter{
				ctx.Writer,
				bytes.NewBuffer([]byte{}),
			}
			ctx.Writer = writer
			ctx.Next()
			newExpiration := 5000
			if len(expiration) > 0 {
				newExpiration = int(expiration[0])
			}
			err := redisDb.SetRedisValue(ctx, redisKey, writer.b.String(), int64(newExpiration))
			if err != nil {
				return
			}
			fmt.Println("请求返回数据:", writer.b.String())
		}
	}
}
