package api

import (
	"fmt"
	"gin_admin_api/dto"
	"gin_admin_api/global"
	"gin_admin_api/model"
	"gin_admin_api/utils"
	"github.com/gin-gonic/gin"
)

// Register 用户注册账号
func Register(c *gin.Context) {
	// 1.获取前端传递过来的数据
	var registerDto dto.RegisterDto
	if err := c.ShouldBindJSON(&registerDto); err != nil {
		// 2.校验数据是否合法
		message := utils.ShowErrorMessage(err)
		utils.Fail(c, message)
		return
	}

	// 3.将数据插入到数据库中
	newPassword, err := utils.GeneratePassword(registerDto.Password)
	if err != nil {
		utils.Fail(c, "密码加密错误")
		return
	}
	account := model.AccountEntity{
		UserName: registerDto.UserName,
		Password: newPassword,
	}
	fmt.Println(registerDto.UserName, "=====>")
	// 判断当前用户名是否已经存在
	if result := global.DB.Model(&model.AccountEntity{}).Where("username=?", registerDto.UserName).First(&account); result.RowsAffected != 0 {
		fmt.Println("=====>")
		utils.Fail(c, fmt.Sprintf("%s已经存在不能重复注册", account.UserName))
		return
	}
	tx := global.DB.Create(&account)
	fmt.Println(tx.RowsAffected, tx.Error)
	if tx.RowsAffected > 0 {
		utils.Success(c, "注册成功")
	} else {
		utils.Fail(c, "插入数据错误")
	}
}
