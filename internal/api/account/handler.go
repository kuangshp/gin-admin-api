package account

import (
	"errors"
	"fmt"
	"gin-admin-api/internal/api/account/dto"
	"gin-admin-api/internal/api/account/vo"
	"gin-admin-api/internal/api/base"
	"gin-admin-api/internal/dal/repository"
	"gin-admin-api/pkg/enum"
	"gin-admin-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/kuangshp/go-utils/k"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type IAccount interface {
	CreateAccountApi(ctx *gin.Context)                // 用户注册
	LoginAccountApi(ctx *gin.Context)                 // 用户名和密码登录
	DeleteAccountByIdApi(ctx *gin.Context)            // 根据id删除数据
	ModifyPasswordByIdApi(ctx *gin.Context)           // 根据id修改密码
	UpdateCurrentAccountPasswordApi(ctx *gin.Context) // 修改当前账号密码
	UpdateStatusByIdApi(ctx *gin.Context)             // 根据id修改状态
	GetAccountByIdApi(ctx *gin.Context)               //根据id查询数据
	GetAccountPageApi(ctx *gin.Context)               // 分页获取数据
}

type Account struct {
	*base.BaseApi
	AccountRepository repository.IAccountRepository
}

// CreateAccountApi 创建账号
func (a *Account) CreateAccountApi(ctx *gin.Context) {
	var createAccountDto dto.CreateAccountDto
	if err := a.BindAndValidateJSON(ctx, &createAccountDto); err != nil {
		return
	}
	// 1.判断两次密码是否一致
	if createAccountDto.Password != createAccountDto.ConfirmPassword {
		a.Fail(ctx, "两次密码不一致", nil)
		return
	}
	// 2.判断账号是否已经存在
	exists, err := a.AccountRepository.ExistsByUsername(ctx, createAccountDto.Username)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		a.Fail(ctx, "创建账号失败", err)
		return
	}
	if exists {
		a.Fail(ctx, fmt.Sprintf("%s已经存在,不能重复创建", createAccountDto.Username), nil)
		return
	}
	// 3.对密码加密
	password, err := k.MakePassword(createAccountDto.Password)
	if err != nil {
		a.Fail(ctx, "创建账号失败", err)
		return
	}
	//ctx1 := ctx.Request.Context()
	if err = a.AccountRepository.Create(a.Ctx(ctx), createAccountDto.Username, password); err != nil {
		a.Fail(ctx, "创建失败", err)
		return
	}
	a.Success(ctx, "创建成功")
}

// LoginAccountApi 用户名和密码登录
func (a *Account) LoginAccountApi(ctx *gin.Context) {
	var accountDto dto.AccountDto
	if err := a.BindAndValidateJSON(ctx, &accountDto); err != nil {
		return
	}
	// 根据账号查询数据
	result, err := a.AccountRepository.GetByUsername(ctx, accountDto.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			a.Fail(ctx, "账号或密码错误", nil)
			return
		}
		a.Fail(ctx, "账号或密码错误", err)
		return
	}

	if result.Status == enum.StatusForbidEnum {
		a.Fail(ctx, "当前账号已经被禁用,请联系管理员", nil)
		return
	}
	isValid := k.CheckPassword(result.Password, accountDto.Password)
	if !isValid {
		a.Fail(ctx, "账号或密码错误", nil)
		return
	}
	// 3.生产token返回给前端
	hmacUser := utils.HmacUser{
		AccountId: result.ID,
		Username:  result.Username,
	}
	token, err := utils.GenerateToken(hmacUser)
	if err != nil {
		a.Fail(ctx, "账号或密码错误", err)
		return
	}
	// 更新账号登录信息
	if err = a.AccountRepository.UpdateLoginInfo(ctx, result.ID, ctx.ClientIP()); err != nil {
		a.Fail(ctx, "账号或密码错误", err)
		return
	}
	a.Success(ctx, vo.LoginVo{
		AccountVo: vo.AccountVo{
			ID:            result.ID,
			CreatedAt:     result.CreatedAt,
			UpdatedAt:     result.UpdatedAt,
			Username:      result.Username,
			Name:          result.Name,
			Mobile:        result.Mobile,
			Email:         result.Email,
			Avatar:        result.Avatar,
			IsAdmin:       result.IsAdmin,
			Status:        result.Status,
			LastLoginDate: time.Now(),
			LastLoginIP:   result.LastLoginIP,
		},
		Token: token,
	})
}

// DeleteAccountByIdApi 根据id删除数据
func (a *Account) DeleteAccountByIdApi(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)
	// 1.判断超级管理员不能删除
	accountData, err := a.AccountRepository.GetByID(ctx, idInt)
	if err != nil {
		a.Fail(ctx, "删除失败", err)
		return
	}
	if accountData.IsAdmin == enum.AdminAccount {
		a.Fail(ctx, "超级管理员不能被删除", nil)
		return
	}
	// 2.判断不能自己删除自己
	if accountData.ID == idInt {
		a.Fail(ctx, "自己不能删除自己", nil)
		return
	}

	if err = a.AccountRepository.Delete(ctx, idInt); err != nil {
		a.Fail(ctx, "删除失败", err)
		return
	}
	a.Success(ctx, "删除成功")
}

// ModifyPasswordByIdApi 根据id修改密码
func (a *Account) ModifyPasswordByIdApi(ctx *gin.Context) {
	var modifyAccountPassword dto.ModifyAccountPassword
	if err := a.BindAndValidateJSON(ctx, &modifyAccountPassword); err != nil {
		return
	}
	if modifyAccountPassword.Password != modifyAccountPassword.ConfirmPassword {
		a.Fail(ctx, "两次密码不一致", nil)
		return
	}
	id := ctx.Param("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)
	salt := utils.RandomString(utils.GetRandomNum(5, 10))
	password, err := utils.MakePassword(modifyAccountPassword.Password, salt)
	if err != nil {
		a.Fail(ctx, "修改密码失败", err)
		return
	}
	if err := a.AccountRepository.UpdatePassword(ctx, idInt, password); err != nil {
		a.Fail(ctx, "修改密码失败", err)
		return
	}
	a.Success(ctx, "修改密码成功")
}

// UpdateCurrentAccountPasswordApi 修改当前账号密码
func (a *Account) UpdateCurrentAccountPasswordApi(ctx *gin.Context) {
	accountId := ctx.GetInt64("accountId")
	var modifyCurrentPassword dto.ModifyCurrentPassword
	if err := a.BindAndValidateJSON(ctx, &modifyCurrentPassword); err != nil {
		return
	}
	if modifyCurrentPassword.NewPassword != modifyCurrentPassword.ConfirmPassword {
		a.Fail(ctx, "两次密码不一致", nil)
		return
	}
	accountData, err := a.AccountRepository.GetByID(ctx, accountId)
	if err != nil {
		a.Fail(ctx, "修改密码失败", err)
		return
	}
	isValid := k.CheckPassword(accountData.Password, modifyCurrentPassword.Password)

	if !isValid {
		a.Fail(ctx, "旧密码错误", nil)
		return
	}
	salt := utils.RandomString(utils.GetRandomNum(5, 10))
	password, err2 := utils.MakePassword(modifyCurrentPassword.NewPassword, salt)
	if err2 != nil {
		a.Fail(ctx, "修改密码失败", err2)
		return
	}
	if err := a.AccountRepository.UpdatePassword(ctx, accountId, password); err != nil {
		a.Fail(ctx, "修改密码失败", err)
		return
	}
	a.Success(ctx, "修改密码成功")
}

// UpdateStatusByIdApi 根据id修改状态
func (a *Account) UpdateStatusByIdApi(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt := cast.ToInt64(id)
	accountData, err := a.AccountRepository.GetByID(ctx, idInt)
	if err != nil {
		a.Fail(ctx, "修改状态失败", err)
		return
	}
	status := int64(0)
	if accountData.Status == enum.StatusForbidEnum {
		status = enum.StatusNormalEnum
	} else {
		status = enum.StatusForbidEnum
	}
	if err = a.AccountRepository.UpdateStatus(a.Ctx(ctx), idInt, status); err != nil {
		a.Fail(ctx, "更新失败", err)
		return
	}
	a.Success(ctx, "更新成功")
}

// GetAccountByIdApi 根据id查询数据
func (a *Account) GetAccountByIdApi(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt := cast.ToInt64(id)
	accountData, err := a.AccountRepository.GetByID(ctx, idInt)
	if err != nil {
		a.Fail(ctx, "查询失败", err)
		return
	}
	var resultData = vo.AccountVo{}
	_ = utils.CopyProperties(&resultData, accountData)
	fmt.Println(utils.MapToJson(resultData), "拷贝后数据")
	a.Success(ctx, resultData)
}

// GetAccountPageApi 分页获取数据
func (a *Account) GetAccountPageApi(ctx *gin.Context) {
	username := ctx.DefaultQuery("username", "")
	status := ctx.DefaultQuery("status", "")
	pageNumber, _ := strconv.Atoi(ctx.DefaultQuery("pageNumber", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	offset, limit := a.DbPage(pageNumber, pageSize)
	accountDataList, total, err := a.AccountRepository.GetPage(ctx, username, cast.ToInt64(status), offset, limit)
	if err != nil {
		a.Fail(ctx, "查询数据失败", err)
		return
	}
	var accountList = make([]vo.AccountVo, 0)
	for _, item := range accountDataList {
		var resultData = vo.AccountVo{}
		_ = utils.CopyProperties(&resultData, item)
		fmt.Println(utils.MapToJson(resultData), "拷贝后数据")
		accountList = append(accountList, resultData)
	}
	a.BuildPageData(ctx, accountList, total)
}
func NewAccount(baseApi *base.BaseApi, accountRepository repository.IAccountRepository) IAccount {
	return &Account{
		BaseApi:           baseApi,
		AccountRepository: accountRepository,
	}
}
