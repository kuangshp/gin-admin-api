package dto

// CreateSysRoleResourcesDTO 创建角色和资源中间表请求
type CreateSysRoleResourcesDTO struct {
	ResourcesID int64 `json:"resourcesId" validate:"required,number,gte=1"` // 关联到sys_resources表主键id
	RoleID int64 `json:"roleId" validate:"required,number,gte=1"` // 关联到sys_role表主键id
	CreatedBy int64 `json:"createdBy"` // 创建人
	UpdatedBy int64 `json:"updatedBy"` // 更新人
}

// ModifySysRoleResourcesDTO 修改角色和资源中间表请求
type ModifySysRoleResourcesDTO struct {
	ID int64 `json:"id" validate:"required,number,gte=1"` // 主键id
	CreateSysRoleResourcesDTO
}
