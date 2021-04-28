package login

import (
	"fmt"
	"gin_admin_api/common"
	"gin_admin_api/dto"
	"gin_admin_api/model"
	"gin_admin_api/response"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var loginDto dto.LoginDto
	err := c.BindJSON(&loginDto)
	if err != nil {
		fmt.Println(err)
		response.Fail(c, "数据解析错误")
		return
	}
	// 查询数据库登录操作
	first := common.DB.Where("user_name=?", loginDto.UserName).Where("password=?", loginDto.Password).First(&model.Account{})
	if first.Error == nil {
		response.Success(c, nil)
		return
	} else {
		response.Fail(c, "登录失败")
	}
	fmt.Println(first.Error)
}
