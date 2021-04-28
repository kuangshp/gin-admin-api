package dto

import (
	"fmt"
	"gin_admin_api/model"
	"github.com/go-playground/validator"
	"unicode/utf8"
)

type AccountDto struct {
	UserName string `validate:"required,checkName" json:"username"`
	Password string `validate:"required" json:"password"`
	Mobile   string `validate:"min=11,max=11" json:"mobile"`
}

func ToAccountDto(account model.Account) AccountDto {
	return AccountDto{
		UserName: account.UserName,
		Password: account.Password,
		Mobile:   account.Mobile,
	}
}

func ToAccountModel(account AccountDto) model.Account {
	return model.Account{
		UserName: account.UserName,
		Password: account.Password,
		Mobile:   account.Mobile,
	}
}

// 自定义校验器校验用户名
func checkNameFunc(f validator.FieldLevel) bool {
	count := utf8.RuneCountInString(f.Field().String())
	if count >= 2 && count <= 12 {
		return true
	} else {
		return false
	}
}

var valildate *validator.Validate

func init() {
	valildate = validator.New()
	valildate.RegisterValidation("checkName", checkNameFunc)
}

// 定义校验数据的方法
func ValidatorAccount(account AccountDto) error {
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
