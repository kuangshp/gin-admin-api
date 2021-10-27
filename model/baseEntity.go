package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseEntity struct {
	// 重写id字段
	ID        int32          `gorm:"primaryKey;autoIncrement;comment:主键id" json:"id"`
	CreatedAt time.Time      `gorm:"comment:创建时间;" json:"created_at"`
	UpdatedAt time.Time      `gorm:"comment:更新时间" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"comment:删除时间" json:"deleted_at"`
}
