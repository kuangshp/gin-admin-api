package model

import "database/sql"

type AccountEntity struct {
	BaseEntity
	Username string        `gorm:"username" json:"username"`       // 用户名
	Password string        `gorm:"password" json:"password"`       // 密码
	Status   sql.NullInt64 `gorm:"status,default=1" json:"status"` // 状态1是正常,0是禁用
}

func (t *AccountEntity) TableName() string {
	return "account"
}
