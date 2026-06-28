package vo

// SysRoleResourcesVO 角色和资源中间表视图对象
type SysRoleResourcesVO struct {
	ID int64 `json:"id"` // 主键id
	ResourcesID int64 `json:"resourcesId"` // 关联到sys_resources表主键id
	RoleID int64 `json:"roleId"` // 关联到sys_role表主键id
	CreatedAt int64 `json:"createdAt"` // 创建时间
	UpdatedAt int64 `json:"updatedAt"` // 更新时间
	CreatedBy int64 `json:"createdBy"` // 创建人
	UpdatedBy int64 `json:"updatedBy"` // 更新人
}