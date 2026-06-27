package vo

// SysPostVO 岗位表视图对象
type SysPostVO struct {
	ID int64 `json:"id"` // 主键id
	Name string `json:"name"` // 岗位名称
	Code string `json:"code"` // 岗位编码
	DeptID int64 `json:"deptId"` // 所属部门id
	Sort int64 `json:"sort"` // 排序
	Status int64 `json:"status"` // 状态1正常 2禁用
	Remark string `json:"remark"` // 备注
	CreatedAt int64 `json:"createdAt"` // 创建时间
	UpdatedAt int64 `json:"updatedAt"` // 更新时间
	CreatedBy int64 `json:"createdBy"` // 创建人
	UpdatedBy int64 `json:"updatedBy"` // 更新人
}