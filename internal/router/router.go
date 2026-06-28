package router

import (
	"gin-admin-api/internal/api/account"
	"gin-admin-api/internal/api/auth"
	"gin-admin-api/internal/api/dept"
	"gin-admin-api/internal/api/menu"
	"gin-admin-api/internal/api/post"
	"gin-admin-api/internal/api/resources"
	"gin-admin-api/internal/api/role"
	"github.com/casbin/casbin/v2"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type AdminRouter struct {
	authHandler      auth.IAuth
	accountHandler   account.ISysAccount
	roleHandler      role.IRole
	resourcesHandler resources.IResources
	menuHandler      menu.IMenu
	deptHandler      dept.IDept
	postHandler      post.IPost
}

func NewAdminRouter(authHandler auth.IAuth, accountHandler account.ISysAccount, roleHandler role.IRole, resourcesHandler resources.IResources, menuHandler menu.IMenu, deptHandler dept.IDept, postHandler post.IPost) *AdminRouter {
	return &AdminRouter{
		authHandler:      authHandler,
		accountHandler:   accountHandler,
		roleHandler:      roleHandler,
		resourcesHandler: resourcesHandler,
		menuHandler:      menuHandler,
		deptHandler:      deptHandler,
		postHandler:      postHandler,
	}
}

func (r *AdminRouter) Register(group *gin.RouterGroup, redis *redis.Client, enforcer *casbin.Enforcer, db *gorm.DB) {
	InitAuthRouter(group, redis, r.authHandler)                         // 认证中心
	InitAccountRouter(group, redis, enforcer, db, r.accountHandler)     // 账号中心
	InitRoleRouter(group, redis, enforcer, db, r.roleHandler)           // 角色中心
	InitResourcesRouter(group, redis, enforcer, db, r.resourcesHandler) // 资源中心
	InitMenuRouter(group, redis, enforcer, db, r.menuHandler)           // 菜单中心
	InitDeptRouter(group, redis, enforcer, db, r.deptHandler)           // 部门中心
	InitPostRouter(group, redis, enforcer, db, r.postHandler)           // 岗位中心
}
