package dto

import (
	"gin_admin_api/model"
	"github.com/go-playground/validator"
	"unicode/utf8"
)

type AccountDto struct {
	UserName string `validate:"checkName,required" json:"username"`
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

// 自定义校验器校验用户名
func checkNameFunc(f validator.FieldLevel) bool {
	count := utf8.RuneCountInString(f.Field().String())
	if count >= 6 && count <= 12 {
		return true
	} else {
		return false
	}
}

func init() {
	valildate := validator.New()
	valildate.RegisterValidation("checkName", checkNameFunc)
}
