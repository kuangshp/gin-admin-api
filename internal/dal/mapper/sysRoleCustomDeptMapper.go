package mapper

import (
	"time"
	"gin-admin-api/internal/dal/model"
	"gin-admin-api/internal/dal/dto"
	"gin-admin-api/internal/dal/vo"
)

// ISysRoleCustomDeptMapper 角色自定义数据权限部门表 mapper 接口
type ISysRoleCustomDeptMapper interface {
	// DtoToEntity 将请求结构体映射到数据库实体
	DtoToEntity(d *dto.CreateSysRoleCustomDeptDTO) *model.SysRoleCustomDeptEntity
	// EntityToVO 将数据库实体映射到响应结构体
	EntityToVO(e *model.SysRoleCustomDeptEntity) *vo.SysRoleCustomDeptVO
}

// sysRoleCustomDeptMapper mapper 实现
type sysRoleCustomDeptMapper struct{}

// NewSysRoleCustomDeptMapper 创建 SysRoleCustomDeptMapper 实例
func NewSysRoleCustomDeptMapper() ISysRoleCustomDeptMapper {
	return &sysRoleCustomDeptMapper{}
}

// DtoToEntity 将 CreateSysRoleCustomDeptDTO 映射到 SysRoleCustomDeptEntity
func (m *sysRoleCustomDeptMapper) DtoToEntity(d *dto.CreateSysRoleCustomDeptDTO) *model.SysRoleCustomDeptEntity {
	e := &model.SysRoleCustomDeptEntity{
		RoleID: d.RoleID, // 角色id
		DeptID: d.DeptID, // 可查看部门id
	}
	return e
}

// EntityToVO 将 SysRoleCustomDeptEntity 映射到 SysRoleCustomDeptVO
func (m *sysRoleCustomDeptMapper) EntityToVO(e *model.SysRoleCustomDeptEntity) *vo.SysRoleCustomDeptVO {
	if e == nil {
		return nil
	}
	return &vo.SysRoleCustomDeptVO{
		ID: e.ID, // 主键id
		RoleID: e.RoleID, // 角色id
		DeptID: e.DeptID, // 可查看部门id
        CreatedAt: e.CreatedAt.Unix(), // 创建时间
        UpdatedAt: e.UpdatedAt.Unix(), // 更新时间

	}
}
