package account

import (
	"errors"
	"fmt"
	"gin-admin-api/internal/api/account/dto"
	"gin-admin-api/internal/api/account/mapper"
	"gin-admin-api/internal/api/account/vo"
	_ "gin-admin-api/internal/api/account/vo"
	"gin-admin-api/internal/api/base"
	"gin-admin-api/internal/dal/dao"
	"gin-admin-api/internal/dal/model"
	"gin-admin-api/internal/dal/repository"
	"gin-admin-api/pkg/enum"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/kuangshp/go-utils/k"
	"github.com/kuangshp/gorm-plus"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type ISysAccount interface {
	CreateSysAccountApi(ctx *gin.Context)                // 创建账号
	DeleteSysAccountByIdApi(ctx *gin.Context)            // 根据id删除账号
	ModifySysAccountByIdApi(ctx *gin.Context)            // 根据id修改账号
	GetSysAccountPageApi(ctx *gin.Context)               // 分页获取账号
	GetSysAccountListApi(ctx *gin.Context)               // 获取账号列表
	GetSysAccountDetailApi(ctx *gin.Context)             // 根据id获取账号详情
	GetCurrentSysAccountInfoApi(ctx *gin.Context)        // 获取当前登录账号信息
	ResetPasswordByIdApi(ctx *gin.Context)               // 根据id重置账号密码
	ModifyCurrentSysAccountPasswordApi(ctx *gin.Context) // 修改当前登录账号密码
}

type SysAccount struct {
	*base.BaseApi
	SysAccountRepository     repository.SysAccountRepository
	SysAccountRoleRepository repository.SysAccountRoleRepository
	SysAccountPostRepository repository.SysAccountPostRepository
	SysDeptRepository        repository.SysDeptRepository
	SysPostRepository        repository.SysPostRepository
	RoleResourcesRepository  repository.SysRoleResourcesRepository
	ResourcesRepository      repository.SysResourcesRepository
	SysAccountMapper         mapper.ISysAccountMapper
	Enforcer                 *casbin.Enforcer
}

// CreateSysAccountApi 创建账号
// @Summary 创建账号
// @Description 创建后台账号，并分配角色、岗位关系；deptId 为账号所属部门，postIdList 第一个岗位为主岗
// @Tags 账号中心
// @Accept json
// @Produce json
// @Param data body dto.CreateSysAccountDTO true "创建账号参数"
// @Success 200 {string} string "创建成功"
// @Router /api/v1/admin/account [post]
func (s SysAccount) CreateSysAccountApi(ctx *gin.Context) {
	var req dto.CreateSysAccountDTO
	if !s.BindAndValidateJSON(ctx, &req) {
		return
	}
	// 判断是否已经存在
	if !s.checkAccountUnique(ctx, req.Username, req.Mobile, req.Email, 0) {
		return
	}
	// 判断部门
	if !s.checkDeptValid(ctx, req.DeptID) {
		return
	}
	// 判断岗位
	postIDs, ok := s.buildAccountPostScope(ctx, req.PostIdList)
	if !ok {
		return
	}
	password, err := k.MakePassword(req.Password)
	if err != nil {
		s.Logger.Error("密码加密失败")
		s.Fail(ctx, err, "创建失败")
		return
	}

	var accountEntity *model.SysAccountEntity
	if err = gormplus.TransactionAsCtx(ctx, s.Db, useQuery, func(tx *dao.Query) error {
		accountEntity = s.SysAccountMapper.DtoToEntity(&req, password, enum.StatusNormalEnum)
		if err = s.SysAccountRepository.CreateTx(ctx, tx, accountEntity, gormplus.Create().WithOmit(
			dao.SysAccountEntity.LastLoginIP,
			dao.SysAccountEntity.LastLoginDate,
		).Build()); err != nil {
			s.Logger.Error("创建失败")
			return err
		}

		// 分配岗位，第一个岗位为主岗
		accountPostEntity := s.buildAccountPostEntityList(accountEntity.ID, postIDs)
		if err = s.SysAccountPostRepository.CreateBatchTx(ctx, tx, accountPostEntity); err != nil {
			s.Logger.Error("创建账号岗位失败")
			return err
		}
		// 分配角色
		accountRoleEntity := s.buildAccountRoleEntityList(accountEntity.ID, req.RoleIdList)
		if err = s.SysAccountRoleRepository.CreateBatchTx(ctx, tx, accountRoleEntity); err != nil {
			s.Logger.Error("创建账号角色失败")
			return err
		}
		return nil
	}); err != nil {
		s.Fail(ctx, err, "创建失败")
		return
	}
	// 数据同步到casbin中
	if err = s.syncAccountRolesCasbin(accountEntity.ID, req.RoleIdList); err != nil {
		s.Fail(ctx, err, "同步账号角色权限失败")
		return
	}
	s.Success(ctx, "创建成功")
}

// DeleteSysAccountByIdApi 删除账号
// @Summary 删除账号
// @Description 根据账号 ID 删除账号，并清理账号角色、岗位关系
// @Tags 账号中心
// @Accept json
// @Produce json
// @Param id path int true "账号ID"
// @Success 200 {string} string "删除成功"
// @Router /api/v1/admin/account/{id} [delete]
func (s SysAccount) DeleteSysAccountByIdApi(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id := cast.ToInt64(idStr)
	if _, err := s.SysAccountRepository.FindById(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Fail(ctx, err, "账号不存在")
			return
		}
		s.Fail(ctx, err, "删除失败")
		return
	}
	if err := gormplus.TransactionAsCtx(ctx, s.Db, useQuery, func(tx *dao.Query) error {
		// 删除账号和角色中间表
		if err := s.SysAccountRoleRepository.DeleteByWrapperTx(ctx, tx, func(g gormplus.IGenWrapper[dao.ISysAccountRoleEntityDo]) {
			g.Where(dao.SysAccountRoleEntity.AccountID.Eq(id))
		}, gormplus.Delete().WithPhysicalDelete().Build()); err != nil {
			return err
		}
		// 删除账号和岗位中间表
		if err := s.SysAccountPostRepository.DeleteByWrapperTx(ctx, tx, func(g gormplus.IGenWrapper[dao.ISysAccountPostEntityDo]) {
			g.Where(dao.SysAccountPostEntity.AccountID.Eq(id))
		}, gormplus.Delete().WithPhysicalDelete().Build()); err != nil {
			return err
		}
		// 删除账号
		return s.SysAccountRepository.DeleteByIdTx(ctx, tx, id)
	}); err != nil {
		s.Fail(ctx, err, "删除失败")
		return
	}
	if err := s.syncAccountRolesCasbin(id, nil); err != nil {
		s.Fail(ctx, err, "同步账号角色权限失败")
		return
	}
	s.Success(ctx, "删除成功")
}

// ModifySysAccountByIdApi 修改账号
// @Summary 修改账号
// @Description 根据账号 ID 修改账号基础信息、角色关系和岗位关系；deptId 为账号所属部门，postIdList 第一个岗位为主岗
// @Tags 账号中心
// @Accept json
// @Produce json
// @Param id path int true "账号ID"
// @Param data body dto.ModifySysAccountDTO true "修改账号参数"
// @Success 200 {string} string "修改成功"
// @Router /api/v1/admin/account/{id} [put]
func (s SysAccount) ModifySysAccountByIdApi(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id := cast.ToInt64(idStr)
	var req dto.ModifySysAccountDTO
	if !s.BindAndValidateJSON(ctx, &req) {
		return
	}
	if _, err := s.SysAccountRepository.FindById(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Fail(ctx, err, "账号不存在")
			return
		}
		s.Fail(ctx, err, "修改失败")
		return
	}
	// 判断手机号码、用户名、email是否唯一
	if !s.checkAccountUnique(ctx, req.Username, req.Mobile, req.Email, id) {
		return
	}
	// 判断部门
	if !s.checkDeptValid(ctx, req.DeptID) {
		return
	}
	// 判断岗位
	postIDs, ok := s.buildAccountPostScope(ctx, req.PostIdList)
	if !ok {
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
				dao.SysAccountEntity.DeptID.Value(req.DeptID),
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
		// 账号岗位先删除然后创建
		if err := s.SysAccountPostRepository.DeleteByWrapperTx(ctx, tx, func(g gormplus.IGenWrapper[dao.ISysAccountPostEntityDo]) {
			g.Where(dao.SysAccountPostEntity.AccountID.Eq(id))
		}, gormplus.Delete().WithPhysicalDelete().Build()); err != nil {
			return err
		}
		accountPostEntity := s.buildAccountPostEntityList(id, postIDs)
		if err := s.SysAccountPostRepository.CreateBatchTx(ctx, tx, accountPostEntity); err != nil {
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
	if err := s.syncAccountRolesCasbin(id, req.RoleIdList); err != nil {
		s.Fail(ctx, err, "同步账号角色权限失败")
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
			WhereOrGroupIf(req.Keyword != "", dao.SysAccountEntity.Username.Like("%"+req.Keyword+"%"), dao.SysAccountEntity.Email.Like("%"+req.Keyword+"%"), dao.SysAccountEntity.Mobile.Like("%"+req.Keyword+"%")).
			WhereIf(req.Status != 0, dao.SysAccountEntity.Status.Eq(req.Status)).
			WhereIf(req.DeptID > 0, dao.SysAccountEntity.DeptID.Eq(req.DeptID))
	})
	if err != nil {
		s.Fail(ctx, err, "获取账号分页失败")
		return
	}
	result := s.SysAccountMapper.EntityListToVo(list)
	if !s.fillAccountPostList(ctx, result) {
		return
	}
	s.BuildPageData(ctx, result, total)
	return
}

// GetSysAccountListApi 获取账号列表
// @Summary 获取账号列表
// @Description 根据查询条件获取账号列表
// @Tags 账号中心
// @Accept json
// @Produce json
// @Param username query string false "登录账号"
// @Param deptId query int false "所属部门ID"
// @Param postId query int false "岗位ID"
// @Success 200 {array} vo.SysAccountVO "账号列表"
// @Router /api/v1/admin/account/list [get]
func (s SysAccount) GetSysAccountListApi(ctx *gin.Context) {
	username := ctx.DefaultQuery("username", "")
	deptID := cast.ToInt64(ctx.DefaultQuery("deptId", "0"))
	postID := cast.ToInt64(ctx.DefaultQuery("postId", "0"))
	accountIDs, ok := s.getAccountIDListByPost(ctx, postID)
	if !ok {
		return
	}
	list, err := s.SysAccountRepository.FindListByWrapper(ctx, func(g gormplus.IGenWrapper[dao.ISysAccountEntityDo]) {
		g.WhereIf(username != "", dao.SysAccountEntity.Username.Like("%"+username+"%")).
			WhereIf(deptID > 0, dao.SysAccountEntity.DeptID.Eq(deptID))
		if postID > 0 {
			if len(accountIDs) == 0 {
				g.Where(dao.SysAccountEntity.ID.Eq(-1))
			} else {
				g.Where(dao.SysAccountEntity.ID.In(accountIDs...))
			}
		}
	})
	if err != nil {
		s.Fail(ctx, err, "获取账号列表失败")
		return
	}
	result := s.SysAccountMapper.EntityListToVo(list)
	// 填充账号岗位
	if !s.fillAccountPostList(ctx, result) {
		return
	}
	s.Success(ctx, result)
}

// GetSysAccountDetailApi 获取账号详情
// @Summary 获取账号详情
// @Description 根据账号 ID 获取账号详情
// @Tags 账号中心
// @Accept json
// @Produce json
// @Param id path int true "账号ID"
// @Success 200 {object} vo.SysAccountDetailVO "统一响应，code=0 时 result 为 vo.SysAccountDetailVO，code=1 时 result 为 null"
// @Router /api/v1/admin/account/detail/{id} [get]
func (s SysAccount) GetSysAccountDetailApi(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id := cast.ToInt64(idStr)
	accountEntity, err := s.SysAccountRepository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Fail(ctx, err, "账号不存在")
			return
		}
		s.Fail(ctx, err, "获取账号详情失败")
		return
	}
	accountVO := s.SysAccountMapper.EntityToVo(accountEntity)
	// 查询账号角色
	list, err := s.SysAccountRoleRepository.FindList(ctx, gormplus.QueryOpt().Where(dao.SysAccountRoleEntity.AccountID.Eq(id)).Build())
	var roleIdList = make([]int64, 0)
	if err != nil {
		s.Fail(ctx, err, "查询账号角色失败")
		return
	}
	if len(list) > 0 {
		roleIdList = k.Map(list, func(item *model.SysAccountRoleEntity, index int) int64 {
			return item.RoleID
		})
	}
	// 获取岗位
	accountVO.PostIdList, err = s.getPostIDListByAccount(ctx, id)
	if err != nil {
		s.Fail(ctx, err, "查询账号岗位失败")
		return
	}
	// 获取岗位名称
	accountVO.PostList, err = s.getPostListByAccount(ctx, id)
	if err != nil {
		s.Fail(ctx, err, "查询账号岗位名称失败")
		return
	}
	s.Success(ctx, vo.SysAccountDetailVO{
		SysAccountVO: *accountVO,
		RoleIdList:   roleIdList,
	})
	return
}

// GetCurrentSysAccountInfoApi 获取当前登录账号信息
// @Summary 获取当前登录账号信息
// @Description 获取当前登录账号基础信息、授权角色和授权接口权限
// @Tags 账号中心
// @Accept json
// @Produce json
// @Success 200 {object} vo.SysAccountCurrentInfoVO "统一响应，code=0 时 result 为 vo.SysAccountCurrentInfoVO，code=1 时 result 为 null"
// @Router /api/v1/admin/account/info [get]
func (s SysAccount) GetCurrentSysAccountInfoApi(ctx *gin.Context) {
	accountID, ok := s.getCurrentAccountID(ctx)
	if !ok {
		return
	}
	accountEntity, err := s.SysAccountRepository.FindById(gormplus.SkipDataPermission(ctx.Request.Context()), accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Fail(ctx, err, "账号不存在")
			return
		}
		s.Fail(ctx, err, "获取当前账号信息失败")
		return
	}
	// 获取资源id
	resourcesIdList := s.RoleResourcesRepository.GetResourcesByAccountId(ctx, accountID)
	var apiPermissionList = make([]vo.SysAccountApiPermissionVO, 0)
	if len(resourcesIdList) > 0 {
		resourcesEntity, err := s.ResourcesRepository.FindByIdList(ctx, resourcesIdList, gormplus.QueryOpt().Where(
			dao.SysResourcesEntity.ResourcesType.Eq(enum.ResourcesTypeApiEnum),
		).Build())
		if err == nil && len(resourcesEntity) > 0 {
			apiPermissionList = k.Map(resourcesEntity, func(item *model.SysResourcesEntity, index int) vo.SysAccountApiPermissionVO {
				return vo.SysAccountApiPermissionVO{
					ID:       item.ID,       // 资源id
					Title:    item.Title,    // 接口名称
					URL:      item.URL,      // 接口地址
					Method:   item.Method,   // 请求方式
					ParentID: item.ParentID, // 上一级id
					Sort:     item.Sort,     // 排序
				}
			})
		}
	}

	s.Success(ctx, vo.SysAccountCurrentInfoVO{
		ID:                accountEntity.ID,
		Username:          accountEntity.Username,
		Email:             accountEntity.Email,
		Mobile:            accountEntity.Mobile,
		LastLoginDate:     accountEntity.LastLoginDate.Unix(),
		LastLoginIP:       accountEntity.LastLoginIP,
		Status:            accountEntity.Status,
		Avatar:            accountEntity.Avatar,
		IsAdmin:           accountEntity.IsAdmin,
		DeptID:            accountEntity.DeptID,
		CreatedAt:         accountEntity.CreatedAt.Unix(),
		UpdatedAt:         accountEntity.UpdatedAt.Unix(),
		ApiPermissionList: apiPermissionList,
	})
}

// ResetPasswordByIdApi 重置账号密码
// @Summary 重置账号密码
// @Description 根据账号 ID 重置账号密码
// @Tags 账号中心
// @Accept json
// @Produce json
// @Param data body dto.ResetPasswordDTO true "重置密码参数"
// @Success 200 {string} string "重置密码成功"
// @Router /api/v1/admin/account/resetPassword [post]
func (s SysAccount) ResetPasswordByIdApi(ctx *gin.Context) {
	var req dto.ResetPasswordDTO
	if !s.BindAndValidateJSON(ctx, &req) {
		return
	}
	if req.Password != req.ConfirmPassword {
		s.Fail(ctx, errors.New("两次密码不一致"), "两次密码不一致")
		return
	}
	if _, err := s.SysAccountRepository.FindById(ctx, req.Id); err != nil {
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
	if err = s.SysAccountRepository.UpdateById(ctx, req.Id, gormplus.Update().WithColumns(
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

// 获取当前登录用户信息
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

// 校验用户名、手机号码、邮箱唯一
func (s SysAccount) checkAccountUnique(ctx *gin.Context, username, mobile, email string, excludeID int64) bool {
	queryCtx := gormplus.SkipDataPermission(ctx.Request.Context())
	existsAccount, err := s.SysAccountRepository.FindOneWrapper(queryCtx, func(g gormplus.IGenWrapper[dao.ISysAccountEntityDo]) {
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

// 校验授权的部门
func (s SysAccount) checkDeptValid(ctx *gin.Context, deptID int64) bool {
	if deptID <= 0 {
		s.Fail(ctx, errors.New("所属部门不能为空"), "所属部门不能为空")
		return false
	}
	exists, err := s.SysDeptRepository.Exists(gormplus.SkipDataPermission(ctx.Request.Context()), gormplus.QueryOpt().Where(
		dao.SysDeptEntity.ID.Eq(deptID),
		dao.SysDeptEntity.Status.Eq(enum.StatusNormalEnum),
	).Build())
	if err != nil {
		s.Fail(ctx, err, "所属部门校验失败")
		return false
	}
	if !exists {
		s.Fail(ctx, errors.New("所属部门不存在或已禁用"), "所属部门不存在或已禁用")
		return false
	}
	return true
}

func (s SysAccount) getAccountIDListByPost(ctx *gin.Context, postID int64) ([]int64, bool) {
	if postID <= 0 {
		return nil, true
	}
	list, err := s.SysAccountPostRepository.FindList(ctx, gormplus.QueryOpt().Where(
		dao.SysAccountPostEntity.PostID.Eq(postID),
	).Build())
	if err != nil {
		s.Fail(ctx, err, "查询岗位账号失败")
		return nil, false
	}
	accountIDs := make([]int64, 0, len(list))
	seen := make(map[int64]struct{}, len(list))
	for _, item := range list {
		if item.AccountID <= 0 {
			continue
		}
		if _, ok := seen[item.AccountID]; ok {
			continue
		}
		seen[item.AccountID] = struct{}{}
		accountIDs = append(accountIDs, item.AccountID)
	}
	return accountIDs, true
}

func (s SysAccount) getPostIDListByAccount(ctx *gin.Context, accountID int64) ([]int64, error) {
	accountPostList, err := s.SysAccountPostRepository.FindList(ctx, gormplus.QueryOpt().
		Where(dao.SysAccountPostEntity.AccountID.Eq(accountID)).
		Order(dao.SysAccountPostEntity.IsPrimary.Desc(), dao.SysAccountPostEntity.ID).
		Build())
	if err != nil {
		return nil, err
	}
	return k.Map(accountPostList, func(item *model.SysAccountPostEntity, index int) int64 {
		return item.PostID
	}), nil
}

func (s SysAccount) getPostListByAccount(ctx *gin.Context, accountID int64) ([]vo.PostVO, error) {
	postIDList, err := s.getPostIDListByAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if len(postIDList) == 0 {
		return []vo.PostVO{}, nil
	}
	postEntityList, err := s.SysPostRepository.FindByIdList(ctx, postIDList)
	if err != nil {
		return nil, err
	}
	postMap := make(map[int64]string, len(postEntityList))
	for _, post := range postEntityList {
		postMap[post.ID] = post.Name
	}
	result := make([]vo.PostVO, 0, len(postIDList))
	for _, postID := range postIDList {
		result = append(result, vo.PostVO{
			ID:   postID,
			Name: postMap[postID],
		})
	}
	return result, nil
}

func (s SysAccount) fillAccountPostList(ctx *gin.Context, list []*vo.SysAccountVO) bool {
	if len(list) == 0 {
		return true
	}
	accountIDs := make([]int64, 0, len(list))
	accountMap := make(map[int64]*vo.SysAccountVO, len(list))
	for _, item := range list {
		if item == nil || item.ID <= 0 {
			continue
		}
		accountIDs = append(accountIDs, item.ID)
		accountMap[item.ID] = item
		item.PostIdList = make([]int64, 0)
		item.PostList = make([]vo.PostVO, 0)
	}
	if len(accountIDs) == 0 {
		return true
	}
	accountPostList, err := s.SysAccountPostRepository.FindList(ctx, gormplus.QueryOpt().
		Where(dao.SysAccountPostEntity.AccountID.In(accountIDs...)).
		Order(dao.SysAccountPostEntity.AccountID, dao.SysAccountPostEntity.IsPrimary.Desc(), dao.SysAccountPostEntity.ID).
		Build())
	if err != nil {
		s.Fail(ctx, err, "查询账号岗位失败")
		return false
	}
	// 收集所有岗位ID，批量查询岗位名称
	postIDSet := make(map[int64]struct{})
	for _, item := range accountPostList {
		postIDSet[item.PostID] = struct{}{}
	}
	postIDs := make([]int64, 0, len(postIDSet))
	for id := range postIDSet {
		postIDs = append(postIDs, id)
	}
	postMap := make(map[int64]string, len(postIDs))
	if len(postIDs) > 0 {
		postEntityList, err := s.SysPostRepository.FindByIdList(ctx, postIDs)
		if err != nil {
			s.Fail(ctx, err, "查询岗位信息失败")
			return false
		}
		for _, post := range postEntityList {
			postMap[post.ID] = post.Name
		}
	}
	for _, item := range accountPostList {
		accountVO := accountMap[item.AccountID]
		if accountVO == nil {
			continue
		}
		accountVO.PostIdList = append(accountVO.PostIdList, item.PostID)
		accountVO.PostList = append(accountVO.PostList, vo.PostVO{
			ID:   item.PostID,
			Name: postMap[item.PostID],
		})
	}
	return true
}

// buildAccountPostScope 校验岗位并返回去重后的岗位列表。列表第一个岗位为主岗。
func (s SysAccount) buildAccountPostScope(ctx *gin.Context, postIDList []int64) ([]int64, bool) {
	postIDs := make([]int64, 0, len(postIDList))
	seen := make(map[int64]struct{}, len(postIDList))
	for _, postID := range postIDList {
		if postID <= 0 {
			continue
		}
		if _, ok := seen[postID]; ok {
			continue
		}
		seen[postID] = struct{}{}
		postIDs = append(postIDs, postID)
	}
	if len(postIDs) == 0 {
		s.Fail(ctx, errors.New("授权岗位不能为空"), "授权岗位不能为空")
		return nil, false
	}
	postList, err := s.SysPostRepository.FindByIdList(ctx, postIDs, gormplus.QueryOpt().Where(
		dao.SysPostEntity.Status.Eq(enum.StatusNormalEnum),
	).Build())
	if err != nil {
		s.Fail(ctx, err, "授权岗位校验失败")
		return nil, false
	}
	postMap := make(map[int64]*model.SysPostEntity, len(postList))
	for _, post := range postList {
		postMap[post.ID] = post
	}
	if len(postMap) != len(postIDs) {
		s.Fail(ctx, errors.New("授权岗位不存在、已禁用或无权限"), "授权岗位不存在、已禁用或无权限")
		return nil, false
	}
	return postIDs, true
}

func useQuery(db *gorm.DB) *dao.Query {
	return dao.Use(db)
}

// 组装账号和角色的数据
func (s SysAccount) buildAccountRoleEntityList(accountID int64, roleIDList []int64) []*model.SysAccountRoleEntity {
	roleIDs := k.Filter(k.Distinct(roleIDList), func(item int64, index int) bool {
		return item > 0
	})
	accountRoleEntity := make([]*model.SysAccountRoleEntity, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		accountRoleEntity = append(accountRoleEntity, &model.SysAccountRoleEntity{
			RoleID:    roleID,
			AccountID: accountID,
		})
	}
	return accountRoleEntity
}

// 组装账号和岗位的数据，第一个岗位为主岗。
func (s SysAccount) buildAccountPostEntityList(accountID int64, postIDList []int64) []*model.SysAccountPostEntity {
	accountPostEntity := make([]*model.SysAccountPostEntity, 0, len(postIDList))
	for index, postID := range postIDList {
		isPrimary := int64(1)
		if index == 0 {
			isPrimary = 2
		}
		accountPostEntity = append(accountPostEntity, &model.SysAccountPostEntity{
			AccountID: accountID,
			PostID:    postID,
			IsPrimary: isPrimary,
		})
	}
	return accountPostEntity
}

// syncAccountRolesCasbin 同步账号角色到 Casbin 的 g 分组策略。
// 会先清空 user_{accountId} 的旧角色关系，再按 roleIdList 写入新的角色分组；
// roleIdList 为空时表示仅清空该账号角色关系，常用于删除账号。
func (s SysAccount) syncAccountRolesCasbin(accountId int64, roleIdList []int64) error {
	if s.Enforcer == nil {
		return errors.New("casbin enforcer未初始化")
	}
	sub := fmt.Sprintf("user_%d", accountId)
	if _, err := s.Enforcer.DeleteRolesForUser(sub); err != nil {
		return err
	}
	roleIds := k.Filter(k.Distinct(roleIdList), func(item int64, index int) bool {
		return item > 0
	})
	if len(roleIds) == 0 {
		return nil
	}
	rules := k.Map(roleIds, func(roleId int64, index int) []string {
		return []string{sub, fmt.Sprintf("role_%d", roleId)}
	})
	if len(rules) > 0 {
		if _, err := s.Enforcer.AddGroupingPolicies(rules); err != nil {
			return err
		}
	}
	return nil
}

func NewSysAccount(baseApi *base.BaseApi, enforcer *casbin.Enforcer) ISysAccount {
	return &SysAccount{
		BaseApi:                  baseApi,
		SysAccountRepository:     repository.NewSysAccountRepository(),
		SysAccountRoleRepository: repository.NewSysAccountRoleRepository(),
		SysAccountPostRepository: repository.NewSysAccountPostRepository(),
		SysDeptRepository:        repository.NewSysDeptRepository(),
		SysPostRepository:        repository.NewSysPostRepository(),
		RoleResourcesRepository:  repository.NewSysRoleResourcesRepository(),
		ResourcesRepository:      repository.NewSysResourcesRepository(),
		SysAccountMapper:         mapper.NewSysAccountMapper(),
		Enforcer:                 enforcer,
	}
}
