package mapper

import (
	"time"
	"gin-admin-api/internal/dal/model"
	"gin-admin-api/internal/dal/dto"
	"gin-admin-api/internal/dal/vo"
)

// ISysAccountMapper 后台账号表 mapper 接口
type ISysAccountMapper interface {
	// DtoToEntity 将请求结构体映射到数据库实体
	DtoToEntity(d *dto.CreateSysAccountDTO, lastLoginDate time.Time) *model.SysAccountEntity
	// EntityToVO 将数据库实体映射到响应结构体
	EntityToVO(e *model.SysAccountEntity) *vo.SysAccountVO
}

// sysAccountMapper mapper 实现
type sysAccountMapper struct{}

// NewSysAccountMapper 创建 SysAccountMapper 实例
func NewSysAccountMapper() ISysAccountMapper {
	return &sysAccountMapper{}
}

// DtoToEntity 将 CreateSysAccountDTO 映射到 SysAccountEntity
func (m *sysAccountMapper) DtoToEntity(d *dto.CreateSysAccountDTO, lastLoginDate time.Time) *model.SysAccountEntity {
	e := &model.SysAccountEntity{
		Nickname: d.Nickname, // 昵称
		Username: d.Username, // 登录帐号
		Email: d.Email, // 邮箱
		Mobile: d.Mobile, // 手机号
		Password: d.Password, // 登录密码
		LastLoginDate: lastLoginDate, // 最后一次登录时间
		LastLoginIP: d.LastLoginIP, // 最后一次登录ip
		Status: d.Status, // 状态1是正常,2是禁用
		Avatar: d.Avatar, // 头像
		IsAdmin: d.IsAdmin, // 1是超级管理员，2是普通管理员
		DeptID: d.DeptID, // 所属部门id，用于数据权限计算
	}
	return e
}

// EntityToVO 将 SysAccountEntity 映射到 SysAccountVO
func (m *sysAccountMapper) EntityToVO(e *model.SysAccountEntity) *vo.SysAccountVO {
	if e == nil {
		return nil
	}
	return &vo.SysAccountVO{
		ID: e.ID, // 主键id
		Nickname: e.Nickname, // 昵称
		Username: e.Username, // 登录帐号
		Email: e.Email, // 邮箱
		Mobile: e.Mobile, // 手机号
		Password: e.Password, // 登录密码
        LastLoginDate: e.LastLoginDate.Unix(), // 最后一次登录时间
		LastLoginIP: e.LastLoginIP, // 最后一次登录ip
		Status: e.Status, // 状态1是正常,2是禁用
		Avatar: e.Avatar, // 头像
		IsAdmin: e.IsAdmin, // 1是超级管理员，2是普通管理员
		DeptID: e.DeptID, // 所属部门id，用于数据权限计算
        CreatedAt: e.CreatedAt.Unix(), // 创建时间
        UpdatedAt: e.UpdatedAt.Unix(), // 更新时间
		CreatedBy: e.CreatedBy, // 创建人
		UpdatedBy: e.UpdatedBy, // 更新人

	}
}
