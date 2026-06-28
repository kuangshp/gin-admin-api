package dto

// CreateSysAccountDTO 创建后台账号表请求
type CreateSysAccountDTO struct {
	Password        string `json:"password" validate:"required,min=6,max=16"`                        // 登录密码
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=6,max=16,eqfield=Password"` // 确定密码
	ModifySysAccountDTO
}

// ModifySysAccountDTO 修改后台账号表请求
type ModifySysAccountDTO struct {
	Username   string  `json:"username" validate:"required"`          // 登录帐号
	Email      string  `json:"email"`                                 // 邮箱
	Mobile     string  `json:"mobile"`                                // 手机号
	Status     int64   `json:"status" validate:"omitempty,oneof=1 2"` // 状态1是正常,2是禁用
	Avatar     string  `json:"avatar"`                                // 头像
	DeptID     int64   `json:"deptId" validate:"required,gte=1"`      // 所属部门id，用于数据权限计算
	PostIdList []int64 `json:"postIdList" validate:"required,min=1"`  // 授权的岗位id，第一个为主岗
	RoleIdList []int64 `json:"roleIdList" validate:"required,min=1"`  // 授权的角色id
}

type ResetPasswordDTO struct {
	Id              int64  `json:"id" validate:"required"` // 主键id
	Password        string `json:"password" binding:"required,min=6,max=16"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=6,max=16"`
}

type ModifyCurrentPasswordDTO struct {
	Password        string `json:"password" binding:"required,min=6,max=16"`
	NewPassword     string `json:"newPassword" binding:"required,min=6,max=16"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=6,max=16,eqfield=Password"`
}

type GetSysAccountPageDTO struct {
	Keyword    string `json:"keyword"`    // 用户名,邮箱,手机号码
	Status     int64  `json:"status"`     // 状态1是正常,2是禁用
	DeptID     int64  `json:"deptId"`     // 所属部门id，0表示全部
	PageNumber int64  `json:"pageNumber"` // 当前页
	PageSize   int64  `json:"pageSize"`   // 一页多少条数据
}
