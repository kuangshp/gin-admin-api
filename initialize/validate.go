package initialize

import (
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	"gin_admin_api/utils"
	"gin_admin_api/validators"
)

func InitValidate() {
	registerValidateMobile()
	registerValidateName()
}

func registerValidateMobile() {
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
}

func registerValidateName() {
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
