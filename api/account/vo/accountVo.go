package vo

import (
	"gin-admin-api/model"
)

type AccountVo struct {
	ID                int64           `json:"id,string"`
	CreatedAt         model.LocalTime `json:"createdAt"`
	UpdatedAt         model.LocalTime `json:"updatedAt"`
	Username          string          `json:"username"`          // 用户名
	Name              string          `json:"name"`              // 真实姓名
	Mobile            string          `json:"mobile"`            // 手机号码
	Email             string          `json:"email"`             // 邮箱地址
	Avatar            string          `json:"avatar"`            // 用户头像
	IsAdmin           int64           `json:"isAdmin"`           // 是否为超级管理员:0否,1是
	Status            int64           `json:"status"`            // 状态1是正常,0是禁用
	LastLoginIP       string          `json:"lastLoginIp"`       // 最后登录ip地址
	LastLoginDate     model.LocalTime `json:"lastLoginDate"`     // 最后登录时间
	LastLoginCountry  string          `json:"lastLoginCountry"`  // 最后登录国家
	LastLoginProvince string          `json:"lastLoginProvince"` // 最后登录省份
	LastLoginCity     string          `json:"lastLoginCity"`     // 最后登录城市
}

type LoginVo struct {
	AccountVo
	Token string `json:"token"`
}
