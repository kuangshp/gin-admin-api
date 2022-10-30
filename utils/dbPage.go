package utils

import (
	"gorm.io/gorm"
	"strconv"
)

// Paginate 定义分页函数
func Paginate(pageNumber, pageSize string) func(db *gorm.DB) *gorm.DB {
	pageSizeInt, err1 := strconv.Atoi(pageSize)
	if err1 != nil {
		pageSizeInt = 10
	}
	pageNumberInt, err2 := strconv.Atoi(pageNumber)
	if err2 != nil {
		pageNumberInt = 1
	}
	if pageNumberInt == 0 {
		pageNumberInt = 1
	}
	if pageSizeInt <= 0 {
		pageSizeInt = 10
	}
	return func(db *gorm.DB) *gorm.DB {
		offset := (pageNumberInt - 1) * pageSizeInt
		return db.Offset(offset).Limit(pageSizeInt)
	}
}
