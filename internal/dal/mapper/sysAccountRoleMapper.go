package mapper

import (
	"gin-admin-api/internal/dal/dto"
	"gin-admin-api/internal/dal/model"
	"gin-admin-api/internal/dal/vo"
	"time"
)

// ISysAccountRoleMapper 账号和角色中间表 mapper 接口
type ISysAccountRoleMapper interface {
	// DtoToEntity 将请求结构体映射到数据库实体
	DtoToEntity(d *dto.CreateSysAccountRoleDTO) *model.SysAccountRoleEntity
	// EntityToVo 将数据库实体映射到响应结构体
	EntityToVo(e *model.SysAccountRoleEntity) *vo.SysAccountRoleVo
}

// sysAccountRoleMapper mapper 实现
type sysAccountRoleMapper struct{}

// NewSysAccountRoleMapper 创建 SysAccountRoleMapper 实例
func NewSysAccountRoleMapper() ISysAccountRoleMapper {
	return &sysAccountRoleMapper{}
}

// DtoToEntity 将 CreateSysAccountRoleDTO 映射到 SysAccountRoleEntity
func (m *sysAccountRoleMapper) DtoToEntity(d *dto.CreateSysAccountRoleDTO) *model.SysAccountRoleEntity {
	e := &model.SysAccountRoleEntity{
		AccountID: d.AccountID, // 关联到sys_account表主键id
		RoleID:    d.RoleID,    // 关联到sys_role表主键id
	}
	return e
}

// EntityToVo 将 SysAccountRoleEntity 映射到 SysAccountRoleVo
func (m *sysAccountRoleMapper) EntityToVo(e *model.SysAccountRoleEntity) *vo.SysAccountRoleVo {
	if e == nil {
		return nil
	}
	return &vo.SysAccountRoleVo{
		ID:        e.ID,               // 主键id
		AccountID: e.AccountID,        // 关联到sys_account表主键id
		RoleID:    e.RoleID,           // 关联到sys_role表主键id
		CreatedAt: e.CreatedAt.Unix(), // 创建时间
		UpdatedAt: e.UpdatedAt.Unix(), // 更新时间
		CreatedBy: e.CreatedBy,        // 创建人
		UpdatedBy: e.UpdatedBy,        // 更新人

	}
}
