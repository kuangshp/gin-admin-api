package dto

// CreateSysResourcesDTO 创建资源表请求
type CreateSysResourcesDTO struct {
	Title         string `json:"title" validate:"required"`   // 名称:按钮标题,或菜单标题
	URL           string `json:"uRL" validate:"required"`     // 按钮请求url,或菜单路由
	Method        string `json:"method"`                      // 接口的请求方式
	Icon          string `json:"icon"`                        // 菜单小图标
	ResourcesType int64  `json:"resourcesType"`               // 类型:1表示目录,2表示菜单,3表示接口
	IsCache       int64  `json:"isCache"`                     // 是否缓存:1表示缓存:2不缓存
	IsHidden      int64  `json:"isHidden"`                    // 是否隐藏:1表示不隐藏,2表示隐藏
	IsLink        int64  `json:"isLink"`                      // 是否为外部链接:1表示不是,2表示是
	ParentID      int64  `json:"parentId"`                    // 上一级id，0=顶级
	Sort          int64  `json:"sort"`                        // 菜单,或按钮排序
	Status        int64  `json:"status" validate:"oneof=1 2"` // 状态1是正常,2是禁用
	Description   string `json:"description"`                 // 描述
}

// ModifySysResourcesDTO 修改资源表请求
type ModifySysResourcesDTO struct {
	ID int64 `json:"id" validate:"required,number,gte=1"` // 主键id
	CreateSysResourcesDTO
}

type PageSysResourcesDTO struct {
	PageNumber int64  `json:"pageNumber"`
	PageSize   int64  `json:"pageSize"`
	ParentID   int64  `json:"parentId,optional,default=0"` // 上一级id
	Status     int64  `json:"status,optional"`             // 状态0是正常,1是禁用
	Title      string `json:"title,optional"`              // 名称:按钮标题,或菜单标题
}

type CatalogDTO struct {
	CatalogType int64 `form:"catalogType,optional"` //类型: 1的时候只查询出模块,2的时候查询出模块和菜单3,的时候查询模块、菜单、按钮
}
