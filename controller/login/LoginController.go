package login

import (
	"fmt"
	"gin_admin_api/common"
	dto "gin_admin_api/controller/login/dto"
	"gin_admin_api/model"
	"gin_admin_api/response"
	"gin_admin_api/utils"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var loginDto dto.LoginDto
	// 解析前端传递过来的数据并且验证是否正确
	if err := c.ShouldBindJSON(&loginDto);err != nil {
		message := common.ShowErrorMessage(err)
		response.Fail(c, message)
		return
	}
	// 查询数据库登录操作
	account := model.Account{}
	first := common.DB.Where("user_name=?", loginDto.UserName).First(&account)
	if first.Error == nil {
		// 对账号和密码校验
		if isOk, _ := utils.CheckPassword(account.Password, loginDto.Password); isOk {
			// 生成token返回给前端
			hmacUser := common.HmacUser{
				Id: int(account.ID),
				Username: account.UserName,
			}
			if token, err := common.GenerateToken(hmacUser); err == nil {
				response.Success(c, gin.H{
					"token": token,
					"username": account.UserName,
				})
			}
		} else {
			response.Fail(c, "账号或密码错误")
		}
	} else {
		response.Fail(c, "账号不存在")
	}
	fmt.Println(first.Error)
}
