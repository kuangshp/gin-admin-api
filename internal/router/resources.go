package router

import (
	"gin-admin-api/internal/api/resources"
	"gin-admin-api/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func InitResourcesRouter(Router *gin.RouterGroup, redis *redis.Client, resourcesHandler resources.IResources) {
	authRouter := Router.Group("resources", middleware.AuthMiddleWare(redis))
	authRouter.POST("", resourcesHandler.CreateResourcesApi)           // 创建资源
	authRouter.DELETE("/:id", resourcesHandler.DeleteResourcesByIdApi) // 根据id删除资源
	authRouter.PUT("/:id", resourcesHandler.ModifyResourcesByIdApi)    // 根据id修改资源
	authRouter.POST("page", resourcesHandler.GetResourcesTreePageApi)  // 分页获取资源树
	authRouter.GET("catalog", resourcesHandler.GetResourcesCatalogApi) // 获取资源目录
	authRouter.GET("list", resourcesHandler.GetResourcesListApi)       // 获取资源列表
}
