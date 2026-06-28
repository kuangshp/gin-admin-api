package router

import (
	"gin-admin-api/internal/api/post"
	"gin-admin-api/internal/middleware"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func InitPostRouter(Router *gin.RouterGroup, redis *redis.Client, enforcer *casbin.Enforcer, db *gorm.DB, postHandler post.IPost) {
	authRouter := Router.Group(
		"post",
		middleware.AuthMiddleWare(redis),
		middleware.OperatorMiddleware(),
		middleware.DataPermissionMiddleware(db),
		middleware.CasbinMiddleWare(enforcer),
	)
	authRouter.POST("", postHandler.CreatePostApi)
	authRouter.DELETE("/:id", postHandler.DeletePostByIdApi)
	authRouter.PUT("/:id", postHandler.ModifyPostByIdApi)
	authRouter.POST("page", postHandler.GetPostPageApi)
	authRouter.GET("list", postHandler.GetPostListApi)
	authRouter.GET("detail/:id", postHandler.GetPostDetailByIdApi)
}
