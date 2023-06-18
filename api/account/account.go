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
	"time"
)

type IAccount interface {
	CreateAccountApi(ctx *gin.Context) // 用户注册
	LoginAccountApi(ctx *gin.Context)  // 用户名和密码登录
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
		fmt.Println("1111")
		// 3.生产token返回给前端
		hmacUser := utils.HmacUser{
			AccountId: result.ID,
			Username:  result.Username,
		}
		if token, err := utils.GenerateToken(hmacUser); err == nil {
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
		} else {
			global.Logger.Error("生成token失败")
			utils.Fail(ctx, "账号或密码错误")
			return
		}
	}
}

func NewAccount(db *dao.Query) IAccount {
	return Account{
		db: db,
	}
}
