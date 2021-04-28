package account

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AccountList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "账号列表",
	})
}
