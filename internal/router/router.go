package router

import (
	"gin-admin-api/internal/api/account"
	"gin-admin-api/internal/api/auth"
	"gin-admin-api/internal/api/resources"
	"gin-admin-api/internal/api/role"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type AdminRouter struct {
	authHandler      auth.IAuth
	accountHandler   account.ISysAccount
	roleHandler      role.IRole
	resourcesHandler resources.IResources
}

func NewAdminRouter(authHandler auth.IAuth, accountHandler account.ISysAccount, roleHandler role.IRole, resourcesHandler resources.IResources) *AdminRouter {
	return &AdminRouter{
		authHandler:      authHandler,
		accountHandler:   accountHandler,
		roleHandler:      roleHandler,
		resourcesHandler: resourcesHandler,
	}
}

func (r *AdminRouter) Register(group *gin.RouterGroup, redis *redis.Client) {
	InitAuthRouter(group, redis, r.authHandler)           // 认证中心
	InitAccountRouter(group, redis, r.accountHandler)     // 账号中心
	InitRoleRouter(group, redis, r.roleHandler)           // 角色中心
	InitResourcesRouter(group, redis, r.resourcesHandler) // 资源中心
}
