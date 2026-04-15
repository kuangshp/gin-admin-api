package utils

import "github.com/gin-gonic/gin"

// GetCurrentIP 获取当前请求ip
func GetCurrentIP(ctx *gin.Context) string {
	ip := ctx.Request.Header.Get("X-Real-IP")
	if ip == "" {
		ip = ctx.Request.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}
	return ip
}
