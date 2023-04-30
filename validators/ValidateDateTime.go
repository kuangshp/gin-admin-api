package validators

import (
	"gin-admin-api/constants"
	"github.com/go-playground/validator/v10"
	"regexp"
)

// ValidateDateTime 定义验证时间格式(年月日时分秒)
func ValidateDateTime(fl validator.FieldLevel) bool {
	dateTime := fl.Field().String()
	//使用正则表达式判断是否合法
	ok, _ := regexp.MatchString(constants.RegDateTime, dateTime)
	if !ok {
		return false
	}
	return true
}
