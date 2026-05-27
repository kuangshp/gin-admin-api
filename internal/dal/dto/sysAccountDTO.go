package dto

// CreateSysAccountDTO 创建后台账号表请求
type CreateSysAccountDTO struct {
	Username        string `json:"username" validate:"required"`                                     // 登录帐号
	Email           string `json:"email"`                                                            // 邮箱
	Mobile          string `json:"mobile"`                                                           // 手机号
	Password        string `json:"password" validate:"required,min=6,max=16"`                        // 登录密码
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=6,max=16,eqfield=Password"` // 确定密码
	Avatar          string `json:"avatar"`                                                           // 头像

}

// ModifySysAccountDTO 修改后台账号表请求
type ModifySysAccountDTO struct {
	ID       int64  `json:"id" validate:"required,number,gte=1"`  // 主键id
	Username string `json:"username" validate:"required"`         // 登录帐号
	Email    string `json:"email"`                                // 邮箱
	Mobile   string `json:"mobile"`                               // 手机号
	Status   int64  `json:"status" validate:"required,oneof=1 2"` // 状态1是正常,2是禁用
	Avatar   string `json:"avatar"`                               // 头像
}

type ResetPasswordDTO struct {
	Password        string `json:"password" binding:"required,min=6,max=16"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=6,max=16"`
}

type ModifyCurrentPasswordDTO struct {
	Password        string `json:"password" binding:"required,min=6,max=16"`
	NewPassword     string `json:"newPassword" binding:"required,min=6,max=16"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=6,max=16,eqfield=Password"`
}

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
