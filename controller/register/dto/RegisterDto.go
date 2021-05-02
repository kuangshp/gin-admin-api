package dto

import (
	"fmt"
	"gin_admin_api/model"
	"github.com/go-playground/validator/v10"
	"unicode/utf8"
)
var valildate *validator.Validate

func init() {
	valildate = validator.New()
	valildate.RegisterValidation("checkName", CheckNameFunc)
}

type RegisterDto struct {
	UserName string `binding:"required,checkName" json:"username"`
	Password string `binding:"required" json:"password"`
}

func ToRegisterAccountModel(account RegisterDto) model.Account {
	return model.Account{
		UserName: account.UserName,
		Password: account.Password,
	}
}

// 自定义校验器校验用户名
func CheckNameFunc(f validator.FieldLevel) bool {
	count := utf8.RuneCountInString(f.Field().String())
	if count >= 9 && count <= 12 {
		return true
	} else {
		return false
	}
}

// 定义校验数据的方法
func ValidatorRegister(account RegisterDto) error {
	err := valildate.Struct(account)
	if err != nil {
		// 输出校验错误 .(validator.ValidationErrors)是断言
		for _, e := range err.(validator.ValidationErrors)[:1] {
			fmt.Println("错误字段:", e.Field())
			fmt.Println("错误的值:", e.Value())
			fmt.Println("错误的tag:", e.Tag())
		}
		return err
	} else {
		return nil
	}
}
