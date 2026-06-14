package auth

import (
	"errors"
	"fmt"
	"gin-admin-api/internal/api/auth/dto"
	"gin-admin-api/internal/api/auth/mapper"
	"gin-admin-api/internal/api/auth/vo"
	"gin-admin-api/internal/api/base"
	"gin-admin-api/internal/dal/dao"
	"gin-admin-api/internal/dal/repository"
	"gin-admin-api/pkg/enum"
	"gin-admin-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/kuangshp/go-utils/k"
	"github.com/kuangshp/go-utils/k/captcha"
	"github.com/kuangshp/gorm-plus"
	"gorm.io/gorm"
	"time"
)

type IAuth interface {
	AccountLoginApi(ctx *gin.Context)  // 账号登录
	LogOutApi(ctx *gin.Context)        // 用户退出登录
	GetCaptchaApi(ctx *gin.Context)    // 获取图形验证码图片
	VerifyCaptchaApi(ctx *gin.Context) // 验证图形验证码
}

type Auth struct {
	*base.BaseApi
	SysAccountRepository repository.SysAccountRepository
	AuthMapper           mapper.IAuthMapper
}

// AccountLoginApi 账号登录
// @Summary 账号登录
// @Description 使用账号、邮箱或手机号登录，登录成功后返回账号信息和 token
// @Tags 认证中心
// @Accept json
// @Produce json
// @Param data body dto.AccountLoginDTO true "登录参数"
// @Success 200 {object} vo.AccountLoginVO "统一响应，code=0 时 result 为 vo.AccountLoginVO，code=1 时 result 为 null"
// @Router /api/v1/admin/auth/login [post]
func (a Auth) AccountLoginApi(ctx *gin.Context) {
	req := dto.AccountLoginDTO{}
	if !a.BindAndValidateJSON(ctx, &req) {
		return
	}
	isOk := captcha.Verify(req.CaptchaId, req.CaptchaValue, true)
	if !isOk {
		a.Fail(ctx, errors.New("验证码失败"), "验证码失败")
		return
	}
	accountEntity, err := a.SysAccountRepository.FindOneWrapper(ctx, func(g gormplus.IGenWrapper[dao.ISysAccountEntityDo]) {
		g.WhereOrGroup(
			dao.SysAccountEntity.Username.Eq(req.Username),
			dao.SysAccountEntity.Email.Eq(req.Username),
			dao.SysAccountEntity.Mobile.Eq(req.Username),
		)
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			a.Fail(ctx, err, "用户名或密码错误")
			return
		}
		a.Fail(ctx, err, "登录失败")
		return
	}
	if accountEntity.Status == enum.StatusForbidEnum {
		a.Fail(ctx, errors.New("用户被禁用登录,请联系管理员"), "用户被禁用登录,请联系管理员")
		return
	}
	fmt.Println(k.MapToString(accountEntity))
	fmt.Println("1111", k.CheckPassword(accountEntity.Password, req.Password))
	// 判断账号密码是否正确
	if !k.CheckPassword(accountEntity.Password, req.Password) {
		a.Fail(ctx, errors.New("用户名或密码错误"), "用户名或密码错误")
		return
	}
	clientIP := utils.GetCurrentIP(ctx)

	token, err := utils.GenerateToken(utils.HmacUser{
		Username:  accountEntity.Username,
		AccountId: accountEntity.ID,
		IsAdmin:   accountEntity.IsAdmin,
	})
	if err != nil {
		a.Fail(ctx, err, "登录失败")
		return
	}
	// 修改最后一次登录信息
	if err = a.SysAccountRepository.UpdateById(
		ctx,
		accountEntity.ID,
		gormplus.Update().WithColumns(
			dao.SysAccountEntity.LastLoginIP.Value(clientIP),
			dao.SysAccountEntity.LastLoginDate.Value(time.Now()),
		).Build(),
	); err != nil {
		a.Fail(ctx, errors.New("修改最后一次登录信息失败"), "登录失败")
		return
	}
	redisDb := utils.NewRedisUtils(a.Redis)
	if err = redisDb.SetRedisValue(
		ctx,
		utils.AuthTokenRedisKey(accountEntity.ID, token),
		accountEntity.ID,
		int64(utils.TokenExpiration/time.Second),
	); err != nil {
		a.Fail(ctx, err, "登录失败")
		return
	}
	loginToVo := a.AuthMapper.LoginToVo(accountEntity, clientIP, token)
	a.Success(ctx, loginToVo)
	return
}

// LogOutApi 用户退出登录
// @Summary 用户退出登录
// @Description 删除当前登录用户的 token 缓存
// @Tags 认证中心
// @Accept json
// @Produce json
// @Param token header string true "登录 token"
// @Success 200 {string} string "退出登录成功"
// @Router /api/v1/admin/auth/logout [post]
func (a Auth) LogOutApi(ctx *gin.Context) {
	tokenString := ctx.GetHeader("token")
	accountIDValue, exists := ctx.Get("accountId")
	if !exists {
		a.Fail(ctx, errors.New("未获取到当前登录用户"), "退出登录失败")
		return
	}
	accountID, ok := accountIDValue.(int64)
	if !ok {
		a.Fail(ctx, errors.New("当前登录用户数据错误"), "退出登录失败")
		return
	}
	redisDb := utils.NewRedisUtils(a.Redis)
	redisDb.DelRedisKey(ctx, utils.AuthTokenRedisKey(accountID, tokenString))
	a.Success(ctx, "退出登录成功")
}

// GetCaptchaApi 获取图形验证码图片
// @Summary 获取图形验证码图片
// @Description 生成图形验证码，返回验证码图片 Base64、验证码 ID 和验证码值
// @Tags 认证中心
// @Accept json
// @Produce json
// @Success 200 {object} vo.GetCaptchaVO "统一响应，code=0 时 result 为 vo.GetCaptchaVO，code=1 时 result 为 null"
// @Router /api/v1/admin/auth/captcha [get]
func (a Auth) GetCaptchaApi(ctx *gin.Context) {
	id, s, answer, err := captcha.DriverDigitFunc()
	if err != nil {
		a.Fail(ctx, err, "生成验证码失败")
		return
	}
	captchaVO := vo.GetCaptchaVO{
		Base64:    s,
		CaptchaId: id,
		Code:      answer,
	}
	a.Success(ctx, captchaVO)
	return
}

// VerifyCaptchaApi 验证图形验证码
// @Summary 验证图形验证码
// @Description 根据验证码 ID 和验证码值校验图形验证码
// @Tags 认证中心
// @Accept json
// @Produce json
// @Param data body dto.VerifyCaptchaDTO true "验证码参数"
// @Success 200 {string} string "验证码成功"
// @Router /api/v1/admin/auth/captcha/verify [post]
func (a Auth) VerifyCaptchaApi(ctx *gin.Context) {
	var req dto.VerifyCaptchaDTO
	if !a.BindAndValidateJSON(ctx, &req) {
		return
	}
	isOk := captcha.Verify(req.CaptchaId, req.CaptchaValue, true)
	if !isOk {
		a.Fail(ctx, errors.New("验证码失败"), "验证码失败")
		return
	}
	a.Success(ctx, "验证码成功")
	return
}
func NewAuth(baseApi *base.BaseApi) IAuth {
	return Auth{
		BaseApi:              baseApi,
		SysAccountRepository: repository.NewSysAccountRepository(),
		AuthMapper:           mapper.NewAuthMapper(),
	}
}
