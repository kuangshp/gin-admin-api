package dto

// CreateSysRoleDTO 创建角色表请求
type CreateSysRoleDTO struct {
	Name        string `json:"name" validate:"required"`    // 角色名称
	Description string `json:"description"`                 // 描述
	Status      int64  `json:"status" validate:"oneof=1 2"` // 状态1是正常,2是禁用
	Sort        int64  `json:"sort"`                        // 排序
	CreatedBy   int64  `json:"createdBy"`                   // 创建人
	UpdatedBy   int64  `json:"updatedBy"`                   // 更新人
}

// ModifySysRoleDTO 修改角色表请求
type ModifySysRoleDTO struct {
	ID int64 `json:"id" validate:"required,number,gte=1"` // 主键id
	CreateSysRoleDTO
}
