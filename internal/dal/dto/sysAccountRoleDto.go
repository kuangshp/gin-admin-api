package dto

// CreateSysAccountRoleDTO 创建账号和角色中间表请求
type CreateSysAccountRoleDTO struct {
	AccountID int64 `json:"accountId" validate:"required,number,gte=1"` // 关联到sys_account表主键id
	RoleID int64 `json:"roleId" validate:"required,number,gte=1"` // 关联到sys_role表主键id
	CreatedBy int64 `json:"createdBy"` // 创建人
	UpdatedBy int64 `json:"updatedBy"` // 更新人
}

// ModifySysAccountRoleDTO 修改账号和角色中间表请求
type ModifySysAccountRoleDTO struct {
	ID int64 `json:"id" validate:"required,number,gte=1"` // 主键id
	CreateSysAccountRoleDTO
}
