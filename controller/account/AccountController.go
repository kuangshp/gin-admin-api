package account

import (
	"fmt"
	"gin_admin_api/common"
	"gin_admin_api/dto"
	"gin_admin_api/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 用户注册账号
func Register(c *gin.Context) {
	// 1.获取前端传递过来的数据
	var accountDto dto.AccountDto
	err := c.Bind(&accountDto)
	if err != nil {
		response.Fail(c, "解析前端传递的数据错误")
		return
	}
	// 2.对前端传递过来的数据进行校验
	err = dto.ValidatorAccount(accountDto)
	if err != nil {
		response.Fail(c, "数据校验错误")
		return
	}
	// 3.将数据插入到数据库中
	account := dto.ToAccountModel(accountDto)
	tx := common.DB.Create(&account)
	fmt.Println(tx.RowsAffected, tx.Error)
	if tx.RowsAffected > 0 {
		response.Success(c, gin.H{})
	} else {
		response.Fail(c, "插入数据错误")
	}
}

func AccountList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "账号列表",
	})
}
