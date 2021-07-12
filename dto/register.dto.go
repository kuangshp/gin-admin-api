package dto

import (
	"gin_admin_api/utils"
	"gin_admin_api/validators"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type RegisterDto struct {
	UserName string `binding:"required,checkName" json:"username"`
	Password string `binding:"required" json:"password"`
}

//注册自定义校验器
func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("checkName", validators.CheckName)
		_ = v.RegisterTranslation("checkName", utils.Trans, func(ut ut.Translator) error {
			return ut.Add("checkName", "{0}用户名不符合规则", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("checkName", fe.Field())
			return t
		})
	}
}
