package mapper

import (
	"time"
	"gin-admin-api/internal/dal/model"
	"gin-admin-api/internal/dal/dto"
	"gin-admin-api/internal/dal/vo"
)

// ISysRoleResourcesMapper 角色和资源中间表 mapper 接口
type ISysRoleResourcesMapper interface {
	// DtoToEntity 将请求结构体映射到数据库实体
	DtoToEntity(d *dto.CreateSysRoleResourcesDTO) *model.SysRoleResourcesEntity
	// EntityToVO 将数据库实体映射到响应结构体
	EntityToVO(e *model.SysRoleResourcesEntity) *vo.SysRoleResourcesVO
}

// sysRoleResourcesMapper mapper 实现
type sysRoleResourcesMapper struct{}

// NewSysRoleResourcesMapper 创建 SysRoleResourcesMapper 实例
func NewSysRoleResourcesMapper() ISysRoleResourcesMapper {
	return &sysRoleResourcesMapper{}
}

// DtoToEntity 将 CreateSysRoleResourcesDTO 映射到 SysRoleResourcesEntity
func (m *sysRoleResourcesMapper) DtoToEntity(d *dto.CreateSysRoleResourcesDTO) *model.SysRoleResourcesEntity {
	e := &model.SysRoleResourcesEntity{
		ResourcesID: d.ResourcesID, // 关联到sys_resources表主键id
		RoleID: d.RoleID, // 关联到sys_role表主键id
	}
	return e
}

// EntityToVO 将 SysRoleResourcesEntity 映射到 SysRoleResourcesVO
func (m *sysRoleResourcesMapper) EntityToVO(e *model.SysRoleResourcesEntity) *vo.SysRoleResourcesVO {
	if e == nil {
		return nil
	}
	return &vo.SysRoleResourcesVO{
		ID: e.ID, // 主键id
		ResourcesID: e.ResourcesID, // 关联到sys_resources表主键id
		RoleID: e.RoleID, // 关联到sys_role表主键id
        CreatedAt: e.CreatedAt.Unix(), // 创建时间
        UpdatedAt: e.UpdatedAt.Unix(), // 更新时间
		CreatedBy: e.CreatedBy, // 创建人
		UpdatedBy: e.UpdatedBy, // 更新人

	}
}
