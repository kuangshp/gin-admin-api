package auth

import (
	"errors"
	"gin-admin-api/internal/api/base"
	"gin-admin-api/internal/dal/dao"
	"gin-admin-api/internal/dal/dto"
	"gin-admin-api/internal/dal/mapper"
	"gin-admin-api/internal/dal/repository"
	"gin-admin-api/internal/dal/vo"
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
	SysAccountMapper     mapper.ISysAccountMapper
}

func (a Auth) AccountLoginApi(ctx *gin.Context) {
	req := dto.AccountLoginDTO{}
	if err := a.BindAndValidateJSON(ctx, &req); err != nil {
		return
	}
	isOk := captcha.Verify(req.CaptchaId, req.CaptchaValue, true)
	if !isOk {
		a.Fail(ctx, errors.New("验证码失败"), "验证码失败")
		return
	}
	accountEntity, err := a.SysAccountRepository.FindOneWrapper(ctx, func(g gormplus.IGenWrapper[dao.ISysAccountEntityDo]) {
		g.OrGroup(
			dao.SysAccountEntity.Username.Eq(req.Username),
			dao.SysAccountEntity.Email.Eq(req.Username),
			dao.SysAccountEntity.Mobile.Eq(req.Username),
		)
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		a.Fail(ctx, err, "用户名或密码错误")
		return
	}
	if accountEntity.Status == enum.StatusForbidEnum {
		err = errors.New("用户被禁用登录,请联系管理员")
		return
	}
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
		dao.SysAccountEntity.LastLoginIP.Value(clientIP),
		dao.SysAccountEntity.LastLoginDate.Value(time.Now()),
	); err != nil {
		a.Fail(ctx, errors.New("修改最后一次登录信息失败"), "登录失败")
		return
	}
	loginToVo := a.SysAccountMapper.LoginToVo(accountEntity, clientIP, token)
	a.Success(ctx, loginToVo)
	return
}

func (a Auth) LogOutApi(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

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

func (a Auth) VerifyCaptchaApi(ctx *gin.Context) {
	var req dto.VerifyCaptchaDTO
	if err := a.BindAndValidateJSON(ctx, &req); err != nil {
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
func NewAuth(baseApi *base.BaseApi, accountRepository repository.SysAccountRepository, sysAccountMapper mapper.ISysAccountMapper) IAuth {
	return Auth{
		BaseApi:              baseApi,
		SysAccountRepository: accountRepository,
		SysAccountMapper:     sysAccountMapper,
	}
}
