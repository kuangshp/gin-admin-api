package vo

type AccountLoginVO struct {
	ID            int64  `json:"id"`            // 主键id
	Username      string `json:"username"`      // 登录帐号
	Email         string `json:"email"`         // 邮箱
	Mobile        string `json:"mobile"`        // 手机号
	LastLoginDate int64  `json:"lastLoginDate"` // 最后一次登录时间
	LastLoginIP   string `json:"lastLoginIp"`   // 最后一次登录ip
	IsAdmin       int64  `json:"isAdmin"`       // 1是超级管理员，2是普通管理员
	Token         string `json:"token"`         // token
}

type GetCaptchaVO struct {
	Base64    string `json:"base64"`    // base64
	CaptchaId string `json:"captchaId"` // 唯一码
	Code      string `json:"code"`      // 验证码
}
