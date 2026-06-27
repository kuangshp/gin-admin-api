package dto

// CreateSysPostDTO 创建岗位表请求
type CreateSysPostDTO struct {
	Name string `json:"name" validate:"required,max=50"` // 岗位名称
	Code string `json:"code" validate:"required,max=60"` // 岗位编码
	DeptID int64 `json:"deptId" validate:"required,number,gte=1"` // 所属部门id
	Sort int64 `json:"sort"` // 排序
	Status int64 `json:"status"` // 状态1正常 2禁用
	Remark string `json:"remark" validate:"omitempty,max=255"` // 备注
	CreatedBy int64 `json:"createdBy"` // 创建人
	UpdatedBy int64 `json:"updatedBy"` // 更新人
}

// ModifySysPostDTO 修改岗位表请求
type ModifySysPostDTO struct {
	ID int64 `json:"id" validate:"required,number,gte=1"` // 主键id
	CreateSysPostDTO
}
