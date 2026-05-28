package account

import (
	"errors"
	"gin-admin-api/internal/api/base"
	"gin-admin-api/internal/dal/dao"
	"gin-admin-api/internal/dal/dto"
	"gin-admin-api/internal/dal/mapper"
	"gin-admin-api/internal/dal/model"
	"gin-admin-api/internal/dal/repository"
	"gin-admin-api/pkg/enum"
	"gin-admin-api/pkg/utils"
	"strconv"

	"github.com/kuangshp/go-utils/k"

	"github.com/gin-gonic/gin"
	gormplus "github.com/kuangshp/gorm-plus"
	"gorm.io/gorm"
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
	SysAccountRepository     repository.SysAccountRepository
	SysAccountRoleRepository repository.SysAccountRoleRepository
	SysAccountMapper         mapper.ISysAccountMapper
}

func (s SysAccount) CreateSysAccountApi(ctx *gin.Context) {
	var req dto.CreateSysAccountDTO
	if !s.BindAndValidateJSON(ctx, &req) {
		return
	}
	// 判断是否已经存在
	if !s.checkAccountUnique(ctx, req.Username, req.Mobile, req.Email, 0) {
		return
	}
	password, err := k.MakePassword(req.Password)
	if err != nil {
		s.Logger.Error("密码加密失败")
		s.Fail(ctx, err, "创建失败")
		return
	}

	if err = gormplus.TransactionAs(s.Db, func(db *gorm.DB) *dao.Query {
		return dao.Use(db)
	}, func(tx *dao.Query) error {
		accountEntity := s.SysAccountMapper.DtoToEntity(&req, password, enum.StatusNormalEnum)
		if err = s.SysAccountRepository.CreateTx(ctx, tx, accountEntity, dao.SysAccountEntity.LastLoginIP, dao.SysAccountEntity.LastLoginDate); err != nil {
			s.Logger.Error("创建失败")
			return err
		}
		// 2.分配角色
		if len(req.RoleIdList) > 0 {
			var accountRoleEntity = make([]*model.SysAccountRoleEntity, 0)
			for _, item := range req.RoleIdList {
				accountRoleEntity = append(accountRoleEntity, &model.SysAccountRoleEntity{
					RoleID:    item,
					AccountID: accountEntity.ID,
				})
			}
			if err = s.SysAccountRoleRepository.CreateBatchTx(ctx, tx, accountRoleEntity); err != nil {
				s.Logger.Error("创建账号角色失败")
				return err
			}
		}
		return nil
	}); err != nil {
		s.Fail(ctx, err, "创建失败")
		return
	}
	s.Success(ctx, "创建成功")
}

// DeleteSysAccountByIdApi 删除账号
// @Summary 删除账号
// @Description 根据账号 ID 删除账号，并清理账号角色关系
// @Tags 账号中心
// @Accept json
// @Produce json
// @Param id path int true "账号ID"
// @Success 200 {object} gin.H "成功响应"
// @Failure 200 {object} gin.H "失败响应"
// @Router /api/v1/admin/account/{id} [delete]
func (s SysAccount) DeleteSysAccountByIdApi(ctx *gin.Context) {
	id, ok := s.getIDParam(ctx)
	if !ok {
		return
	}
	if _, err := s.SysAccountRepository.FindById(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Fail(ctx, err, "账号不存在")
			return
		}
		s.Fail(ctx, err, "删除失败")
		return
	}
	if err := gormplus.TransactionAsCtx(ctx, s.Db, func(db *gorm.DB) *dao.Query {
		return dao.Use(db)
	}, func(tx *dao.Query) error {
		if err := s.SysAccountRoleRepository.DeleteByWrapperTx(ctx, tx, func(g gormplus.IGenWrapper[dao.ISysAccountRoleEntityDo]) {
			g.WhereIf(true, dao.SysAccountRoleEntity.AccountID.Eq(id))
		}); err != nil {
			return err
		}
		return s.SysAccountRepository.DeleteByIdTx(ctx, tx, id)
	}); err != nil {
		s.Fail(ctx, err, "删除失败")
		return
	}
	s.Success(ctx, "删除成功")
}

// ModifySysAccountByIdApi 修改账号
// @Summary 修改账号
// @Description 根据账号 ID 修改账号基础信息
// @Tags 账号中心
// @Accept json
// @Produce json
// @Param id path int true "账号ID"
// @Param data body dto.ModifySysAccountDTO true "修改账号参数"
// @Success 200 {object} gin.H "成功响应"
// @Failure 200 {object} gin.H "失败响应"
// @Router /api/v1/admin/account/modify/{id} [put]
// @Router /api/v1/admin/account/modify/{id} [patch]
func (s SysAccount) ModifySysAccountByIdApi(ctx *gin.Context) {
	id, ok := s.getIDParam(ctx)
	if !ok {
		return
	}
	var req dto.ModifySysAccountDTO
	if !s.BindAndValidateJSON(ctx, &req) {
		return
	}
	req.ID = id
	if _, err := s.SysAccountRepository.FindById(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Fail(ctx, err, "账号不存在")
			return
		}
		s.Fail(ctx, err, "修改失败")
		return
	}
	if !s.checkAccountUnique(ctx, req.Username, req.Mobile, req.Email, id) {
		return
	}
	if err := s.SysAccountRepository.UpdateById(
		ctx,
		id,
		dao.SysAccountEntity.Username.Value(req.Username),
		dao.SysAccountEntity.Email.Value(req.Email),
		dao.SysAccountEntity.Mobile.Value(req.Mobile),
		dao.SysAccountEntity.Status.Value(req.Status),
		dao.SysAccountEntity.Avatar.Value(req.Avatar),
	); err != nil {
		s.Fail(ctx, err, "修改失败")
		return
	}
	s.Success(ctx, "修改成功")
}

// GetSysAccountPageApi 分页获取账号
// @Summary 分页获取账号
// @Description 根据查询条件分页获取账号列表
// @Tags 账号中心
// @Accept json
// @Produce json
// @Param pageNumber query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param username query string false "登录账号"
// @Param mobile query string false "手机号"
// @Param email query string false "邮箱"
// @Param status query int false "状态：1正常，2禁用"
// @Success 200 {object} gin.H "成功响应，result 为分页数据"
// @Failure 200 {object} gin.H "失败响应"
// @Router /api/v1/admin/account [get]
func (s SysAccount) GetSysAccountPageApi(ctx *gin.Context) {
	pageSize, pageNumber := utils.GetQueryPage(ctx.Request)
	status, _ := strconv.ParseInt(ctx.Query("status"), 10, 64)
	username := ctx.Query("username")
	mobile := ctx.Query("mobile")
	email := ctx.Query("email")

	list, total, err := s.SysAccountRepository.FindPageByWrapper(ctx, pageNumber, pageSize, func(g gormplus.IGenWrapper[dao.ISysAccountEntityDo]) {
		g.RawWhereIf(username != "", "username LIKE ?", "%"+username+"%").
			RawWhereIf(mobile != "", "mobile LIKE ?", "%"+mobile+"%").
			RawWhereIf(email != "", "email LIKE ?", "%"+email+"%").
			WhereIf(status != 0, dao.SysAccountEntity.Status.Eq(status))
	})
	if err != nil {
		s.Fail(ctx, err, "获取账号分页失败")
		return
	}
	s.BuildPageData(ctx, s.SysAccountMapper.EntityListToVo(list), total)
}

// GetSysAccountListApi 获取账号列表
// @Summary 获取账号列表
// @Description 根据查询条件获取账号列表
// @Tags 账号中心
// @Accept json
// @Produce json
// @Param username query string false "登录账号"
// @Param mobile query string false "手机号"
// @Param email query string false "邮箱"
// @Param status query int false "状态：1正常，2禁用"
// @Success 200 {object} gin.H "成功响应，result 为账号列表"
// @Failure 200 {object} gin.H "失败响应"
// @Router /api/v1/admin/account/list [get]
func (s SysAccount) GetSysAccountListApi(ctx *gin.Context) {
	status, _ := strconv.ParseInt(ctx.Query("status"), 10, 64)
	username := ctx.Query("username")
	mobile := ctx.Query("mobile")
	email := ctx.Query("email")

	list, err := s.SysAccountRepository.FindListByWrapper(ctx, func(g gormplus.IGenWrapper[dao.ISysAccountEntityDo]) {
		g.RawWhereIf(username != "", "username LIKE ?", "%"+username+"%").
			RawWhereIf(mobile != "", "mobile LIKE ?", "%"+mobile+"%").
			RawWhereIf(email != "", "email LIKE ?", "%"+email+"%").
			WhereIf(status != 0, dao.SysAccountEntity.Status.Eq(status))
	})
	if err != nil {
		s.Fail(ctx, err, "获取账号列表失败")
		return
	}
	s.Success(ctx, s.SysAccountMapper.EntityListToVo(list))
}

// GetSysAccountDetailApi 获取账号详情
// @Summary 获取账号详情
// @Description 根据账号 ID 获取账号详情
// @Tags 账号中心
// @Accept json
// @Produce json
// @Param id path int true "账号ID"
// @Success 200 {object} gin.H "成功响应，result 为 vo.SysAccountVO"
// @Failure 200 {object} gin.H "失败响应"
// @Router /api/v1/admin/account/{id} [get]
func (s SysAccount) GetSysAccountDetailApi(ctx *gin.Context) {
	id, ok := s.getIDParam(ctx)
	if !ok {
		return
	}
	accountEntity, err := s.SysAccountRepository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Fail(ctx, err, "账号不存在")
			return
		}
		s.Fail(ctx, err, "获取账号详情失败")
		return
	}
	s.Success(ctx, s.SysAccountMapper.EntityToVo(accountEntity))
}

// ResetPasswordByIdApi 重置账号密码
// @Summary 重置账号密码
// @Description 根据账号 ID 重置账号密码
// @Tags 账号中心
// @Accept json
// @Produce json
// @Param id path int true "账号ID"
// @Param data body dto.ResetPasswordDTO true "重置密码参数"
// @Success 200 {object} gin.H "成功响应"
// @Failure 200 {object} gin.H "失败响应"
// @Router /api/v1/admin/account/modifyPassword/{id} [put]
// @Router /api/v1/admin/account/modifyPassword/{id} [patch]
func (s SysAccount) ResetPasswordByIdApi(ctx *gin.Context) {
	id, ok := s.getIDParam(ctx)
	if !ok {
		return
	}
	var req dto.ResetPasswordDTO
	if !s.BindAndValidateJSON(ctx, &req) {
		return
	}
	if req.Password != req.ConfirmPassword {
		s.Fail(ctx, errors.New("两次密码不一致"), "两次密码不一致")
		return
	}
	if _, err := s.SysAccountRepository.FindById(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Fail(ctx, err, "账号不存在")
			return
		}
		s.Fail(ctx, err, "重置密码失败")
		return
	}
	password, err := k.MakePassword(req.Password)
	if err != nil {
		s.Fail(ctx, err, "重置密码失败")
		return
	}
	if err = s.SysAccountRepository.UpdateById(ctx, id, dao.SysAccountEntity.Password.Value(password)); err != nil {
		s.Fail(ctx, err, "重置密码失败")
		return
	}
	s.Success(ctx, "重置密码成功")
}

// ModifyCurrentSysAccountPasswordApi 修改当前账号密码
// @Summary 修改当前账号密码
// @Description 修改当前登录账号密码
// @Tags 账号中心
// @Accept json
// @Produce json
// @Param data body dto.ModifyCurrentPasswordDTO true "修改当前账号密码参数"
// @Success 200 {object} gin.H "成功响应"
// @Failure 200 {object} gin.H "失败响应"
// @Router /api/v1/admin/account/modifyCurrentPassword [put]
// @Router /api/v1/admin/account/modifyCurrentPassword [patch]
func (s SysAccount) ModifyCurrentSysAccountPasswordApi(ctx *gin.Context) {
	accountID, ok := s.getCurrentAccountID(ctx)
	if !ok {
		return
	}
	var req dto.ModifyCurrentPasswordDTO
	if !s.BindAndValidateJSON(ctx, &req) {
		return
	}
	if req.NewPassword != req.ConfirmPassword {
		s.Fail(ctx, errors.New("两次密码不一致"), "两次密码不一致")
		return
	}
	accountEntity, err := s.SysAccountRepository.FindById(ctx, accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Fail(ctx, err, "账号不存在")
			return
		}
		s.Fail(ctx, err, "修改当前账号密码失败")
		return
	}
	if !k.CheckPassword(accountEntity.Password, req.Password) {
		s.Fail(ctx, errors.New("原密码错误"), "原密码错误")
		return
	}
	password, err := k.MakePassword(req.NewPassword)
	if err != nil {
		s.Fail(ctx, err, "修改当前账号密码失败")
		return
	}
	if err = s.SysAccountRepository.UpdateById(ctx, accountID, dao.SysAccountEntity.Password.Value(password)); err != nil {
		s.Fail(ctx, err, "修改当前账号密码失败")
		return
	}
	s.Success(ctx, "修改当前账号密码成功")
}

func (s SysAccount) getIDParam(ctx *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		s.Fail(ctx, err, "参数id错误")
		return 0, false
	}
	return id, true
}

func (s SysAccount) getCurrentAccountID(ctx *gin.Context) (int64, bool) {
	accountIDValue, exists := ctx.Get("accountId")
	if !exists {
		s.Fail(ctx, errors.New("未获取到当前登录用户"), "未获取到当前登录用户")
		return 0, false
	}
	accountID, ok := accountIDValue.(int64)
	if !ok || accountID <= 0 {
		s.Fail(ctx, errors.New("当前登录用户数据错误"), "当前登录用户数据错误")
		return 0, false
	}
	return accountID, true
}

func (s SysAccount) checkAccountUnique(ctx *gin.Context, username, mobile, email string, excludeID int64) bool {
	sql := "username = ?"
	args := []any{username}
	if mobile != "" {
		sql += " OR mobile = ?"
		args = append(args, mobile)
	}
	if email != "" {
		sql += " OR email = ?"
		args = append(args, email)
	}
	existsAccount, err := s.SysAccountRepository.FindOneWrapper(ctx, func(g gormplus.IGenWrapper[dao.ISysAccountEntityDo]) {
		g.WhereIf(excludeID > 0, dao.SysAccountEntity.ID.Neq(excludeID)).
			RawWhere("("+sql+")", args...)
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Fail(ctx, err, "账号重复校验失败")
		return false
	}
	if existsAccount == nil {
		return true
	}
	switch {
	case existsAccount.Username == username:
		s.Fail(ctx, errors.New("用户名称已经存在,不能重复"), "用户名称已经存在,不能重复")
	case mobile != "" && existsAccount.Mobile == mobile:
		s.Fail(ctx, errors.New("手机号码已经存在,不能重复"), "手机号码已经存在,不能重复")
	case email != "" && existsAccount.Email == email:
		s.Fail(ctx, errors.New("邮箱已经存在,不能重复"), "邮箱已经存在,不能重复")
	default:
		s.Fail(ctx, errors.New("账号信息已经存在,不能重复"), "账号信息已经存在,不能重复")
	}
	return false
}

func NewSysAccount(baseApi *base.BaseApi, accountRepository repository.SysAccountRepository, sysAccountRoleRepository repository.SysAccountRoleRepository, sysAccountMapper mapper.ISysAccountMapper) ISysAccount {
	return &SysAccount{
		BaseApi:                  baseApi,
		SysAccountRepository:     accountRepository,
		SysAccountRoleRepository: sysAccountRoleRepository,
		SysAccountMapper:         sysAccountMapper,
	}
}
