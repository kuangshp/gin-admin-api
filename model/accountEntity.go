package model

import (
	"fmt"
	"time"
)

type AccountEntity struct {
	BaseEntity
	UserName string     `gorm:"type:varchar(50);column:username;not null;unique;comment:账号" json:"username"`
	Password string     `gorm:"type:varchar(200);not null;comment:账号密码" json:"password"`
	Mobile   string     `gorm:"type:varchar(11);default null;unique;comment:手机号码" json:"mobile"`
	NickName string     `gorm:"type:varchar(10);default null;unique;comment:昵称" json:"nick_name"`
	Address  string     `gorm:"type:varchar(100);comment:地址" json:"address"`
	Avatar   string     `gorm:"type:varchar(100);comment:头像" json:"avatar"`
	Desc     string     `gorm:"type:varchar(100);comment:描述" json:"desc"`
	Gender   string     `gorm:"type:varchar(10);comment:性别" json:"gender"`
	BirthDay *time.Time `gorm:"type:datetime;comment:出生年月" json:"birth_day"`
	RoleId   uint32     `gorm:"type:int(11);comment:角色ID" json:"role_id"`
}

// TableName 自定义表名
func (AccountEntity) TableName() string {
	return "account"
}

func init() {
	fmt.Println("执行了=========")
	//global.DB.AutoMigrate(&AccountEntity{})
}
