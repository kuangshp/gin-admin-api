package vo

// SysDeptVO 部门表视图对象
type SysDeptVO struct {
	ID int64 `json:"id"` // 主键id
	Name string `json:"name"` // 部门名称
	ParentID int64 `json:"parentId"` // 上级部门id
	FullID string `json:"fullId"` // 全层级ID，例：1,5,12
	FullName string `json:"fullName"` // 全层级名称
	Sort int64 `json:"sort"` // 排序
	Status int64 `json:"status"` // 状态1正常 2禁用
	LeaderID int64 `json:"leaderId"` // 负责人账号id，关联sys_account.id
	Phone string `json:"phone"` // 联系电话
	Email string `json:"email"` // 邮箱
	CreatedAt int64 `json:"createdAt"` // 创建时间
	UpdatedAt int64 `json:"updatedAt"` // 更新时间
	CreatedBy int64 `json:"createdBy"` // 创建人
	UpdatedBy int64 `json:"updatedBy"` // 更新人
}