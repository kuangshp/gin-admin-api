package model

import (
	"gorm.io/gorm"
	"time"
)

type Account struct {
	gorm.Model // 继承这个主要是因为要软删除字段
	// 重写那几个字段
	ID        uint      `gorm:"primarykey;autoIncrement;comment:主键ID" json:"id"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"更新时间" json:"updated_at"`
	UserName  string    `gorm:"type:varchar(50);column(username);not null;unique;comment:账号" json:"username"`
	Password  string    `gorm:"type:varchar(200);not null;comment:账号密码" json:"password"`
	Mobile    string    `gorm:"varchar(11);default null;comment:手机号码" json:"mobile"`
}

// 自定义表名
func (Account) TableName() string {
	return "account"
}
