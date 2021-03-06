package api

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"gin_admin_api/global"
	"gin_admin_api/model"
	"gin_admin_api/utils"
	"gin_admin_api/vo"
)

// AccountById 根据id查询数据
func AccountById(c *gin.Context) {
	// 获取参数
	id := c.Param("id")
	// 根据id查询数据
	account := &model.AccountEntity{}
	tx := global.DB.Select("id", "user_name", "mobile", "created_at", "updated_at").Where("id=?", id).First(&account)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		utils.Fail(c, "查询错误")
		return
	}
	// 成功返回
	utils.Success(c, vo.ToAccountModelToRes(*account))
}

// AccountList 查询全部的账号信息
func AccountList(c *gin.Context) {
	fmt.Println(global.ServerConfig.DataSource, "测试")
	// 定义一个切片来存储查询出来的数据
	account := make([]model.AccountEntity, 10)
	var total int64
	//tx := global.DB.Select([]string{"id", "username",  "created_at", "updated_at"}).Find(&account).Count(&total)
	tx := global.DB.Find(&account).Count(&total)
	if tx.Error != nil {
		utils.Fail(c, "查询数据错误")
		fmt.Println(tx.Error)
		return
	} else {
		res := vo.ToAccountModelListToRes(account)
		utils.Success(c, gin.H{
			"data": res,
			"total": total,
		})
	}
}
