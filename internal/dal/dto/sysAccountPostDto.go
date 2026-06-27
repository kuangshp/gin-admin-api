package dto

// CreateSysAccountPostDTO 创建账号岗位关联表请求
type CreateSysAccountPostDTO struct {
	AccountID int64 `json:"accountId" validate:"required,number,gte=1"` // 账号id
	PostID int64 `json:"postId" validate:"required,number,gte=1"` // 岗位id
	CreatedBy int64 `json:"createdBy"` // 创建人
	UpdatedBy int64 `json:"updatedBy"` // 更新人
}

// ModifySysAccountPostDTO 修改账号岗位关联表请求
type ModifySysAccountPostDTO struct {
	ID int64 `json:"id" validate:"required,number,gte=1"` // 主键id
	CreateSysAccountPostDTO
}
