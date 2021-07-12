package api

import (
	"gin_admin_api/dto"
	"gin_admin_api/response"
	"gin_admin_api/utils"
	"github.com/gin-gonic/gin"
)

// Register 用户注册账号
func Register(c *gin.Context) {
	// 1.获取前端传递过来的数据
	var registerDto dto.RegisterDto
	if err := c.ShouldBindJSON(&registerDto);err != nil {
		// 2.校验数据是否合法
		message := utils.ShowErrorMessage(err)
		response.Fail(c, message)
		return
	}
	//// 3.将数据插入到数据库中
	//newPassword, err := utils.GeneratePassword(registerDto.Password)
	//if err != nil {
	//	response.Fail(c, "密码加密错误")
	//	return
	//}
	//account := model.AccountEntity{
	//	UserName: registerDto.UserName,
	//	Password: newPassword,
	//}
	//tx := global.DB.Create(&account)
	//fmt.Println(tx.RowsAffected, tx.Error)
	//if tx.RowsAffected > 0 {
	//	response.Success(c, nil)
	//} else {
	//	response.Fail(c, "插入数据错误")
	//}
}
