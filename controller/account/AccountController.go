package account

import (
	"fmt"
	"gin_admin_api/common"
	"gin_admin_api/model"
	"gin_admin_api/response"
	"github.com/gin-gonic/gin"
)

// 查询全部的账号信息
func AccountList(c *gin.Context) {
	// 定义一个切片来存储查询出来的数据
	account := make([]model.Account, 10)
	tx := common.DB.Select([]string{"id", "user_name", "mobile", "created_at", "updated_at"}).Find(&account)
	if tx.Error != nil {
		response.Fail(c, "查询数据错误")
		fmt.Println(tx.Error)
		return
	} else {
		response.Success(c, gin.H{
			"data": account,
		})
	}
}