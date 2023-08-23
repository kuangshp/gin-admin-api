package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Response(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
		"result":  data,
	})
}

// Success 成功的请求
func Success(ctx *gin.Context, data interface{}) {
	Response(ctx, 0, "请求成功", data)
}

// Fail 失败的请求
func Fail(ctx *gin.Context, message string) {
	Response(ctx, 1, message, nil)
}

type PageVo struct {
	Data       interface{} `json:"data"`       // 数据
	Total      int64       `json:"total"`      // 总条数
	PageSize   int64       `json:"pageSize"`   // 当前条数
	PageNumber int64       `json:"pageNumber"` // 当前页数
}

// BuildPageData 构造分页查询器
func BuildPageData(ctx *gin.Context, data interface{}, total int64) {
	size, number := GetQueryPage(ctx.Request)
	Success(ctx, PageVo{
		Data:       If(data == nil, make([]interface{}, 0), data),
		Total:      total,
		PageSize:   size,
		PageNumber: number,
	})
}
