package initialize

import (
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	"gin-admin-api/utils"
	"gin-admin-api/validators"
)

func InitValidate() {
	//注册手机号码验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", validators.ValidateMobile)
		_ = v.RegisterTranslation("mobile", utils.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	// 注册email的验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("myEmail", validators.ValidateEmail)
		_ = v.RegisterTranslation("myEmail", utils.Trans, func(ut ut.Translator) error {
			return ut.Add("myEmail", "{0} 邮箱号码非法!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("myEmail", fe.Field())
			return t
		})
	}
	// 注册年月日的验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("date", validators.ValidateDate)
		_ = v.RegisterTranslation("date", utils.Trans, func(ut ut.Translator) error {
			return ut.Add("date", "{0} 时间格式非法必须为2006-01-02格式!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("date", fe.Field())
			return t
		})
	}
	// 注册年月日时分秒的验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("dateTime", validators.ValidateDateTime)
		_ = v.RegisterTranslation("dateTime", utils.Trans, func(ut ut.Translator) error {
			return ut.Add("dateTime", "{0} 时间格式非法必须为2006-01-02 15:04:05格式!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("dateTime", fe.Field())
			return t
		})
	}

	// 校对用户名
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("checkName", validators.ValidateName)
		_ = v.RegisterTranslation("checkName", utils.Trans, func(ut ut.Translator) error {
			return ut.Add("checkName", "{0}用户名长度必须是3-12位字符", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("checkName", fe.Field())
			return t
		})
	}
}
