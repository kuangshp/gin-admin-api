package api

import (
	"fmt"
	"gin_admin_api/dto"
	"gin_admin_api/global"
	"gin_admin_api/model"
	"gin_admin_api/utils"
	"github.com/gin-gonic/gin"
)

// Login @Summary 用户登录接口
// @Tags 用户登录
// @title 输入用户名和密码登录系统
// @version 1.0
// @description 用户登录
// @Produce json
// @Param loginDto body dto.LoginDto true "用户登录参数"
// @Success 200 {object} Res {"code":0,"data":null,"message":"操作成功"}
// @Router /login [post]
func Login(c *gin.Context) {
	var loginDto dto.LoginDto
	// 解析前端传递过来的数据并且验证是否正确
	if err := c.ShouldBindJSON(&loginDto); err != nil {
		message := utils.ShowErrorMessage(err)
		utils.Fail(c, message)
		return
	}
	// 查询数据库登录操作
	account := model.AccountEntity{}
	first := global.DB.Where("username=?", loginDto.UserName).First(&account)
	if first.Error == nil {
		// 对账号和密码校验
		if isOk, _ := utils.CheckPassword(account.Password, loginDto.Password); isOk {
			// 生成token返回给前端
			hmacUser := utils.HmacUser{
				Id:       int(account.ID),
				Username: account.UserName,
			}
			if token, err := utils.GenerateToken(hmacUser); err == nil {
				utils.Success(c, gin.H{
					"token":    token,
					"username": account.UserName,
				})
			}
		} else {
			utils.Fail(c, "账号或密码错误")
		}
	} else {
		utils.Fail(c, "账号不存在")
	}
	fmt.Println(first.Error)
}
