package vo

// SysAccountRoleVO 账号和角色中间表视图对象
type SysAccountRoleVO struct {
	ID        int64 `json:"id"`        // 主键id
	AccountID int64 `json:"accountId"` // 关联到sys_account表主键id
	RoleID    int64 `json:"roleId"`    // 关联到sys_role表主键id
	CreatedAt int64 `json:"createdAt"` // 创建时间
	UpdatedAt int64 `json:"updatedAt"` // 更新时间
	CreatedBy int64 `json:"createdBy"` // 创建人
	UpdatedBy int64 `json:"updatedBy"` // 更新人
}
