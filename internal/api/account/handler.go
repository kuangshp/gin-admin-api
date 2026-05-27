package account

import (
	"gin-admin-api/internal/api/base"
	"gin-admin-api/internal/dal/repository"
	"github.com/gin-gonic/gin"
)

type ISysAccount interface {
	CreateSysAccountApi(ctx *gin.Context)                // 创建账号
	DeleteSysAccountByIdApi(ctx *gin.Context)            // 根据id删除账号
	ModifySysAccountByIdApi(ctx *gin.Context)            // 根据id修改账号
	GetSysAccountPageApi(ctx *gin.Context)               // 分页获取账号
	GetSysAccountListApi(ctx *gin.Context)               // 获取账号列表
	GetSysAccountDetailApi(ctx *gin.Context)             // 根据id获取账号详情
	ResetPasswordByIdApi(ctx *gin.Context)               // 根据id重置账号密码
	ModifyCurrentSysAccountPasswordApi(ctx *gin.Context) // 修改当前登录账号密码
}

type SysAccount struct {
	*base.BaseApi
	SysAccountRepository repository.SysAccountRepository
}

func (s SysAccount) CreateSysAccountApi(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (s SysAccount) DeleteSysAccountByIdApi(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (s SysAccount) ModifySysAccountByIdApi(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (s SysAccount) GetSysAccountPageApi(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (s SysAccount) GetSysAccountListApi(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (s SysAccount) GetSysAccountDetailApi(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (s SysAccount) ResetPasswordByIdApi(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (s SysAccount) ModifyCurrentSysAccountPasswordApi(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func NewSysAccount(baseApi *base.BaseApi, accountRepository repository.SysAccountRepository) ISysAccount {
	return &SysAccount{
		BaseApi:              baseApi,
		SysAccountRepository: accountRepository,
	}
}
