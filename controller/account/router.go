package account

import (
	"github.com/gin-gonic/gin"
)

func AccountRouter(router *gin.RouterGroup) {
	router.POST("/", Register)
	router.GET("/", AccountList)
}
