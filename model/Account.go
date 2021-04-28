package model

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	UserName string `gorm:"type:varchar(50);column(username);not null;unique;comment:账号"`
	Password string `gorm:"type:varchar(200);not null;comment:账号密码"`
	Mobile   string `gorm:"varchar(11);not null;unique;comment:手机号码"`
}

// 自定义表名
func (Account) TableName() string {
	return "account"
}
