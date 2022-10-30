package validators

import (
	"fmt"
	"gin-admin-api/constants"
	"github.com/go-playground/validator/v10"
	"regexp"
)

// ValidateMobile 定义验证手机号码的校验器
func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	fmt.Println(mobile, "当前手机号码")
	//使用正则表达式判断是否合法
	ok, _ := regexp.MatchString(constants.RegMobile, mobile)
	if !ok {
		return false
	}
	return true
}
