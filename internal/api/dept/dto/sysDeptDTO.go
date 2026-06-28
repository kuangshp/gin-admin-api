package dto

type CreateSysDeptDTO struct {
	Name     string `json:"name" validate:"required,max=50"`         // 部门名称
	ParentID int64  `json:"parentId" validate:"gte=0"`               // 上级部门id，0表示根部门
	Sort     int64  `json:"sort"`                                    // 排序，值越小越靠前
	Status   int64  `json:"status" validate:"omitempty,oneof=1 2"`   // 状态：1正常，2禁用；不传默认1
	LeaderID int64  `json:"leaderId" validate:"omitempty,gte=1"`     // 负责人账号id，关联sys_account.id
	Phone    string `json:"phone" validate:"omitempty,max=20"`       // 联系电话
	Email    string `json:"email" validate:"omitempty,email,max=60"` // 邮箱
}

type SysDeptPageDTO struct {
	Name       string `json:"name"`       // 部门名称，模糊查询
	Status     int64  `json:"status"`     // 状态：1正常，2禁用；0表示全部
	ParentID   int64  `json:"parentId"`   // 上级部门id，默认0查询根部门
	PageNumber int64  `json:"pageNumber"` // 页码
	PageSize   int64  `json:"pageSize"`   // 每页条数
}
