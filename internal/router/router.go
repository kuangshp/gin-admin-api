package router

import (
	"gin-admin-api/internal/api/account"
	"gin-admin-api/internal/api/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type AdminRouter struct {
	authHandler    auth.IAuth
	accountHandler account.ISysAccount
}

func NewAdminRouter(authHandler auth.IAuth, accountHandler account.ISysAccount) *AdminRouter {
	return &AdminRouter{
		authHandler:    authHandler,
		accountHandler: accountHandler,
	}
}

func (r *AdminRouter) Register(group *gin.RouterGroup, redis *redis.Client) {
	InitAuthRouter(group, redis, r.authHandler)       // 认证中心
	InitAccountRouter(group, redis, r.accountHandler) // 账号中心
}
