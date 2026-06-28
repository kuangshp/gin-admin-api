package vo

// SysRoleVO 角色表视图对象
type SysRoleVO struct {
	ID int64 `json:"id"` // 主键id
	Name string `json:"name"` // 角色名称
	Description string `json:"description"` // 描述
	Status int64 `json:"status"` // 状态1是正常,2是禁用
	Sort int64 `json:"sort"` // 排序
	DataScope int64 `json:"dataScope"` // 数据范围：1=全部 2=本部门 3=本部门及下级 4=仅本人 5=自定义部门
	CreatedAt int64 `json:"createdAt"` // 创建时间
	UpdatedAt int64 `json:"updatedAt"` // 更新时间
	CreatedBy int64 `json:"createdBy"` // 创建人
	UpdatedBy int64 `json:"updatedBy"` // 更新人
}