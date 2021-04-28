package register


import (
"github.com/gin-gonic/gin"
)

func RegisterRouter(router *gin.RouterGroup) {
	router.POST("/", Register)
}
