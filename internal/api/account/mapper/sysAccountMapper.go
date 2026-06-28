package mapper

import (
	"gin-admin-api/internal/api/account/dto"
	"gin-admin-api/internal/api/account/vo"
	"gin-admin-api/internal/dal/model"
)

// ISysAccountMapper 后台账号表 mapper 接口
type ISysAccountMapper interface {
	// DtoToEntity 将请求结构体映射到数据库实体
	DtoToEntity(d *dto.CreateSysAccountDTO, password string, status int64) *model.SysAccountEntity
	// EntityToVo 将数据库实体映射到响应结构体
	EntityToVo(e *model.SysAccountEntity, deptName string) *vo.SysAccountVO
	EntityListToVo(list []*model.SysAccountEntity, deptNameMap map[int64]string) []*vo.SysAccountVO
}

// sysAccountMapper mapper 实现
type sysAccountMapper struct{}

// NewSysAccountMapper 创建 SysAccountMapper 实例
func NewSysAccountMapper() ISysAccountMapper {
	return &sysAccountMapper{}
}

// DtoToEntity 将 CreateSysAccountDTO 映射到 SysAccountEntity
func (m *sysAccountMapper) DtoToEntity(d *dto.CreateSysAccountDTO, password string, status int64) *model.SysAccountEntity {
	e := &model.SysAccountEntity{
		Username: d.Username, // 登录帐号
		Email:    d.Email,    // 邮箱
		Mobile:   d.Mobile,   // 手机号
		Password: password,   // 登录密码
		Status:   status,     // 状态1是正常,2是禁用
		Avatar:   d.Avatar,   // 头像
		DeptID:   d.DeptID,   // 所属部门id，用于数据权限计算
	}
	return e
}

// EntityToVo 将 SysAccountEntity 映射到 SysAccountVo
func (m *sysAccountMapper) EntityToVo(e *model.SysAccountEntity, deptName string) *vo.SysAccountVO {
	if e == nil {
		return nil
	}
	return &vo.SysAccountVO{
		ID:            e.ID,                   // 主键id
		Username:      e.Username,             // 登录帐号
		Email:         e.Email,                // 邮箱
		Mobile:        e.Mobile,               // 手机号
		Password:      e.Password,             // 登录密码
		LastLoginDate: e.LastLoginDate.Unix(), // 最后一次登录时间
		LastLoginIP:   e.LastLoginIP,          // 最后一次登录ip
		Status:        e.Status,               // 状态1是正常,2是禁用
		Avatar:        e.Avatar,               // 头像
		IsAdmin:       e.IsAdmin,              // 1是超级管理员，2是普通管理员
		DeptID:        e.DeptID,               // 所属部门id
		DeptName:      deptName,               // 所属部门
		CreatedAt:     e.CreatedAt.Unix(),     // 创建时间
		UpdatedAt:     e.UpdatedAt.Unix(),     // 更新时间
		CreatedBy:     e.CreatedBy,            // 创建人
		UpdatedBy:     e.UpdatedBy,            // 更新人
	}
}

func (m *sysAccountMapper) EntityListToVo(list []*model.SysAccountEntity, deptNameMap map[int64]string) []*vo.SysAccountVO {
	result := make([]*vo.SysAccountVO, 0, len(list))
	for _, item := range list {
		result = append(result, m.EntityToVo(item, deptNameMap[item.DeptID]))
	}
	return result
}
