package account

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateAccount(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
	})
}
