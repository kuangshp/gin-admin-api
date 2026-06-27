package mapper

import (
	"time"
	"gin-admin-api/internal/dal/model"
	"gin-admin-api/internal/dal/dto"
	"gin-admin-api/internal/dal/vo"
)

// ISysRoleMapper 角色表 mapper 接口
type ISysRoleMapper interface {
	// DtoToEntity 将请求结构体映射到数据库实体
	DtoToEntity(d *dto.CreateSysRoleDTO) *model.SysRoleEntity
	// EntityToVO 将数据库实体映射到响应结构体
	EntityToVO(e *model.SysRoleEntity) *vo.SysRoleVO
}

// sysRoleMapper mapper 实现
type sysRoleMapper struct{}

// NewSysRoleMapper 创建 SysRoleMapper 实例
func NewSysRoleMapper() ISysRoleMapper {
	return &sysRoleMapper{}
}

// DtoToEntity 将 CreateSysRoleDTO 映射到 SysRoleEntity
func (m *sysRoleMapper) DtoToEntity(d *dto.CreateSysRoleDTO) *model.SysRoleEntity {
	e := &model.SysRoleEntity{
		Name: d.Name, // 角色名称
		Description: d.Description, // 描述
		Status: d.Status, // 状态1是正常,2是禁用
		Sort: d.Sort, // 排序
	}
	return e
}

// EntityToVO 将 SysRoleEntity 映射到 SysRoleVO
func (m *sysRoleMapper) EntityToVO(e *model.SysRoleEntity) *vo.SysRoleVO {
	if e == nil {
		return nil
	}
	return &vo.SysRoleVO{
		ID: e.ID, // 主键id
		Name: e.Name, // 角色名称
		Description: e.Description, // 描述
		Status: e.Status, // 状态1是正常,2是禁用
		Sort: e.Sort, // 排序
        CreatedAt: e.CreatedAt.Unix(), // 创建时间
        UpdatedAt: e.UpdatedAt.Unix(), // 更新时间
		CreatedBy: e.CreatedBy, // 创建人
		UpdatedBy: e.UpdatedBy, // 更新人

	}
}
