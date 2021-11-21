package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseEntity struct {
	// 重写id字段
	ID        int32          `json:"id" gorm:"primaryKey;autoIncrement;comment:主键id"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间;"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"comment:删除时间"`
}
