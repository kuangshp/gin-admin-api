package register

import (
	"fmt"
	"gin_admin_api/common"
	"gin_admin_api/controller/register/dto"
	"gin_admin_api/model"
	"gin_admin_api/response"
	"gin_admin_api/utils"
	"github.com/gin-gonic/gin"
)

// 用户注册账号
func Register(c *gin.Context) {
	// 1.获取前端传递过来的数据
	var registerDto dto.RegisterDto
	err := c.Bind(&registerDto)
	if err != nil {
		response.Fail(c, "解析前端传递的数据错误")
		return
	}
	// 2.对前端传递过来的数据进行校验
	err = dto.ValidatorRegister(registerDto)
	if err != nil {
		response.Fail(c, "数据校验错误")
		return
	}
	// 3.将数据插入到数据库中
	newPassword, err := utils.GeneratePassword(registerDto.Password)
	if err != nil {
		response.Fail(c, "密码加密错误")
		return
	}
	account := model.Account{
		UserName: registerDto.UserName,
		Password: newPassword,
	}
	tx := common.DB.Create(&account)
	fmt.Println(tx.RowsAffected, tx.Error)
	if tx.RowsAffected > 0 {
		response.Success(c, nil)
	} else {
		response.Fail(c, "插入数据错误")
	}
}
