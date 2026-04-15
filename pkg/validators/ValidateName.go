package validators

import (
	"github.com/go-playground/validator/v10"
	"unicode/utf8"
)

//ValidateName 自定义校验器校验用户名
func ValidateName(f validator.FieldLevel) bool {
	count := utf8.RuneCountInString(f.Field().String())
	if count >= 3 && count <= 12 {
		return true
	} else {
		return false
	}
}
