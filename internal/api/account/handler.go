package account

import (
	"errors"
	"fmt"
	"gin-admin-api/internal/config"
	"gin-admin-api/internal/api/account/dto"
	"gin-admin-api/internal/api/account/vo"
	"gin-admin-api/internal/dal/repository"
	"gin-admin-api/pkg/enum"
	"gin-admin-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/kuangshp/go-utils/k"
	"github.com/spf13/cast"
	"go.uber.org/zap"
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
	db                *gorm.DB
	cfg               *config.ServerConfig
	redis             *redis.Client
	logger            *zap.Logger
	accountRepository repository.IAccountRepository
}

// CreateAccountApi 创建账号
func (a Account) CreateAccountApi(ctx *gin.Context) {
	var createAccountDto dto.CreateAccountDto
	if err := ctx.ShouldBindJSON(&createAccountDto); err != nil {
		message := utils.ShowErrorMessage(err)
		utils.Fail(ctx, message)
		return
	}
	// 1.判断两次密码是否一致
	if createAccountDto.Password != createAccountDto.ConfirmPassword {
		utils.Fail(ctx, "两次密码不一致")
		return
	}
	// 2.判断账号是否已经存在
	exists, err := a.accountRepository.ExistsByUsername(ctx, createAccountDto.Username)
	if err != nil {
		a.logger.Error("查询账号失败" + err.Error())
		utils.Fail(ctx, "创建账号失败")
		return
	}
	if exists {
		utils.Fail(ctx, fmt.Sprintf("%s已经存在,不能重复创建", createAccountDto.Username))
		return
	}
	// 3.对密码加密
	salt := utils.RandomString(utils.GetRandomNum(5, 10))
	password, err := utils.MakePassword(createAccountDto.Password, salt)
	if err != nil {
		a.logger.Error("密码加密失败" + err.Error())
		utils.Fail(ctx, "创建账号失败")
		return
	}
	if err := a.accountRepository.Create(ctx, createAccountDto.Username, password, salt); err != nil {
		a.logger.Error("创建账号失败" + err.Error())
		utils.Fail(ctx, "创建失败")
		return
	}
	utils.Success(ctx, "创建成功")
	return
}

// LoginAccountApi 用户名和密码登录
func (a Account) LoginAccountApi(ctx *gin.Context) {
	var accountDto dto.AccountDto
	if err := ctx.ShouldBindJSON(&accountDto); err != nil {
		message := utils.ShowErrorMessage(err)
		utils.Fail(ctx, message)
		return
	}
	// 根据账号查询数据
	result, err := a.accountRepository.GetByUsername(ctx, accountDto.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			a.logger.Error("账号不存在" + accountDto.Username)
			utils.Fail(ctx, "账号或密码错误")
			return
		}
		a.logger.Error("查询账号失败" + err.Error())
		utils.Fail(ctx, "账号或密码错误")
		return
	}
	if result.Status == enum.StatusForbidEnum {
		utils.Fail(ctx, "当前账号已经被禁用,请联系管理员")
		return
	}
	isValid := k.CheckPassword(result.Password, accountDto.Password)
	if !isValid {
		utils.Fail(ctx, "账号或密码错误")
		return
	}
	// 3.生产token返回给前端
	hmacUser := utils.HmacUser{
		AccountId: result.ID,
		Username:  result.Username,
	}
	token, err := utils.GenerateToken(hmacUser)
	if err != nil {
		a.logger.Error("生成token失败")
		utils.Fail(ctx, "账号或密码错误")
		return
	}
	// 更新账号登录信息
	if err := a.accountRepository.UpdateLoginInfo(ctx, result.ID, ctx.ClientIP()); err != nil {
		a.logger.Error("更新登录信息失败")
		utils.Fail(ctx, "账号或密码错误")
		return
	}
	utils.Success(ctx, vo.LoginVo{
		AccountVo: vo.AccountVo{
			ID:            result.ID,
			CreatedAt:     result.CreatedAt,
			UpdatedAt:     result.UpdatedAt,
			Username:      result.Username, // 用户名
			Name:          result.Name,     // 真实姓名
			Mobile:        result.Mobile,   // 手机号码
			Email:         result.Email,    // 邮箱地址
			Avatar:        result.Avatar,   // 用户头像
			IsAdmin:       result.IsAdmin,  // 是否为超级管理员:0否,1是
			Status:        result.Status,   // 状态1是正常,0是禁用
			LastLoginDate: time.Now(),
			LastLoginIP:   result.LastLoginIP,
		},
		Token: token,
	})
	return
}

// DeleteAccountByIdApi 根据id删除数据
func (a Account) DeleteAccountByIdApi(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)
	// 1.判断超级管理员不能删除
	accountData, err := a.accountRepository.GetByID(ctx, idInt)
	if err != nil {
		a.logger.Error("根据id查询数据失败" + err.Error())
		utils.Fail(ctx, "删除失败")
		return
	}
	if accountData.IsAdmin == enum.AdminAccount {
		utils.Fail(ctx, "超级管理员不能被删除")
		return
	}
	// 2.判断不能自己删除自己
	if accountData.ID == idInt {
		utils.Fail(ctx, "自己不能删除自己")
		return
	}
	if err := a.accountRepository.Delete(ctx, idInt); err != nil {
		utils.Fail(ctx, "删除失败")
		return
	}
	utils.Success(ctx, "删除成功")
	return
}

// ModifyPasswordByIdApi 根据id修改密码
func (a Account) ModifyPasswordByIdApi(ctx *gin.Context) {
	var modifyAccountPassword dto.ModifyAccountPassword
	if err := ctx.ShouldBindJSON(&modifyAccountPassword); err != nil {
		message := utils.ShowErrorMessage(err)
		utils.Fail(ctx, message)
		return
	}
	if modifyAccountPassword.Password != modifyAccountPassword.ConfirmPassword {
		utils.Fail(ctx, "两次密码不一致")
		return
	}
	id := ctx.Param("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)
	salt := utils.RandomString(utils.GetRandomNum(5, 10))
	password, err := utils.MakePassword(modifyAccountPassword.Password, salt)
	if err != nil {
		a.logger.Error("密码加密失败" + err.Error())
		utils.Fail(ctx, "修改密码失败")
		return
	}
	if err := a.accountRepository.UpdatePassword(ctx, idInt, password); err != nil {
		utils.Fail(ctx, "修改密码失败")
		return
	}
	utils.Success(ctx, "修改密码成功")
	return
}

// UpdateCurrentAccountPasswordApi 修改当前账号密码
func (a Account) UpdateCurrentAccountPasswordApi(ctx *gin.Context) {
	accountId := ctx.GetInt64("accountId")
	var modifyCurrentPassword dto.ModifyCurrentPassword
	if err := ctx.ShouldBindJSON(&modifyCurrentPassword); err != nil {
		message := utils.ShowErrorMessage(err)
		utils.Fail(ctx, message)
		return
	}
	if modifyCurrentPassword.NewPassword != modifyCurrentPassword.ConfirmPassword {
		utils.Fail(ctx, "两次密码不一致")
		return
	}
	accountData, err := a.accountRepository.GetByID(ctx, accountId)
	if err != nil {
		a.logger.Error("根据id查询数据失败" + err.Error())
		utils.Fail(ctx, "修改密码失败")
		return
	}
	isValid := k.CheckPassword(accountData.Password, modifyCurrentPassword.Password)

	if !isValid {
		utils.Fail(ctx, "旧密码错误")
		return
	}
	salt := utils.RandomString(utils.GetRandomNum(5, 10))
	password, err2 := utils.MakePassword(modifyCurrentPassword.NewPassword, salt)
	if err2 != nil {
		a.logger.Error("密码加密失败" + err2.Error())
		utils.Fail(ctx, "修改密码失败")
		return
	}
	if err := a.accountRepository.UpdatePassword(ctx, accountId, password); err != nil {
		utils.Fail(ctx, "修改密码失败")
		return
	}
	utils.Success(ctx, "修改密码成功")
	return
}

// UpdateStatusByIdApi 根据id修改状态
func (a Account) UpdateStatusByIdApi(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt := cast.ToInt64(id)
	accountData, err := a.accountRepository.GetByID(ctx, idInt)
	if err != nil {
		a.logger.Error("根据id查询数据失败" + err.Error())
		utils.Fail(ctx, "修改状态失败")
		return
	}
	status := int64(0)
	if accountData.Status == enum.StatusForbidEnum {
		status = enum.StatusNormalEnum
	} else {
		status = enum.StatusForbidEnum
	}
	if err := a.accountRepository.UpdateStatus(ctx, idInt, status); err != nil {
		utils.Fail(ctx, "更新失败")
		return
	}
	utils.Success(ctx, "更新成功")
	return
}

// GetAccountByIdApi 根据id查询数据
func (a Account) GetAccountByIdApi(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt := cast.ToInt64(id)
	accountData, err := a.accountRepository.GetByID(ctx, idInt)
	if err != nil {
		utils.Fail(ctx, "查询失败")
		return
	}
	var resultData = vo.AccountVo{}
	_ = utils.CopyProperties(&resultData, accountData)
	fmt.Println(utils.MapToJson(resultData), "拷贝后数据")
	utils.Success(ctx, resultData)
	return
}

// GetAccountPageApi 分页获取数据
func (a Account) GetAccountPageApi(ctx *gin.Context) {
	username := ctx.DefaultQuery("username", "")
	status := ctx.DefaultQuery("status", "")
	pageNumber, _ := strconv.Atoi(ctx.DefaultQuery("pageNumber", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if pageNumber <= 0 {
		pageNumber = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (pageNumber - 1) * pageSize
	accountDataList, total, err := a.accountRepository.GetPage(ctx, username, cast.ToInt64(status), offset, pageSize)
	if err != nil {
		a.logger.Error("查询数据失败" + err.Error())
	}
	var accountList = make([]vo.AccountVo, 0)
	for _, item := range accountDataList {
		var resultData = vo.AccountVo{}
		_ = utils.CopyProperties(&resultData, item)
		fmt.Println(utils.MapToJson(resultData), "拷贝后数据")
		accountList = append(accountList, resultData)
	}
	utils.BuildPageData(ctx, accountList, total)
	return
}

func NewAccount(db *gorm.DB, cfg *config.ServerConfig, redis *redis.Client, logger *zap.Logger, accountRepository repository.IAccountRepository) IAccount {
	return Account{
		db:                db,
		cfg:               cfg,
		redis:             redis,
		logger:            logger,
		accountRepository: accountRepository,
	}
}
