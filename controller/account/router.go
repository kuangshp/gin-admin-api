package account

import (
	"github.com/gin-gonic/gin"
)

func AccountRouter(router *gin.RouterGroup) {
	router.GET("/", AccountList)
	router.GET("/:id", AccountById)
}
