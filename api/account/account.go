package account

import (
	"database/sql"
	"gin-admin-api/api/account/dto"
	"gin-admin-api/enum"
	"gin-admin-api/global"
	"gin-admin-api/model"
	"gin-admin-api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

type IAccount interface {
	Register(ctx *gin.Context)                     // 用户注册
	Login(ctx *gin.Context)                        // 用户登录
	DeleteAccountById(ctx *gin.Context)            // 根据id修改账号
	ModifyPasswordById(ctx *gin.Context)           // 根据id修改账号密码
	UpdateStatusById(ctx *gin.Context)             // 根据id修改状态
	UpdateCurrentAccountPassword(ctx *gin.Context) // 修改当前账号密码
	GetAccountById(ctx *gin.Context)               // 根据id获取账号信息
	GetAccountPage(ctx *gin.Context)               // 分页获取账号数据
}

type Account struct {
	db *gorm.DB
}

func (a Account) Register(ctx *gin.Context) {
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
	// 2.对密码加密
	password, err := utils.MakePassword(createAccountDto.Password)
	if err != nil {
		global.Logger.Error("密码加密失败" + err.Error())
		utils.Fail(ctx, "创建账号失败")
		return
	}
	// 3.创建账号信息
	if result := a.db.Create(&model.AccountEntity{
		Username: createAccountDto.UserName,
		Password: password,
		Status:   sql.NullInt64{Valid: true, Int64: enum.Normal},
	}).Error; result != nil {
		global.Logger.Error("创建账号失败" + result.Error())
		utils.Fail(ctx, "创建账号失败")
	}
	utils.Success(ctx, "创建成功")
	return
}

func (a Account) Login(ctx *gin.Context) {
	var accountDto dto.AccountDto
	if err := ctx.ShouldBindJSON(&accountDto); err != nil {
		message := utils.ShowErrorMessage(err)
		utils.Fail(ctx, message)
		return
	}
	// 1.根据账号名去查询密码信息
	var accountEntity model.AccountEntity
	if result := a.db.Where("username=?", accountDto.UserName).Select([]string{"password", "id", "username", "status"}).First(&accountEntity); result.RowsAffected == 0 {
		global.Logger.Error("根据用户名查询数据失败" + result.Error.Error())
		utils.Fail(ctx, "账号或密码错误")
		return
	}
	if accountEntity.Status.Int64 == enum.Forbid {
		utils.Fail(ctx, "当前账号不允许登录,请联系管理员")
		return
	}
	// 2.判断密码是否正确
	isOk, err := utils.CheckPassword(accountEntity.Password, accountDto.Password)
	if err != nil {
		global.Logger.Error("校验密码错误" + err.Error())
		utils.Fail(ctx, "账号或密码错误")
		return
	}
	if !isOk {
		utils.Fail(ctx, "账号或密码错误")
		return
	}
	// 3.生产token返回给前端
	hmacUser := utils.HmacUser{
		Id:       int(accountEntity.Id),
		Username: accountEntity.Username,
	}
	if token, err := utils.GenerateToken(hmacUser); err == nil {
		utils.Success(ctx, gin.H{
			"id":       accountEntity.Id,
			"username": accountEntity.Username,
			"token":    token,
		})
		return
	} else {
		global.Logger.Error("生成token失败")
		utils.Fail(ctx, "账号或密码错误")
		return
	}
}

func (a Account) DeleteAccountById(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)
	if result := a.db.Where("id=?", idInt).Delete(&model.AccountEntity{}).Error; result != nil {
		global.Logger.Error("根据id删除账号失败" + result.Error())
		utils.Fail(ctx, "删除失败")
		return
	}
	utils.Success(ctx, "删除成功")
	return
}

func (a Account) ModifyPasswordById(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)

	var modifyAccountPassword dto.ModifyAccountPassword
	if err := ctx.ShouldBindJSON(&modifyAccountPassword); err != nil {
		message := utils.ShowErrorMessage(err)
		utils.Fail(ctx, message)
		return
	}
	// 1.判断两次密码是否一致
	if modifyAccountPassword.Password != modifyAccountPassword.ConfirmPassword {
		utils.Fail(ctx, "两次密码不一致")
		return
	}
	// 2.对密码加密
	password, err := utils.MakePassword(modifyAccountPassword.Password)
	if err != nil {
		global.Logger.Error("密码加密失败" + err.Error())
		utils.Fail(ctx, "创建账号失败")
		return
	}
	if result := a.db.Where("id=?", idInt).Updates(&model.AccountEntity{
		Password: password,
	}).Error; result != nil {
		global.Logger.Error("修改密码失败" + result.Error())
		utils.Fail(ctx, "修改密码失败")
		return
	}
	utils.Success(ctx, "修改密码成功")
	return
}

func (a Account) UpdateStatusById(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a Account) UpdateCurrentAccountPassword(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a Account) GetAccountById(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a Account) GetAccountPage(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func NewAccount(db *gorm.DB) IAccount {
	return Account{
		db: db,
	}
}
