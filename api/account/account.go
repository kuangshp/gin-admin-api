package account

import (
	"fmt"
	"gin-admin-api/api/account/dto"
	"gin-admin-api/api/account/vo"
	"gin-admin-api/dao"
	"gin-admin-api/enum"
	"gin-admin-api/global"
	"gin-admin-api/model"
	"gin-admin-api/utils"
	"github.com/gin-gonic/gin"
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
}

type Account struct {
	db *dao.Query
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
	var queryAccountBuilder = a.db.Account
	// 2.判断账号是否已经存在
	if result, err := queryAccountBuilder.WithContext(ctx).Where(queryAccountBuilder.Username.Eq(createAccountDto.Username)).
		Select(queryAccountBuilder.ID, queryAccountBuilder.Username).First(); err != gorm.ErrRecordNotFound {
		utils.Fail(ctx, fmt.Sprintf("%s已经存在,不能重复创建", result.Username))
		return
	}
	// 3.对密码加密
	salt := utils.RandomString(utils.GetRandomNum(5, 10))
	password, err := utils.MakePassword(createAccountDto.Password, salt)
	if err != nil {
		global.Logger.Error("密码加密失败" + err.Error())
		utils.Fail(ctx, "创建账号失败")
		return
	}
	if result := queryAccountBuilder.Select(queryAccountBuilder.Username, queryAccountBuilder.Password, queryAccountBuilder.Salt,
		queryAccountBuilder.Status, queryAccountBuilder.IsAdmin).
		Create(&model.Account{
			Username: createAccountDto.Username,
			Password: password,
			Salt:     salt,
			Status:   enum.Normal,
			IsAdmin:  enum.NormalAccount,
		}); result != nil {
		global.Logger.Error("创建账号失败" + result.Error())
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
	var queryAccountBuilder = a.db.Account
	if result, err := queryAccountBuilder.Where(queryAccountBuilder.Username.Eq(accountDto.Username)).First(); err == gorm.ErrRecordNotFound {
		global.Logger.Error("账号不存在" + accountDto.Username)
		utils.Fail(ctx, "账号或密码错误")
		return
	} else {
		if result.Status == enum.Forbid {
			utils.Fail(ctx, "当前账号已经被禁用,请联系管理员")
			return
		}
		isValid, err1 := utils.CheckPassword(result.Password, accountDto.Password, result.Salt)
		if err1 != nil || !isValid {
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
			global.Logger.Error("生成token失败")
			utils.Fail(ctx, "账号或密码错误")
			return
		}
		// 更新账号
		if _, err := queryAccountBuilder.Where(queryAccountBuilder.ID.Eq(result.ID)).
			Select(queryAccountBuilder.ExpireTime, queryAccountBuilder.Token, queryAccountBuilder.LastLoginDate, queryAccountBuilder.LastLoginIP).
			Updates(&model.Account{
				Token:         token,
				ExpireTime:    model.LocalTime{Time: time.Now().Add(7 * time.Hour * 24)},
				LastLoginDate: model.LocalTime{Time: time.Now()},
				LastLoginIP:   ctx.ClientIP(), //最后登录id
			}); err != nil {
			global.Logger.Error("更新表的时候失败")
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
				LastLoginDate: model.LocalTime{Time: time.Now()},
				LastLoginIP:   result.LastLoginIP,
			},
			Token: token,
		})
		return
	}
}

// DeleteAccountByIdApi 根据id删除数据
func (a Account) DeleteAccountByIdApi(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)
	var queryAccountBuilder = a.db.Account
	// 1.判断超级管理员不能删除
	accountData, err := queryAccountBuilder.Where(queryAccountBuilder.ID.Eq(idInt)).First()
	if err != nil {
		global.Logger.Error("根据id查询数据失败" + err.Error())
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
	if _, err := queryAccountBuilder.Where(queryAccountBuilder.ID.Eq(idInt)).Delete(); err == nil {
		utils.Success(ctx, "删除成功")
		return
	} else {
		utils.Fail(ctx, "删除失败")
		return
	}
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
		global.Logger.Error("密码加密失败" + err.Error())
		utils.Fail(ctx, "修改密码失败")
		return
	}
	var queryAccountBuilder = a.db.Account
	if _, err := queryAccountBuilder.Where(queryAccountBuilder.ID.Eq(idInt)).
		Select(queryAccountBuilder.Password, queryAccountBuilder.Salt).
		Updates(&model.Account{
			Password: password,
			Salt:     salt,
		}); err != nil {
		utils.Fail(ctx, "修改密码失败")
		return
	}
	utils.Success(ctx, "修改密码成功")
	return
}

// UpdateCurrentAccountPasswordApi 修改当前账号密码
func (a Account) UpdateCurrentAccountPasswordApi(ctx *gin.Context) {
	accountId := ctx.GetInt64("accountId")
	fmt.Println(accountId, "====")
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
	var queryAccountBuilder = a.db.Account
	accountData, err := queryAccountBuilder.Select(queryAccountBuilder.Password, queryAccountBuilder.Salt).Where(queryAccountBuilder.ID.Eq(accountId)).First()
	if err != nil {
		global.Logger.Error("根据id查询数据失败" + err.Error())
		utils.Fail(ctx, "修改密码失败")
		return
	}
	isValid, err1 := utils.CheckPassword(accountData.Password, modifyCurrentPassword.Password, accountData.Salt)
	if err1 != nil {
		global.Logger.Error("校验旧密码失败" + err1.Error())
		utils.Fail(ctx, "修改密码失败")
		return
	}
	if !isValid {
		utils.Fail(ctx, "旧密码错误")
		return
	}
	salt := utils.RandomString(utils.GetRandomNum(5, 10))
	password, err := utils.MakePassword(modifyCurrentPassword.NewPassword, salt)
	if err != nil {
		global.Logger.Error("密码加密失败" + err.Error())
		utils.Fail(ctx, "修改密码失败")
		return
	}
	if _, err := queryAccountBuilder.Where(queryAccountBuilder.ID.Eq(accountId)).
		Select(queryAccountBuilder.Password, queryAccountBuilder.Salt).
		Updates(&model.Account{
			Password: password,
			Salt:     salt,
		}); err != nil {
		utils.Fail(ctx, "修改密码失败")
		return
	}
	utils.Success(ctx, "修改密码成功")
	return
}

// UpdateStatusByIdApi 根据id修改状态
func (a Account) UpdateStatusByIdApi(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)
	var queryAccountBuilder = a.db.Account
	accountData, err := queryAccountBuilder.Where(queryAccountBuilder.ID.Eq(idInt)).Select(queryAccountBuilder.Status).First()
	if err != nil {
		global.Logger.Error("根据id查询数据失败" + err.Error())
		utils.Fail(ctx, "修改状态失败")
		return
	}
	status := 0
	if accountData.Status == enum.Forbid {
		status = enum.Normal
	} else {
		status = enum.Forbid
	}
	if _, err1 := queryAccountBuilder.Where(queryAccountBuilder.ID.Eq(idInt)).Updates(map[string]interface{}{
		"status": status,
	}); err1 != nil {
		utils.Fail(ctx, "更新失败")
		return
	}
	utils.Success(ctx, "更新成功")
	return
}

// GetAccountByIdApi 根据id查询数据
func (a Account) GetAccountByIdApi(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)
	var queryAccountBuilder = a.db.Account
	accountData, err := queryAccountBuilder.Where(queryAccountBuilder.ID.Eq(idInt)).Omit(queryAccountBuilder.Password).First()
	if err != nil {
		utils.Fail(ctx, "查询失败")
		return
	}
	utils.GetIpToAddress("106.226.190.66")
	utils.Success(ctx, vo.AccountVo{
		ID:            accountData.ID,
		CreatedAt:     accountData.CreatedAt,
		UpdatedAt:     accountData.UpdatedAt,
		Username:      accountData.Username,      // 用户名
		Name:          accountData.Name,          // 真实姓名
		Mobile:        accountData.Mobile,        // 手机号码
		Email:         accountData.Email,         // 邮箱地址
		Avatar:        accountData.Avatar,        // 用户头像
		IsAdmin:       accountData.IsAdmin,       // 是否为超级管理员:0否,1是
		Status:        accountData.Status,        // 状态1是正常,0是禁用
		LastLoginIP:   accountData.LastLoginIP,   // 最后登录ip地址
		LastLoginDate: accountData.LastLoginDate, // 最后登录时间
	})
	return
}
func NewAccount(db *dao.Query) IAccount {
	return Account{
		db: db,
	}
}
