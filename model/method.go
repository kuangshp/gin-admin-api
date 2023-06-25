package model

import "gorm.io/gen"

// Method 自定义sql查询
type Method interface {
	SimpleFindByNameAndAge(username string, status int64) (gen.T, error)
}
