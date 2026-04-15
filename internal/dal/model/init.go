package model

import (
	"gin-admin-api/internal/dal/model/entity"
)

// GetAllModels 返回所有需要 AutoMigrate 的 model 列表
// 新增 model 后在此追加即可
func GetAllModels() []interface{} {
	return []interface{}{
		&entity.AccountEntity{},
	}
}
