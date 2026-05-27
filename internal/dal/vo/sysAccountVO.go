package vo

// SysAccountVO 后台账号表视图对象
type SysAccountVO struct {
	ID            int64  `json:"id"`            // 主键id
	Username      string `json:"username"`      // 登录帐号
	Email         string `json:"email"`         // 邮箱
	Mobile        string `json:"mobile"`        // 手机号
	Password      string `json:"password"`      // 登录密码
	LastLoginDate int64  `json:"lastLoginDate"` // 最后一次登录时间
	LastLoginIP   string `json:"lastLoginIP"`   // 最后一次登录ip
	Status        int64  `json:"status"`        // 状态1是正常,2是禁用
	Avatar        string `json:"avatar"`        // 头像
	IsAdmin       int64  `json:"isAdmin"`       // 1是超级管理员，2是普通管理员
	CreatedAt     int64  `json:"createdAt"`     // 创建时间
	UpdatedAt     int64  `json:"updatedAt"`     // 更新时间
	CreatedBy     int64  `json:"createdBy"`     // 创建人
	UpdatedBy     int64  `json:"updatedBy"`     // 更新人
}

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
