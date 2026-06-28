package dto

// CreateSysRoleDTO 创建角色表请求
type CreateSysRoleDTO struct {
	Name string `json:"name" validate:"required,max=50"` // 角色名称
	Description string `json:"description" validate:"omitempty,max=255"` // 描述
	Status int64 `json:"status"` // 状态1是正常,2是禁用
	Sort int64 `json:"sort"` // 排序
	DataScope int64 `json:"dataScope" validate:"required,oneof=1 2 3 4 5"` // 数据范围：1=全部 2=本部门 3=本部门及下级 4=仅本人 5=自定义部门
	CreatedBy int64 `json:"createdBy"` // 创建人
	UpdatedBy int64 `json:"updatedBy"` // 更新人
}

// ModifySysRoleDTO 修改角色表请求
type ModifySysRoleDTO struct {
	ID int64 `json:"id" validate:"required,number,gte=1"` // 主键id
	CreateSysRoleDTO
}
