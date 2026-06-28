package dto

// CreateSysRoleDTO 创建角色表请求
type CreateSysRoleDTO struct {
	Name            string  `json:"name" validate:"required"`                      // 角色名称
	Description     string  `json:"description"`                                   // 描述
	Status          int64   `json:"status" validate:"oneof=1 2"`                   // 状态1是正常,2是禁用
	Sort            int64   `json:"sort"`                                          // 排序
	ResourcesIdList []int64 `json:"resourcesIdList" validate:"required,min=1"`     // 资源id
	DataScope       int64   `json:"dataScope" validate:"required,oneof=1 2 3 4 5"` // 数据范围：1=全部 2=本部门 3=本部门及下级 4=仅本人 5=自定义部门
	DeptIdList      []int64 `json:"deptIdList"`                                    // 自定义部门时候传递部门id
}

type RolePageDTO struct {
	Name       string `json:"name"`       // 角色名称
	Status     int64  `json:"status"`     // 状态1是正常,2是禁用
	PageNumber int64  `json:"pageNumber"` // 当前页
	PageSize   int64  `json:"pageSize"`   // 一页多少条数据
}
