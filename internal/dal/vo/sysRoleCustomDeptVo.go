package vo

// SysRoleCustomDeptVO 角色自定义数据权限部门表视图对象
type SysRoleCustomDeptVO struct {
	ID int64 `json:"id"` // 主键id
	RoleID int64 `json:"roleId"` // 角色id
	DeptID int64 `json:"deptId"` // 可查看部门id
	CreatedAt int64 `json:"createdAt"` // 创建时间
	UpdatedAt int64 `json:"updatedAt"` // 更新时间
}