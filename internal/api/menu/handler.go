package menu

import (
	"gin-admin-api/internal/api/base"
	"gin-admin-api/internal/api/menu/vo"
	"gin-admin-api/internal/dal/dao"
	"gin-admin-api/internal/dal/model"
	"gin-admin-api/internal/dal/repository"
	"gin-admin-api/pkg/enum"
	"github.com/gin-gonic/gin"
	"github.com/kuangshp/go-utils/k"
	"github.com/kuangshp/gorm-plus"
)

type IMenu interface {
	GetMenusApi(ctx *gin.Context)
}

type Menu struct {
	*base.BaseApi
	ResourcesRepository     repository.SysResourcesRepository
	RoleResourcesRepository repository.SysRoleResourcesRepository
}

// GetMenusApi 获取菜单
// @Summary 获取菜单
// @Description 获取当前登录账号可访问的菜单资源列表，超级管理员返回全部正常资源
// @Tags 菜单中心
// @Accept json
// @Produce json
// @Success 200 {array} vo.MenusVO "菜单列表"
// @Router /api/v1/admin/menu [get]
func (m Menu) GetMenusApi(ctx *gin.Context) {
	accountId := ctx.GetInt64("accountId")
	isAdmin := ctx.GetInt64("isAdmin")
	// 普通用户只能查看授权的资源
	var (
		resourcesEntities = make([]*model.SysResourcesEntity, 0)
		err               error
	)
	if isAdmin != enum.AdminAccount {
		// 根据角色查询授权的资源id
		resourcesIdList := m.RoleResourcesRepository.GetResourcesByAccountId(ctx, accountId)
		if len(resourcesIdList) == 0 {
			m.Success(ctx, make([]vo.MenusVO, 0))
			return
		}

		// 根据资源id查询资源
		resourcesEntities, err = m.ResourcesRepository.FindByIdList(ctx, resourcesIdList)
		if err != nil {
			m.Success(ctx, make([]vo.MenusVO, 0))
			return
		}
	} else {
		// 超管直接查询
		resourcesEntities, err = m.ResourcesRepository.FindList(ctx, gormplus.QueryOpt().Where(
			dao.SysResourcesEntity.Status.Eq(enum.StatusNormalEnum),
		).Order(dao.SysResourcesEntity.Sort).Build())
		if err != nil {
			m.Success(ctx, make([]vo.MenusVO, 0))
			return
		}
	}
	if len(resourcesEntities) == 0 {
		m.Success(ctx, make([]vo.MenusVO, 0))
		return
	}
	menusVO := k.Map(resourcesEntities, func(item *model.SysResourcesEntity, index int) vo.MenusVO {
		return vo.MenusVO{
			ID:            item.ID,            // 主键id
			Title:         item.Title,         // 名称:按钮标题,或菜单标题
			URL:           item.URL,           // 按钮请求url,或菜单路由
			Icon:          item.Icon,          // 菜单小图标
			ResourcesType: item.ResourcesType, // 类型:0表示目录,1表示菜单,2表示接口
			IsCache:       item.IsCache,       // 是否缓存:0表示缓存:1不缓存
			IsHidden:      item.IsHidden,      // 是否隐藏:0表示不隐藏,1表示隐藏
			IsLink:        item.IsLink,        // 是否为外部链接:0表示不是,1表示是
			ParentID:      item.ParentID,      // 上一级id
			Sort:          item.Sort,          // 菜单,或按钮排序
		}
	})
	m.Success(ctx, menusVO)
	return
}

func NewMenu(baseApi *base.BaseApi) IMenu {
	return Menu{
		BaseApi:                 baseApi,
		ResourcesRepository:     repository.NewSysResourcesRepository(),
		RoleResourcesRepository: repository.NewSysRoleResourcesRepository(),
	}
}
