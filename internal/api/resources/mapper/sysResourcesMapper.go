package mapper

import (
	"gin-admin-api/internal/api/resources/dto"
	"gin-admin-api/internal/api/resources/vo"
	"gin-admin-api/internal/dal/model"
)

// ISysResourcesMapper 资源表 mapper 接口
type ISysResourcesMapper interface {
	// DtoToEntity 将请求结构体映射到数据库实体
	DtoToEntity(d *dto.CreateSysResourcesDTO) *model.SysResourcesEntity
	// EntityToVo 将数据库实体映射到响应结构体
	EntityToVo(e *model.SysResourcesEntity) *vo.SysResourcesVO
}

// sysResourcesMapper mapper 实现
type sysResourcesMapper struct{}

// NewSysResourcesMapper 创建 SysResourcesMapper 实例
func NewSysResourcesMapper() ISysResourcesMapper {
	return &sysResourcesMapper{}
}

// DtoToEntity 将 CreateSysResourcesDTO 映射到 SysResourcesEntity
func (m *sysResourcesMapper) DtoToEntity(d *dto.CreateSysResourcesDTO) *model.SysResourcesEntity {
	e := &model.SysResourcesEntity{
		Title:         d.Title,         // 名称:按钮标题,或菜单标题
		URL:           d.URL,           // 按钮请求url,或菜单路由
		Method:        d.Method,        // 接口的请求方式
		Icon:          d.Icon,          // 菜单小图标
		ResourcesType: d.ResourcesType, // 类型:1表示目录,2表示菜单,3表示接口
		IsCache:       d.IsCache,       // 是否缓存:1表示缓存:2不缓存
		IsHidden:      d.IsHidden,      // 是否隐藏:1表示不隐藏,2表示隐藏
		IsLink:        d.IsLink,        // 是否为外部链接:1表示不是,2表示是
		ParentID:      d.ParentID,      // 上一级id，0=顶级
		Sort:          d.Sort,          // 菜单,或按钮排序
		Status:        d.Status,        // 状态1是正常,2是禁用
		Description:   d.Description,   // 描述
	}
	return e
}

// EntityToVo 将 SysResourcesEntity 映射到 SysResourcesVO
func (m *sysResourcesMapper) EntityToVo(e *model.SysResourcesEntity) *vo.SysResourcesVO {
	if e == nil {
		return nil
	}
	return &vo.SysResourcesVO{
		ID:            e.ID,               // 主键id
		Title:         e.Title,            // 名称:按钮标题,或菜单标题
		URL:           e.URL,              // 按钮请求url,或菜单路由
		Method:        e.Method,           // 接口的请求方式
		Icon:          e.Icon,             // 菜单小图标
		ResourcesType: e.ResourcesType,    // 类型:1表示目录,2表示菜单,3表示接口
		IsCache:       e.IsCache,          // 是否缓存:1表示缓存:2不缓存
		IsHidden:      e.IsHidden,         // 是否隐藏:1表示不隐藏,2表示隐藏
		IsLink:        e.IsLink,           // 是否为外部链接:1表示不是,2表示是
		ParentID:      e.ParentID,         // 上一级id，0=顶级
		Sort:          e.Sort,             // 菜单,或按钮排序
		Status:        e.Status,           // 状态1是正常,2是禁用
		Description:   e.Description,      // 描述
		IsAdminHave:   e.IsAdminHave,      // 是否超管独有,1表示是,0表示不是
		CreatedAt:     e.CreatedAt.Unix(), // 创建时间
		UpdatedAt:     e.UpdatedAt.Unix(), // 更新时间
	}
}
