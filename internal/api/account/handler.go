package account

import (
	"errors"
	"gin-admin-api/internal/api/account/dto"
	"gin-admin-api/internal/api/account/mapper"
	_ "gin-admin-api/internal/api/account/vo"
	"gin-admin-api/internal/api/base"
	"gin-admin-api/internal/dal/dao"
	"gin-admin-api/internal/dal/model"
	"gin-admin-api/internal/dal/repository"
	"gin-admin-api/pkg/enum"
	"github.com/kuangshp/go-utils/k"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kuangshp/gorm-plus"
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

// CreateSysAccountApi 创建账号
// @Summary 创建账号
// @Description 创建后台账号，并分配角色关系
// @Tags 账号中心
// @Accept json
// @Produce json
// @Param data body dto.CreateSysAccountDTO true "创建账号参数"
// @Success 200 {string} string "创建成功"
// @Router /api/v1/admin/account/register [post]
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

	if err = gormplus.TransactionAsCtx(ctx, s.Db, useQuery, func(tx *dao.Query) error {
		accountEntity := s.SysAccountMapper.DtoToEntity(&req, password, enum.StatusNormalEnum)
		if err = s.SysAccountRepository.CreateTx(ctx, tx, accountEntity, gormplus.Create().WithOmit(
			dao.SysAccountEntity.LastLoginIP,
			dao.SysAccountEntity.LastLoginDate,
		).Build()); err != nil {
			s.Logger.Error("创建失败")
			return err
		}
		// 2.分配角色
		if len(req.RoleIdList) > 0 {
			accountRoleEntity := s.buildAccountRoleEntityList(accountEntity.ID, req.RoleIdList)
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
// @Success 200 {string} string "删除成功"
// @Router /api/v1/admin/account/{id} [delete]
func (s SysAccount) DeleteSysAccountByIdApi(ctx *gin.Context) {
	id, ok := s.getIdParam(ctx)
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
	if err := gormplus.TransactionAsCtx(ctx, s.Db, useQuery, func(tx *dao.Query) error {
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
// @Description 根据账号 ID 修改账号基础信息和角色关系
// @Tags 账号中心
// @Accept json
// @Produce json
// @Param id path int true "账号ID"
// @Param data body dto.ModifySysAccountDTO true "修改账号参数"
// @Success 200 {string} string "修改成功"
// @Router /api/v1/admin/account/modify/{id} [put]
func (s SysAccount) ModifySysAccountByIdApi(ctx *gin.Context) {
	id, ok := s.getIdParam(ctx)
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
	if err := gormplus.TransactionAsCtx(ctx, s.Db, useQuery, func(tx *dao.Query) error {
		if err := s.SysAccountRepository.UpdateByIdTx(
			ctx,
			tx,
			id,
			gormplus.Update().WithColumns(
				dao.SysAccountEntity.Username.Value(req.Username),
				dao.SysAccountEntity.Email.Value(req.Email),
				dao.SysAccountEntity.Mobile.Value(req.Mobile),
				dao.SysAccountEntity.Status.Value(req.Status),
				dao.SysAccountEntity.Avatar.Value(req.Avatar),
			).Build(),
		); err != nil {
			return err
		}
		// 先全部删除,然后创建
		if err := s.SysAccountRoleRepository.DeleteByWrapperTx(ctx, tx, func(g gormplus.IGenWrapper[dao.ISysAccountRoleEntityDo]) {
			g.Where(dao.SysAccountRoleEntity.AccountID.Eq(id))
		}, gormplus.Delete().WithPhysicalDelete().Build()); err != nil {
			return err
		}
		// 创建账号角色
		accountRoleEntity := s.buildAccountRoleEntityList(id, req.RoleIdList)
		if len(accountRoleEntity) == 0 {
			return nil
		}
		return s.SysAccountRoleRepository.CreateBatchTx(ctx, tx, accountRoleEntity)
	}); err != nil {
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
// @Param data body dto.GetSysAccountPageDTO true "分页查询参数"
// @Success 200 {array} vo.SysAccountVO
// @Router /api/v1/admin/account/page [post]
func (s SysAccount) GetSysAccountPageApi(ctx *gin.Context) {
	var req dto.GetSysAccountPageDTO
	if !s.BindAndValidateJSON(ctx, &req) {
		return
	}
	list, total, err := s.SysAccountRepository.FindPageByWrapper(ctx, req.PageNumber, req.PageSize, func(g gormplus.IGenWrapper[dao.ISysAccountEntityDo]) {
		g.
			WhereIf(req.Username != "", dao.SysAccountEntity.Username.Like("%"+req.Username+"%")).
			WhereIf(req.Email != "", dao.SysAccountEntity.Email.Like("%"+req.Email+"%")).
			WhereIf(req.Mobile != "", dao.SysAccountEntity.Mobile.Like("%"+req.Mobile+"%")).
			WhereIf(req.Status != 0, dao.SysAccountEntity.Status.Eq(req.Status))
	})
	if err != nil {
		s.Fail(ctx, err, "获取账号分页失败")
		return
	}
	s.BuildPageData(ctx, s.SysAccountMapper.EntityListToVo(list), total)
	return
}

// GetSysAccountListApi 获取账号列表
// @Summary 获取账号列表
// @Description 根据查询条件获取账号列表
// @Tags 账号中心
// @Accept json
// @Produce json
// @Param username query string false "登录账号"
// @Success 200 {array} vo.SysAccountVO "账号列表"
// @Router /api/v1/admin/account/list [get]
func (s SysAccount) GetSysAccountListApi(ctx *gin.Context) {
	username := ctx.DefaultQuery("username", "")
	list, err := s.SysAccountRepository.FindListByWrapper(ctx, func(g gormplus.IGenWrapper[dao.ISysAccountEntityDo]) {
		g.WhereIf(username != "", dao.SysAccountEntity.Username.Like("%"+username+"%"))
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
// @Success 200 {object} vo.SysAccountVO "统一响应，code=0 时 result 为 vo.SysAccountVO，code=1 时 result 为 null"
// @Router /api/v1/admin/account/{id} [get]
func (s SysAccount) GetSysAccountDetailApi(ctx *gin.Context) {
	id, ok := s.getIdParam(ctx)
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
// @Success 200 {string} string "重置密码成功"
// @Router /api/v1/admin/account/resetPassword/{id} [post]
func (s SysAccount) ResetPasswordByIdApi(ctx *gin.Context) {
	id, ok := s.getIdParam(ctx)
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
	if err = s.SysAccountRepository.UpdateById(ctx, id, gormplus.Update().WithColumns(
		dao.SysAccountEntity.Password.Value(password),
	).Build()); err != nil {
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
// @Success 200 {string} string "修改当前账号密码成功"
// @Router /api/v1/admin/account/modifyCurrentPassword [post]
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
	if err = s.SysAccountRepository.UpdateById(ctx, accountID, gormplus.Update().WithColumns(
		dao.SysAccountEntity.Password.Value(password),
	).Build()); err != nil {
		s.Fail(ctx, err, "修改当前账号密码失败")
		return
	}
	s.Success(ctx, "修改当前账号密码成功")
}

func (s SysAccount) getIdParam(ctx *gin.Context) (int64, bool) {
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
	existsAccount, err := s.SysAccountRepository.FindOneWrapper(ctx, func(g gormplus.IGenWrapper[dao.ISysAccountEntityDo]) {
		g.WhereIf(excludeID > 0, dao.SysAccountEntity.ID.Neq(excludeID)).
			WhereGroupFn(func(w gormplus.IGenWrapper[dao.ISysAccountEntityDo]) {
				w.Where(dao.SysAccountEntity.Username.Eq(username)).
					OrWhereIf(mobile != "", dao.SysAccountEntity.Mobile.Eq(mobile)).
					OrWhereIf(email != "", dao.SysAccountEntity.Email.Eq(email))
			})
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

func useQuery(db *gorm.DB) *dao.Query {
	return dao.Use(db)
}

func (s SysAccount) buildAccountRoleEntityList(accountID int64, roleIDList []int64) []*model.SysAccountRoleEntity {
	accountRoleEntity := make([]*model.SysAccountRoleEntity, 0, len(roleIDList))
	for _, roleID := range roleIDList {
		accountRoleEntity = append(accountRoleEntity, &model.SysAccountRoleEntity{
			RoleID:    roleID,
			AccountID: accountID,
		})
	}
	return accountRoleEntity
}

func NewSysAccount(baseApi *base.BaseApi) ISysAccount {
	return &SysAccount{
		BaseApi:                  baseApi,
		SysAccountRepository:     repository.NewSysAccountRepository(),
		SysAccountRoleRepository: repository.NewSysAccountRoleRepository(),
		SysAccountMapper:         mapper.NewSysAccountMapper(),
	}
}
