package vo

// SysAccountPostVO 账号岗位关联表视图对象
type SysAccountPostVO struct {
	ID int64 `json:"id"` // 主键id
	AccountID int64 `json:"accountId"` // 账号id
	PostID int64 `json:"postId"` // 岗位id
	CreatedAt int64 `json:"createdAt"` // 创建时间
	UpdatedAt int64 `json:"updatedAt"` // 更新时间
	CreatedBy int64 `json:"createdBy"` // 创建人
	UpdatedBy int64 `json:"updatedBy"` // 更新人
}