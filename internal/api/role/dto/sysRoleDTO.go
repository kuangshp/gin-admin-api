package dto

// CreateSysRoleDTO 创建角色表请求
type CreateSysRoleDTO struct {
	Name            string  `json:"name" validate:"required"`                  // 角色名称
	Description     string  `json:"description"`                               // 描述
	Status          int64   `json:"status" validate:"oneof=1 2"`               // 状态1是正常,2是禁用
	Sort            int64   `json:"sort"`                                      // 排序
	ResourcesIdList []int64 `json:"resourcesIdList" validate:"required,min=1"` // 资源id
}

type RolePageDTO struct {
	Name       string `json:"name"`   // 角色名称
	Status     int64  `json:"status"` // 状态1是正常,2是禁用
	PageNumber int64  `json:"pageNumber"`
	PageSize   int64  `json:"pageSize"`
}
