package dto

import (
	"fmt"
	dto "gin_admin_api/controller/register/dto"
	"github.com/go-playground/validator"
)

var valildate *validator.Validate

func init() {
	valildate = validator.New()
	valildate.RegisterValidation("checkName", dto.CheckNameFunc)
}

type LoginDto struct {
	UserName string `validate:"required,checkName" json:"username"`
	Password string `validate:"required" json:"password"`
}

// 定义校验数据的方法
func ValidatorLogin(account LoginDto) error {
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
