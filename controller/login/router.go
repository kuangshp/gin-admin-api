package login

import (
"github.com/gin-gonic/gin"
)

func LoginRouter(router *gin.RouterGroup) {
	router.POST("/", Login)
}
