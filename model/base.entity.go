package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseEntity struct {
	gorm.Model // 继承这个主要是因为要软删除字段
	// 重写那几个字段
	ID        uint      `gorm:"primarykey;autoIncrement;comment:主键ID" json:"id"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"更新时间" json:"updated_at"`
}
