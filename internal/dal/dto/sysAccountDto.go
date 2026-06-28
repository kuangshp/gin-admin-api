package dto

// CreateSysAccountDTO 创建后台账号表请求
type CreateSysAccountDTO struct {
	Nickname string `json:"nickname" validate:"omitempty,max=60"` // 昵称
	Username string `json:"username" validate:"required,max=60"` // 登录帐号
	Email string `json:"email" validate:"omitempty,email,max=60"` // 邮箱
	Mobile string `json:"mobile" validate:"omitempty,mobile,max=11"` // 手机号
	Password string `json:"password" validate:"required,max=255"` // 登录密码
	LastLoginDate int64 `json:"lastLoginDate"` // 最后一次登录时间
	LastLoginIP string `json:"lastLoginIP" validate:"omitempty,ip,max=30"` // 最后一次登录ip
	Status int64 `json:"status" validate:"required,gte=1"` // 状态1是正常,2是禁用
	Avatar string `json:"avatar" validate:"omitempty,max=500"` // 头像
	IsAdmin int64 `json:"isAdmin" validate:"required,oneof=1 2"` // 1是超级管理员，2是普通管理员
	DeptID int64 `json:"deptId" validate:"required,number,gte=1"` // 所属部门id，用于数据权限计算
	CreatedBy int64 `json:"createdBy"` // 创建人
	UpdatedBy int64 `json:"updatedBy"` // 更新人
}

// ModifySysAccountDTO 修改后台账号表请求
type ModifySysAccountDTO struct {
	ID int64 `json:"id" validate:"required,number,gte=1"` // 主键id
	CreateSysAccountDTO
}
