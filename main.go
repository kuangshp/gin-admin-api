package main

import (
	"fmt"
	"gin_admin_api/common"
	_ "gin_admin_api/common"
	"gin_admin_api/model"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func init()  {
	common.DB.AutoMigrate(&model.Account{})
}

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
		})
	})
	port := viper.GetString("server.port")
	fmt.Println("当前端口", port)
	if port != "" {
		router.Run(":" + port)
	} else {
		router.Run()
	}
}


