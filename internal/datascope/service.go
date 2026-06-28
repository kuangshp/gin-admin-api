package datascope

import (
	"context"
	"errors"
	"gin-admin-api/internal/dal/model"
	"gin-admin-api/pkg/enum"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type roleScopeRow struct {
	RoleID    int64 `gorm:"column:role_id"`
	DataScope int64 `gorm:"column:data_scope"`
	DeptID    int64 `gorm:"column:dept_id"`
}

func (s *Service) BuildContext(ctx context.Context, accountID int64) (*Context, error) {
	if accountID <= 0 {
		return nil, errors.New("账号ID为空")
	}

	var account model.SysAccountEntity
	if err := s.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL AND status = ?", accountID, enum.StatusNormalEnum).
		First(&account).Error; err != nil {
		return nil, err
	}

	scopeCtx := &Context{
		AccountID: account.ID,
		IsAdmin:   account.IsAdmin,
		DeptID:    account.DeptID,
		DataScope: enum.DataScopeSelf,
	}
	if scopeCtx.IsSuperAdmin() {
		scopeCtx.DataScope = enum.DataScopeAll
		return scopeCtx, nil
	}

	var roleIDs []int64
	if err := s.db.WithContext(ctx).
		Model(&model.SysAccountRoleEntity{}).
		Where("account_id = ? AND deleted_at IS NULL", account.ID).
		Pluck("role_id", &roleIDs).Error; err != nil {
		return nil, err
	}
	scopeCtx.RoleIDs = distinctPositive(roleIDs)
	if len(scopeCtx.RoleIDs) == 0 {
		scopeCtx.DeptIDs = []int64{}
		return scopeCtx, nil
	}

	var rows []roleScopeRow
	if err := s.db.WithContext(ctx).
		Table("sys_role r").
		Select("r.id AS role_id, r.data_scope, rcd.dept_id").
		Joins("LEFT JOIN sys_role_custom_dept rcd ON rcd.role_id = r.id AND rcd.deleted_at IS NULL").
		Where("r.id IN ? AND r.deleted_at IS NULL AND r.status = ?", scopeCtx.RoleIDs, enum.StatusNormalEnum).
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	return s.mergeRoleScopes(ctx, scopeCtx, rows)
}

func (s *Service) mergeRoleScopes(ctx context.Context, scopeCtx *Context, rows []roleScopeRow) (*Context, error) {
	hasAll := false
	hasDeptAndChild := false
	hasDept := false
	hasCustom := false
	deptIDs := make([]int64, 0)
	hasValidRole := false

	for _, row := range rows {
		hasValidRole = true
		switch row.DataScope {
		case enum.DataScopeAll:
			hasAll = true
		case enum.DataScopeDeptAndChild:
			hasDeptAndChild = true
		case enum.DataScopeDept:
			hasDept = true
		case enum.DataScopeCustomDept:
			hasCustom = true
			if row.DeptID > 0 {
				deptIDs = append(deptIDs, row.DeptID)
			}
		}
	}

	if !hasValidRole {
		scopeCtx.DeptIDs = []int64{}
		return scopeCtx, nil
	}
	if hasAll {
		scopeCtx.DataScope = enum.DataScopeAll
		scopeCtx.DeptIDs = nil
		return scopeCtx, nil
	}
	if hasDeptAndChild && scopeCtx.DeptID > 0 {
		childIDs, err := s.FindChildDeptIDs(ctx, scopeCtx.DeptID)
		if err != nil {
			return nil, err
		}
		deptIDs = append(deptIDs, childIDs...)
	}
	if hasDept && scopeCtx.DeptID > 0 {
		deptIDs = append(deptIDs, scopeCtx.DeptID)
	}

	if hasDeptAndChild || hasDept || hasCustom {
		scopeCtx.DataScope = enum.DataScopeCustomDept
		scopeCtx.DeptIDs = distinctPositive(deptIDs)
		return scopeCtx, nil
	}

	scopeCtx.DataScope = enum.DataScopeSelf
	scopeCtx.DeptIDs = nil
	return scopeCtx, nil
}

func (s *Service) FindChildDeptIDs(ctx context.Context, deptID int64) ([]int64, error) {
	if deptID <= 0 {
		return []int64{}, nil
	}
	var ids []int64
	err := s.db.WithContext(ctx).
		Model(&model.SysDeptEntity{}).
		Where("deleted_at IS NULL AND status = ? AND (id = ? OR FIND_IN_SET(?, full_id))", enum.StatusNormalEnum, deptID, deptID).
		Pluck("id", &ids).Error
	if err != nil {
		return nil, err
	}
	return distinctPositive(ids), nil
}

func distinctPositive(values []int64) []int64 {
	seen := make(map[int64]struct{}, len(values))
	result := make([]int64, 0, len(values))
	for _, value := range values {
		if value <= 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
