package dto

type CreateSysPostDTO struct {
	Name   string `json:"name" validate:"required,max=50"`       // 岗位名称
	Code   string `json:"code" validate:"required,max=60"`       // 岗位编码，全局唯一
	Sort   int64  `json:"sort"`                                  // 排序，值越小越靠前
	Status int64  `json:"status" validate:"omitempty,oneof=1 2"` // 状态：1正常，2禁用；不传默认1
	Remark string `json:"remark" validate:"omitempty,max=255"`   // 备注
}

type ModifySysPostDTO struct {
	Name   string `json:"name" validate:"required,max=50"`       // 岗位名称
	Code   string `json:"code" validate:"required,max=60"`       // 岗位编码，全局唯一
	Sort   int64  `json:"sort"`                                  // 排序，值越小越靠前
	Status int64  `json:"status" validate:"omitempty,oneof=1 2"` // 状态：1正常，2禁用；不传默认1
	Remark string `json:"remark" validate:"omitempty,max=255"`   // 备注
}

type SysPostPageDTO struct {
	Name       string `json:"name"`       // 岗位名称，模糊查询
	Code       string `json:"code"`       // 岗位编码，模糊查询
	Status     int64  `json:"status"`     // 状态：1正常，2禁用；0表示全部
	PageNumber int64  `json:"pageNumber"` // 页码
	PageSize   int64  `json:"pageSize"`   // 每页条数
}
