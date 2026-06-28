package datascope

import (
	"fmt"
	"gin-admin-api/pkg/enum"
	"strings"

	"gorm.io/gorm"
)

func Inject(scope *Context) func(db *gorm.DB, tableName string) {
	return func(db *gorm.DB, tableName string) {
		if scope == nil || scope.IsSuperAdmin() || scope.DataScope == enum.DataScopeAll {
			return
		}

		switch strings.ToLower(tableName) {
		case "sys_account", "sys_post":
			injectDeptOrSelf(db, tableName, scope)
		case "sys_dept":
			injectDeptTable(db, tableName, scope)
		default:
			injectSelf(db, tableName, scope)
		}
	}
}

func injectDeptOrSelf(db *gorm.DB, tableName string, scope *Context) {
	if scope.DataScope == enum.DataScopeSelf {
		injectSelf(db, tableName, scope)
		return
	}
	if len(scope.DeptIDs) == 0 {
		db.Where("1 = 0")
		return
	}
	db.Where(fmt.Sprintf("%s.dept_id IN ?", tableName), scope.DeptIDs)
}

func injectDeptTable(db *gorm.DB, tableName string, scope *Context) {
	if scope.DataScope == enum.DataScopeSelf {
		injectSelf(db, tableName, scope)
		return
	}
	if len(scope.DeptIDs) == 0 {
		db.Where("1 = 0")
		return
	}
	db.Where(fmt.Sprintf("%s.id IN ?", tableName), scope.DeptIDs)
}

func injectSelf(db *gorm.DB, tableName string, scope *Context) {
	if scope.AccountID <= 0 {
		db.Where("1 = 0")
		return
	}
	db.Where(fmt.Sprintf("%s.created_by = ?", tableName), scope.AccountID)
}
