package dto

// CreateSysRoleCustomDeptDTO 创建角色自定义数据权限部门表请求
type CreateSysRoleCustomDeptDTO struct {
	RoleID int64 `json:"roleId" validate:"required,number,gte=1"` // 角色id
	DeptID int64 `json:"deptId" validate:"required,number,gte=1"` // 可查看部门id
}

// ModifySysRoleCustomDeptDTO 修改角色自定义数据权限部门表请求
type ModifySysRoleCustomDeptDTO struct {
	ID int64 `json:"id" validate:"required,number,gte=1"` // 主键id
	CreateSysRoleCustomDeptDTO
}
