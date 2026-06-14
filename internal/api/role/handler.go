package role

import (
	"errors"
	"fmt"
	"gin-admin-api/internal/api/base"
	"gin-admin-api/internal/api/role/dto"
	"gin-admin-api/internal/api/role/mapper"
	"gin-admin-api/internal/api/role/vo"
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
	"strings"
)

type IRole interface {
	CreateRoleApi(ctx *gin.Context)        // 创建角色
	DeleteRoleByIdApi(ctx *gin.Context)    // 根据id删除角色
	ModifyRoleByIdApi(ctx *gin.Context)    // 根据id修改角色
	GetRolePageApi(ctx *gin.Context)       // 分页获取角色
	GetRoleListApi(ctx *gin.Context)       // 获取角色列表
	GetRoleDetailByIdApi(ctx *gin.Context) // 根据角色id获取角色详情
}

type Role struct {
	*base.BaseApi
	RoleRepository          repository.SysRoleRepository
	RoleResourcesRepository repository.SysRoleResourcesRepository
	ResourcesRepository     repository.SysResourcesRepository
	RoleMapper              mapper.ISysRoleMapper
	Enforcer                *casbin.Enforcer
}

// CreateRoleApi 创建角色
// @Summary 创建角色
// @Description 创建后台角色，并分配角色资源关系
// @Tags 角色中心
// @Accept json
// @Produce json
// @Param data body dto.CreateSysRoleDTO true "创建角色参数"
// @Success 200 {string} string "创建成功"
// @Router /api/v1/admin/role [post]
func (r Role) CreateRoleApi(ctx *gin.Context) {
	var req dto.CreateSysRoleDTO
	if !r.BindAndValidateJSON(ctx, &req) {
		return
	}
	exists, err := r.RoleRepository.Exists(ctx, gormplus.QueryOpt().Where(dao.SysRoleEntity.Name.Eq(req.Name)).Build())
	if err != nil {
		r.Fail(ctx, err, "角色名称重复校验失败")
		return
	}
	if exists {
		r.Fail(ctx, errors.New("角色名称已经存在,不能重复"), "角色名称已经存在,不能重复")
		return
	}
	var roleEntity *model.SysRoleEntity
	if err = gormplus.TransactionAsCtx(ctx, r.Db, func(db *gorm.DB) *dao.Query {
		return dao.Use(db)
	}, func(tx *dao.Query) error {
		// 1.创建角色
		roleEntity = r.RoleMapper.DtoToEntity(&req)
		if err = r.RoleRepository.CreateTx(ctx, tx, roleEntity); err != nil {
			return err
		}
		if len(req.ResourcesIdList) > 0 {
			// 2.创建角色资源
			sysRoleResourcesEntity := make([]*model.SysRoleResourcesEntity, 0)
			for _, item := range req.ResourcesIdList {
				sysRoleResourcesEntity = append(sysRoleResourcesEntity, &model.SysRoleResourcesEntity{
					RoleID:      roleEntity.ID,
					ResourcesID: item,
				})
			}
			if err = r.RoleResourcesRepository.CreateBatchTx(ctx, tx, sysRoleResourcesEntity); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		r.Fail(ctx, err, "操作失败")
		return
	}
	if err = r.syncRoleResourcesCasbin(ctx, roleEntity.ID, req.ResourcesIdList); err != nil {
		r.Fail(ctx, err, "同步角色权限失败")
		return
	}
	r.Success(ctx, "成功")
	return
}

// DeleteRoleByIdApi 删除角色
// @Summary 删除角色
// @Description 根据角色 ID 删除角色，并清理角色资源关系
// @Tags 角色中心
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {string} string "删除成功"
// @Router /api/v1/admin/role/{id} [delete]
func (r Role) DeleteRoleByIdApi(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt := cast.ToInt64(id)
	if err := gormplus.TransactionAsCtx(ctx, r.Db, func(db *gorm.DB) *dao.Query {
		return dao.Use(db)
	}, func(tx *dao.Query) error {
		// 删除中间件表
		if err := r.RoleResourcesRepository.DeleteByWrapperTx(ctx, tx, func(g gormplus.IGenWrapper[dao.ISysRoleResourcesEntityDo]) {
			g.Where(dao.SysRoleResourcesEntity.RoleID.Eq(idInt))
		}, gormplus.Delete().WithPhysicalDelete().Build()); err != nil {
			return err
		}
		// 删除角色
		return r.RoleRepository.DeleteByIdTx(ctx, tx, idInt)
	}); err != nil {
		r.Fail(ctx, err, "操作失败")
		return
	}
	if err := r.syncRoleResourcesCasbin(ctx, idInt, nil); err != nil {
		r.Fail(ctx, err, "同步角色权限失败")
		return
	}
	r.Success(ctx, "删除成功")
}

// ModifyRoleByIdApi 修改角色
// @Summary 修改角色
// @Description 根据角色 ID 修改角色基础信息和角色资源关系
// @Tags 角色中心
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param data body dto.CreateSysRoleDTO true "修改角色参数"
// @Success 200 {string} string "修改成功"
// @Router /api/v1/admin/role/{id} [put]
func (r Role) ModifyRoleByIdApi(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt := cast.ToInt64(id)
	var req dto.CreateSysRoleDTO
	if !r.BindAndValidateJSON(ctx, &req) {
		return
	}
	exists, err := r.RoleRepository.Exists(ctx, gormplus.QueryOpt().Where(
		dao.SysRoleEntity.Name.Eq(req.Name),
		dao.SysRoleEntity.ID.Neq(idInt),
	).Build())
	if err != nil {
		r.Fail(ctx, err, "角色名称重复校验失败")
		return
	}
	if exists {
		r.Fail(ctx, errors.New("角色名称已经存在,不能重复"), "角色名称已经存在,不能重复")
		return
	}
	// 修改角色，删除之前的角色资源,重新创建角色资源
	if err = gormplus.TransactionAsCtx(ctx, r.Db, func(db *gorm.DB) *dao.Query {
		return dao.Use(db)
	}, func(tx *dao.Query) error {
		// 更新角色信息
		if err = r.RoleRepository.UpdateByIdTx(ctx, tx, idInt, gormplus.Update().WithColumns(
			dao.SysRoleEntity.Name.Value(req.Name),
			dao.SysRoleEntity.Sort.Value(req.Sort),
			dao.SysRoleEntity.Status.Value(req.Status),
			dao.SysRoleEntity.Description.Value(req.Description),
		).Build()); err != nil {
			return err
		}
		// 删除之前的角色资源
		if err = r.RoleResourcesRepository.DeleteByWrapperTx(ctx, tx, func(g gormplus.IGenWrapper[dao.ISysRoleResourcesEntityDo]) {
			g.Where(dao.SysRoleResourcesEntity.RoleID.Eq(idInt))
		}, gormplus.Delete().WithPhysicalDelete().Build()); err != nil {
			return err
		}
		// 创建
		if len(req.ResourcesIdList) > 0 {
			sysRoleResourcesEntity := make([]*model.SysRoleResourcesEntity, 0)
			for _, item := range req.ResourcesIdList {
				sysRoleResourcesEntity = append(sysRoleResourcesEntity, &model.SysRoleResourcesEntity{
					RoleID:      idInt,
					ResourcesID: item,
				})
			}
			if err = r.RoleResourcesRepository.CreateBatchTx(ctx, tx, sysRoleResourcesEntity); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		r.Fail(ctx, err, "操作失败")
		return
	}
	if err = r.syncRoleResourcesCasbin(ctx, idInt, req.ResourcesIdList); err != nil {
		r.Fail(ctx, err, "同步角色权限失败")
		return
	}
	r.Success(ctx, "操作成功")
	return
}

// GetRolePageApi 分页获取角色
// @Summary 分页获取角色
// @Description 根据查询条件分页获取角色列表
// @Tags 角色中心
// @Accept json
// @Produce json
// @Param data body dto.RolePageDTO true "分页查询参数"
// @Success 200 {array} vo.SysRoleVO "角色分页列表"
// @Router /api/v1/admin/role/page [post]
func (r Role) GetRolePageApi(ctx *gin.Context) {
	var req dto.RolePageDTO
	if !r.BindAndValidateJSON(ctx, &req) {
		return
	}
	list, total, err := r.RoleRepository.FindPageByWrapper(ctx, req.PageNumber, req.PageSize, func(g gormplus.IGenWrapper[dao.ISysRoleEntityDo]) {
		g.WhereIf(req.Name != "", dao.SysRoleEntity.Name.Like("%"+req.Name+"%")).
			WhereIf(req.Status > 0, dao.SysRoleEntity.Status.Eq(req.Status))
	}, gormplus.QueryOpt().Order(dao.SysRoleEntity.Sort).Build())
	if err != nil {
		r.Fail(ctx, err, "获取分页数据失败")
		return
	}
	r.BuildPageData(ctx, k.Map(list, func(item *model.SysRoleEntity, index int) vo.SysRoleVO {
		return *r.RoleMapper.EntityToVo(item)
	}), total)
	return
}

// GetRoleListApi 获取角色列表
// @Summary 获取角色列表
// @Description 获取正常状态的角色列表
// @Tags 角色中心
// @Accept json
// @Produce json
// @Success 200 {array} vo.SysRoleVO "角色列表"
// @Router /api/v1/admin/role/list [get]
func (r Role) GetRoleListApi(ctx *gin.Context) {
	list, err := r.RoleRepository.FindList(ctx, gormplus.QueryOpt().Where(
		dao.SysRoleEntity.Status.Eq(enum.StatusNormalEnum),
	).Order(dao.SysRoleEntity.Sort).Build())
	if err != nil {
		r.Success(ctx, make([]interface{}, 0))
		return
	}
	r.Success(ctx, k.Map(list, func(item *model.SysRoleEntity, index int) vo.SysRoleVO {
		return *r.RoleMapper.EntityToVo(item)
	}))
}

// GetRoleDetailByIdApi 获取角色详情
// @Summary 获取角色详情
// @Description 根据角色 ID 获取角色详情和已授权资源 ID 列表
// @Tags 角色中心
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} vo.SysRoleDetailVO "角色详情"
// @Router /api/v1/admin/role/detail/{id} [get]
func (r Role) GetRoleDetailByIdApi(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt := cast.ToInt64(id)
	roleEntity, err := r.RoleRepository.FindById(ctx, idInt)
	if err != nil {
		r.Fail(ctx, err, "传递的角色id错误")
		return
	}
	// 根据角色获取授权的资源
	list, err := r.RoleResourcesRepository.FindList(ctx, gormplus.QueryOpt().Where(dao.SysRoleResourcesEntity.RoleID.Eq(idInt)).Build())
	resourcesIdList := make([]int64, 0)
	if err == nil && len(list) > 0 {
		resourcesIdList = k.Map(list, func(item *model.SysRoleResourcesEntity, index int) int64 {
			return item.ResourcesID
		})
	}
	roleVO := r.RoleMapper.EntityToVo(roleEntity)
	r.Success(ctx, vo.SysRoleDetailVO{
		SysRoleVO:       *roleVO,
		ResourcesIdList: resourcesIdList,
	})
}

// syncRoleResourcesCasbin 同步角色资源到 Casbin 的 p 策略。
// 会先清空 role_{roleId} 的旧资源权限，再按 resourcesIdList 写入新的接口权限；
// resourcesIdList 为空时表示仅清空该角色权限，常用于删除角色。
func (r Role) syncRoleResourcesCasbin(ctx *gin.Context, roleId int64, resourcesIdList []int64) error {
	if r.Enforcer == nil {
		return errors.New("casbin enforcer未初始化")
	}
	sub := fmt.Sprintf("role_%d", roleId)

	// 清除该角色所有旧的 p 策略
	if _, err := r.Enforcer.RemoveFilteredPolicy(0, sub); err != nil {
		return err
	}

	if len(resourcesIdList) == 0 {
		return nil
	}
	resourceIds := k.Filter(k.Distinct(resourcesIdList), func(item int64, index int) bool {
		return item > 0
	})
	if len(resourceIds) == 0 {
		return nil
	}

	resourcesEntities, err := r.ResourcesRepository.FindList(ctx, gormplus.QueryOpt().Where(
		dao.SysResourcesEntity.ID.In(resourceIds...),
		dao.SysResourcesEntity.ResourcesType.Eq(3),
		dao.SysResourcesEntity.Status.Eq(enum.StatusNormalEnum),
	).Build())
	if err != nil {
		return fmt.Errorf("查询角色资源失败: %w", err)
	}

	rules := make([][]string, 0, len(resourcesEntities))
	for _, res := range resourcesEntities {
		url := strings.TrimSpace(res.URL)
		method := strings.ToUpper(strings.TrimSpace(res.Method))
		if url == "" || method == "" {
			continue
		}
		rules = append(rules, []string{sub, url, method})
	}

	if len(rules) > 0 {
		if _, err = r.Enforcer.AddPolicies(rules); err != nil {
			return err
		}
	}
	return nil
}

func NewRole(baseApi *base.BaseApi, enforcer *casbin.Enforcer) IRole {
	return Role{
		BaseApi:                 baseApi,
		RoleRepository:          repository.NewSysRoleRepository(),
		RoleResourcesRepository: repository.NewSysRoleResourcesRepository(),
		ResourcesRepository:     repository.NewSysResourcesRepository(),
		RoleMapper:              mapper.NewSysRoleMapper(),
		Enforcer:                enforcer,
	}
}
