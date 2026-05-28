package mapper

import (
	"gin-admin-api/internal/api/auth/vo"
	"gin-admin-api/internal/dal/model"
	"time"
)

type IAuthMapper interface {
	LoginToVo(e *model.SysAccountEntity, lastLoginIp, token string) *vo.AccountLoginVO
}

type authMapper struct{}

func NewAuthMapper() IAuthMapper {
	return &authMapper{}
}

func (m *authMapper) LoginToVo(e *model.SysAccountEntity, lastLoginIp, token string) *vo.AccountLoginVO {
	return &vo.AccountLoginVO{
		ID:            e.ID,              // 主键id
		Username:      e.Username,        // 登录帐号
		Email:         e.Email,           // 邮箱
		Mobile:        e.Mobile,          // 手机号
		LastLoginDate: time.Now().Unix(), // 最后一次登录时间
		LastLoginIP:   lastLoginIp,       // 最后一次登录ip
		IsAdmin:       e.IsAdmin,         // 1是超级管理员，2是普通管理员
		Token:         token,             // token
	}
}
