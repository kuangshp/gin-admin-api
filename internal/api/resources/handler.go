package resources

import (
	"errors"
	"fmt"
	"gin-admin-api/internal/api/base"
	"gin-admin-api/internal/api/resources/dto"
	"gin-admin-api/internal/api/resources/mapper"
	"gin-admin-api/internal/api/resources/vo"
	"gin-admin-api/internal/dal/dao"
	"gin-admin-api/internal/dal/model"
	"gin-admin-api/internal/dal/repository"
	"gin-admin-api/pkg/enum"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/kuangshp/go-utils/k"
	"github.com/kuangshp/gorm-plus"
	"github.com/spf13/cast"
	"strings"
)

type IResources interface {
	CreateResourcesApi(ctx *gin.Context)      // 创建资源
	DeleteResourcesByIdApi(ctx *gin.Context)  // 根据主键删除资源
	ModifyResourcesByIdApi(ctx *gin.Context)  // 根据主键修改资源
	GetResourcesTreePageApi(ctx *gin.Context) // 获取资源树
	GetResourcesCatalogApi(ctx *gin.Context)  //  查询目录或者目录和菜单
	GetResourcesListApi(ctx *gin.Context)     // 获取全部的目录及菜单
	GetDetailByIdApi(ctx *gin.Context)        // 根据id获取详情
}

type Resources struct {
	*base.BaseApi
	ResourcesRepository     repository.SysResourcesRepository
	RoleResourcesRepository repository.SysRoleResourcesRepository
	ResourcesMapper         mapper.ISysResourcesMapper
	Enforcer                *casbin.Enforcer
}

// CreateResourcesApi 创建资源
// @Summary 创建资源
// @Description 创建目录、菜单或接口资源
// @Tags 资源中心
// @Accept json
// @Produce json
// @Param data body dto.CreateSysResourcesDTO true "创建资源参数"
// @Success 200 {string} string "操作成功"
// @Router /api/v1/admin/resources [post]
func (r Resources) CreateResourcesApi(ctx *gin.Context) {
	var req dto.CreateSysResourcesDTO
	if !r.BindAndValidateJSON(ctx, &req) {
		return
	}
	exists, err := r.ResourcesRepository.Exists(ctx, gormplus.QueryOpt().Where(dao.SysResourcesEntity.URL.Eq(req.URL)).Build())
	if err != nil || exists {
		r.Fail(ctx, err, "url已经存在")
		return
	}
	if err = r.ResourcesRepository.Create(ctx, &model.SysResourcesEntity{
		Title:         req.Title,             // 名称:按钮标题,或菜单标题
		URL:           req.URL,               // 按钮请求url,或菜单路由
		Method:        req.Method,            // 接口的请求方式
		Icon:          req.Icon,              // 菜单小图标
		ResourcesType: req.ResourcesType,     // 类型:1表示目录,2表示菜单,3表示接口
		IsCache:       req.IsCache,           // 是否缓存:1表示缓存:2不缓存
		IsHidden:      req.IsHidden,          // 是否隐藏:1表示不隐藏,2表示隐藏
		IsLink:        req.IsLink,            // 是否为外部链接:1表示不是,2表示是
		ParentID:      req.ParentID,          // 上一级id
		Sort:          req.Sort,              // 菜单,或按钮排序
		Description:   req.Description,       // 描述
		Status:        enum.StatusNormalEnum, // 状态:1表示正常,2表示禁用
	}); err != nil {
		r.Fail(ctx, err, "创建资源失败")
		return
	}
	r.Success(ctx, "操作成功")
	return
}

// DeleteResourcesByIdApi 删除资源
// @Summary 删除资源
// @Description 根据资源 ID 删除资源；存在子资源或已被角色绑定时不能删除
// @Tags 资源中心
// @Accept json
// @Produce json
// @Param id path int true "资源ID"
// @Success 200 {string} string "操作成功"
// @Router /api/v1/admin/resources/{id} [delete]
func (r Resources) DeleteResourcesByIdApi(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt := cast.ToInt64(id)
	resourceEntity, err := r.ResourcesRepository.FindById(ctx, idInt)
	if err != nil {
		r.Fail(ctx, err, "传递的主键id错误")
		return
	}
	exists, err := r.ResourcesRepository.Exists(ctx, gormplus.QueryOpt().Where(dao.SysResourcesEntity.ParentID.Eq(idInt)).Build())
	if exists || err != nil {
		r.Fail(ctx, err, "资源下存在子资源，无法删除")
		return
	}
	exist, err := r.RoleResourcesRepository.Exists(ctx, gormplus.QueryOpt().Where(dao.SysRoleResourcesEntity.ResourcesID.Eq(idInt)).Build())
	if err != nil {
		r.Fail(ctx, err, "根据资源id查询角色资源失败")
		return
	}
	if exist {
		r.Fail(ctx, err, "资源已经被角色绑定，不能直接删除")
		return
	}
	if err = r.ResourcesRepository.DeleteById(ctx, idInt); err != nil {
		r.Fail(ctx, err, "操作失败")
	}
	// 兜底清理 casbin 中残留的该资源策略
	if resourceEntity.ResourcesType == enum.ResourcesTypeApiEnum && resourceEntity.URL != "" && resourceEntity.Method != "" {
		_, _ = r.Enforcer.RemoveFilteredPolicy(1, resourceEntity.URL, resourceEntity.Method)
	}
	r.Success(ctx, "操作成功")
	return
}

// ModifyResourcesByIdApi 修改资源
// @Summary 修改资源
// @Description 根据资源 ID 修改目录、菜单或接口资源
// @Tags 资源中心
// @Accept json
// @Produce json
// @Param id path int true "资源ID"
// @Param data body dto.CreateSysResourcesDTO true "修改资源参数"
// @Success 200 {string} string "操作成功"
// @Router /api/v1/admin/resources/{id} [put]
func (r Resources) ModifyResourcesByIdApi(ctx *gin.Context) {
	var req dto.CreateSysResourcesDTO
	if !r.BindAndValidateJSON(ctx, &req) {
		return
	}
	id := ctx.Param("id")
	idInt := cast.ToInt64(id)
	exists, err := r.ResourcesRepository.Exists(ctx, gormplus.QueryOpt().Where(
		dao.SysResourcesEntity.URL.Eq(req.URL),
		dao.SysResourcesEntity.ID.Neq(idInt),
	).Build())
	if err != nil || exists {
		r.Fail(ctx, err, "url已经存在")
		return
	}
	// 判断url或者唯一标识是否存在
	if err = r.ResourcesRepository.UpdateById(ctx, idInt,
		gormplus.Update().WithColumns(
			dao.SysResourcesEntity.Title.Value(req.Title),
			dao.SysResourcesEntity.URL.Value(req.URL),
			dao.SysResourcesEntity.Method.Value(req.Method),
			dao.SysResourcesEntity.Icon.Value(req.Icon),
			dao.SysResourcesEntity.ResourcesType.Value(req.ResourcesType),
			dao.SysResourcesEntity.IsCache.Value(req.IsCache),
			dao.SysResourcesEntity.IsHidden.Value(req.IsHidden),
			dao.SysResourcesEntity.IsLink.Value(req.IsLink),
			dao.SysResourcesEntity.ParentID.Value(req.ParentID),
			dao.SysResourcesEntity.Sort.Value(req.Sort),
			dao.SysResourcesEntity.Description.Value(req.Description),
			dao.SysResourcesEntity.Status.Value(req.Status),
		).Build(),
	); err != nil {
		r.Fail(ctx, err, "操作失败")
		return
	}
	if err = r.syncRelatedRoleResourcesCasbin(ctx, idInt); err != nil {
		r.Fail(ctx, err, "同步角色权限失败")
		return
	}
	r.Success(ctx, "操作成功")
	return
}

// GetResourcesTreePageApi 分页获取资源树
// @Summary 分页获取资源树
// @Description 根据查询条件分页获取资源列表，并标记是否存在子节点
// @Tags 资源中心
// @Accept json
// @Produce json
// @Param data body dto.PageSysResourcesDTO true "分页查询参数"
// @Success 200 {array} vo.SysResourcesVO "资源分页列表"
// @Router /api/v1/admin/resources/pageTree [post]
func (r Resources) GetResourcesTreePageApi(ctx *gin.Context) {
	var req dto.PageSysResourcesDTO
	if !r.BindAndValidateJSON(ctx, &req) {
		return
	}
	resourcesEntities, count, err := r.ResourcesRepository.FindPageByWrapper(ctx, req.PageNumber, req.PageSize, func(g gormplus.IGenWrapper[dao.ISysResourcesEntityDo]) {
		g.
			WhereIf(req.Title != "", dao.SysResourcesEntity.Title.Like("%"+req.Title+"%")).
			WhereIf(req.Status > 0, dao.SysResourcesEntity.Status.Eq(req.Status)).
			WhereIf(req.ParentID >= 0, dao.SysResourcesEntity.ParentID.Eq(req.ParentID))
	})
	if err != nil || len(resourcesEntities) == 0 {
		r.BuildPageData(ctx, make([]interface{}, 0), 0)
		return
	}

	idList := k.Map(resourcesEntities, func(item *model.SysResourcesEntity, index int) int64 {
		return item.ID
	})
	// 判断是否有子节点
	resourcesByParentIdList, err := r.ResourcesRepository.FindList(ctx, gormplus.QueryOpt().Where(dao.SysResourcesEntity.ParentID.In(idList...)).Build())
	resourcesParentMap := make(map[int64][]*model.SysResourcesEntity)
	if err == nil && len(resourcesByParentIdList) > 0 {
		resourcesParentMap = k.GroupBy(resourcesByParentIdList, func(item *model.SysResourcesEntity) int64 {
			return item.ParentID
		})
	}
	var resourcesTreeVo = make([]vo.SysResourcesVO, 0)
	for _, item := range resourcesEntities {
		currentItem := r.ResourcesMapper.EntityToVo(item)
		currentItem.HasChildren = k.If(resourcesParentMap[item.ID] == nil, false, true)
		resourcesTreeVo = append(resourcesTreeVo, *currentItem)
	}
	r.BuildPageData(ctx, resourcesTreeVo, count)
	return
}

// GetResourcesCatalogApi 获取资源目录
// @Summary 获取资源目录
// @Description 根据资源类型获取当前账号可访问的目录、菜单或接口资源
// @Tags 资源中心
// @Accept json
// @Produce json
// @Param catalogType query int false "资源类型:1只查询模块,2查询模块和菜单,3查询模块、菜单、按钮"
// @Success 200 {array} vo.ResourcesVO "资源目录列表"
// @Router /api/v1/admin/resources/catalog [get]
func (r Resources) GetResourcesCatalogApi(ctx *gin.Context) {
	var req dto.CatalogDTO
	if err := ctx.ShouldBindQuery(&req); err != nil {
		r.Fail(ctx, err, "参数错误")
		return
	}
	accountId := ctx.GetInt64("accountId")
	isAdmin := ctx.GetInt64("isAdmin")
	// 1的时候只查询出模块,2的时候查询出模块和菜单3,的时候查询模块、菜单、按钮
	tx := dao.SysResourcesEntity.WithContext(ctx)
	if req.CatalogType == enum.ResourcesTypeCatalogEnum {
		tx = tx.Where(dao.SysResourcesEntity.ResourcesType.Eq(req.CatalogType))
	} else if req.CatalogType == enum.ResourcesTypeMenuEnum {
		tx = tx.Where(dao.SysResourcesEntity.ResourcesType.In([]int64{enum.ResourcesTypeCatalogEnum, enum.ResourcesTypeMenuEnum}...))
	} else if req.CatalogType == enum.ResourcesTypeApiEnum {
		tx = tx.Where(dao.SysResourcesEntity.ResourcesType.In([]int64{enum.ResourcesTypeCatalogEnum, enum.ResourcesTypeMenuEnum, enum.ResourcesTypeApiEnum}...))
	}
	// 判断如果当前不是超管的时候,只返回授权的资源
	if int64(isAdmin) != 1 {
		resourcesIdList := r.RoleResourcesRepository.GetResourcesByAccountId(ctx, accountId)
		if len(resourcesIdList) == 0 {
			r.Success(ctx, make([]vo.ResourcesVO, 0))
			return
		}
		tx = tx.Where(dao.SysResourcesEntity.ID.In(resourcesIdList...))
	}
	resourcesEntities, err := tx.Find()
	if err != nil || len(resourcesEntities) == 0 {
		r.Success(ctx, make([]vo.ResourcesVO, 0))
		return
	}
	resp := k.Map(resourcesEntities, func(item *model.SysResourcesEntity, index int) vo.ResourcesVO {
		return vo.ResourcesVO{
			ID:       item.ID,       // 主键id
			Title:    item.Title,    // 名称:按钮标题,或菜单标题
			ParentID: item.ParentID, // 上一级id
			Sort:     item.Sort,     // 排序
		}
	})
	r.Success(ctx, resp)
	return
}

// GetResourcesListApi 获取资源列表
// @Summary 获取资源列表
// @Description 获取正常状态的目录及菜单资源列表
// @Tags 资源中心
// @Accept json
// @Produce json
// @Success 200 {array} vo.ResourcesVO "资源列表"
// @Router /api/v1/admin/resources/list [get]
func (r Resources) GetResourcesListApi(ctx *gin.Context) {
	resourcesEntities, err := r.ResourcesRepository.FindListByWrapper(ctx, func(g gormplus.IGenWrapper[dao.ISysResourcesEntityDo]) {
		g.Where(dao.SysResourcesEntity.Status.Eq(enum.StatusNormalEnum))
	})
	if err != nil || len(resourcesEntities) == 0 {
		r.Success(ctx, make([]vo.ResourcesVO, 0))
		return
	}
	resourcesVOS := k.Map(resourcesEntities, func(item *model.SysResourcesEntity, index int) vo.ResourcesVO {
		return vo.ResourcesVO{
			ID:       item.ID,       // 主键id
			Title:    item.Title,    // 名称:按钮标题,或菜单标题
			ParentID: item.ParentID, // 上一级id
		}
	})
	r.Success(ctx, resourcesVOS)
	return
}

func (r Resources) GetDetailByIdApi(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt := cast.ToInt64(id)
	resourcesEntity, err := r.ResourcesRepository.FindById(ctx, idInt)
	if err != nil {
		r.Fail(ctx, err, "主键错误")
		return
	}
	r.Success(ctx, r.ResourcesMapper.EntityToVo(resourcesEntity))
	return
}

// syncRelatedRoleResourcesCasbin 刷新受资源变更影响的角色权限。
// 当接口资源的 url/method/status/type 被修改后，已授权该资源的角色需要重建 Casbin p 策略。
func (r Resources) syncRelatedRoleResourcesCasbin(ctx *gin.Context, resourcesID int64) error {
	roleResources, err := r.RoleResourcesRepository.FindList(ctx, gormplus.QueryOpt().
		Where(dao.SysRoleResourcesEntity.ResourcesID.Eq(resourcesID)).
		Build())
	if err != nil {
		return fmt.Errorf("查询资源关联角色失败: %w", err)
	}
	roleIDs := k.Filter(k.Distinct(k.Map(roleResources, func(item *model.SysRoleResourcesEntity, index int) int64 {
		return item.RoleID
	})), func(item int64, index int) bool {
		return item > 0
	})
	for _, roleID := range roleIDs {
		if err = r.syncRoleResourcesCasbin(ctx, roleID); err != nil {
			return err
		}
	}
	return nil
}

// syncRoleResourcesCasbin 按角色当前授权资源重建 Casbin p 策略。
// 只会同步正常状态的接口资源(resources_type=3)，目录和菜单不进入 Casbin。
func (r Resources) syncRoleResourcesCasbin(ctx *gin.Context, roleID int64) error {
	if r.Enforcer == nil {
		return errors.New("casbin enforcer未初始化")
	}
	sub := fmt.Sprintf("role_%d", roleID)
	if _, err := r.Enforcer.RemoveFilteredPolicy(0, sub); err != nil {
		return err
	}
	roleResources, err := r.RoleResourcesRepository.FindList(ctx, gormplus.QueryOpt().
		Where(dao.SysRoleResourcesEntity.RoleID.Eq(roleID)).
		Build())
	if err != nil {
		return fmt.Errorf("查询角色资源失败: %w", err)
	}
	resourceIDs := k.Filter(k.Distinct(k.Map(roleResources, func(item *model.SysRoleResourcesEntity, index int) int64 {
		return item.ResourcesID
	})), func(item int64, index int) bool {
		return item > 0
	})
	if len(resourceIDs) == 0 {
		return nil
	}
	resourcesEntities, err := r.ResourcesRepository.FindList(ctx, gormplus.QueryOpt().Where(
		dao.SysResourcesEntity.ID.In(resourceIDs...),
		dao.SysResourcesEntity.ResourcesType.Eq(enum.ResourcesTypeApiEnum),
		dao.SysResourcesEntity.Status.Eq(enum.StatusNormalEnum),
	).Build())
	if err != nil {
		return fmt.Errorf("查询资源详情失败: %w", err)
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

func NewResources(baseApi *base.BaseApi, enforcer *casbin.Enforcer) IResources {
	return Resources{
		BaseApi:                 baseApi,
		ResourcesRepository:     repository.NewSysResourcesRepository(),
		RoleResourcesRepository: repository.NewSysRoleResourcesRepository(),
		ResourcesMapper:         mapper.NewSysResourcesMapper(),
		Enforcer:                enforcer,
	}
}
