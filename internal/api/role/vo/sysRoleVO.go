package vo

// SysRoleVO 角色表视图对象
type SysRoleVO struct {
	ID          int64  `json:"id"`          // 主键id
	Name        string `json:"name"`        // 角色名称
	Description string `json:"description"` // 描述
	Status      int64  `json:"status"`      // 状态1是正常,2是禁用
	Sort        int64  `json:"sort"`        // 排序
	CreatedAt   int64  `json:"createdAt"`   // 创建时间
	UpdatedAt   int64  `json:"updatedAt"`   // 更新时间
}
