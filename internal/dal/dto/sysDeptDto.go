package dto

// CreateSysDeptDTO 创建部门表请求
type CreateSysDeptDTO struct {
	Name      string `json:"name" validate:"required,max=50"`           // 部门名称
	ParentID  int64  `json:"parentId" validate:"required,number,gte=1"` // 上级部门id
	FullID    string `json:"fullId" validate:"required,max=255"`        // 全层级ID，例：1,5,12
	FullName  string `json:"fullName" validate:"required,max=255"`      // 全层级名称
	Sort      int64  `json:"sort"`                                      // 排序
	Status    int64  `json:"status"`                                    // 状态1正常 2禁用
	LeaderID  int64  `json:"leaderId" validate:"omitempty,gte=1"`       // 负责人账号id，关联sys_account.id
	Phone     string `json:"phone" validate:"omitempty,e164,max=20"`    // 联系电话
	Email     string `json:"email" validate:"omitempty,email,max=60"`   // 邮箱
	CreatedBy int64  `json:"createdBy"`                                 // 创建人
	UpdatedBy int64  `json:"updatedBy"`                                 // 更新人
}

// ModifySysDeptDTO 修改部门表请求
type ModifySysDeptDTO struct {
	ID int64 `json:"id" validate:"required,number,gte=1"` // 主键id
	CreateSysDeptDTO
}
