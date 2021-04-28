package login

import (
	"fmt"
	"gin_admin_api/common"
	"gin_admin_api/dto"
	"gin_admin_api/model"
	"gin_admin_api/response"
	"gin_admin_api/utils"
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
	// 校验数据是否合法
	err = dto.ValidatorLogin(loginDto)
	if err != nil {
		response.Fail(c, "参数校验失败")
		return
	}
	// 查询数据库登录操作
	account := model.Account{}
	first := common.DB.Where("user_name=?", loginDto.UserName).First(&account)
	if first.Error == nil {
		// 对账号和密码校验
		if isOk, _ := utils.CheckPassword(account.Password, loginDto.Password); isOk {
			response.Success(c, nil)
		} else {
			response.Fail(c, "账号或密码错误")
		}
	} else {
		response.Fail(c, "账号不存在")
	}
	fmt.Println(first.Error)
}
