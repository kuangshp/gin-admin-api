package dto


type RegisterDto struct {
	UserName string `binding:"required" json:"username"`
	Password string `binding:"required" json:"password"`
}


// CheckNameFunc 自定义校验器校验用户名
//func CheckNameFunc(f validator.FieldLevel) bool {
//	count := utf8.RuneCountInString(f.Field().String())
//	if count >= 9 && count <= 12 {
//		return true
//	} else {
//		return false
//	}
//}

// 注册自定义校验器
//func init() {
//	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
//		_ = v.RegisterValidation("checkName", CheckNameFunc)
//		_ = v.RegisterTranslation("checkName", utils.Trans, func(ut ut.Translator) error {
//			return ut.Add("checkName", "{0}用户名不符合规则", true)
//		}, func(ut ut.Translator, fe validator.FieldError) string {
//			t, _ := ut.T("checkName", fe.Field())
//			return t
//		})
//	}
//}