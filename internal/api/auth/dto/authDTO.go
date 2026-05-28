package dto

type AccountLoginDTO struct {
	Username     string `json:"username" validate:"required"`     // 登录帐号
	Password     string `json:"password" validate:"required"`     // 登录密码
	CaptchaId    string `json:"captchaId" validate:"required"`    // 图形验证码id
	CaptchaValue string `json:"captchaValue" validate:"required"` // 图形验证码
}

type VerifyCaptchaDTO struct {
	CaptchaId    string `json:"captchaId" validate:"required"`    // 图形验证码id
	CaptchaValue string `json:"captchaValue" validate:"required"` // 图形验证码
}
