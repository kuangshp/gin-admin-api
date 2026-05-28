package initialize

import (
	"gin-admin-api/internal/config"
	"gin-admin-api/internal/middleware"
	"gin-admin-api/internal/router"
	"github.com/go-redis/redis/v8"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// NewRouter Wire Provider：注册所有路由，返回 *gin.Engine
func NewRouter(
	cfg *config.ServerConfig,
	logger *zap.Logger,
	adminRouter *router.AdminRouter,
	redis *redis.Client,
) *gin.Engine {
	// 配置启动模式
	gin.SetMode(cfg.Mode)
	r := gin.New()

	// 全局中间件
	r.Use(
		middleware.CorsMiddleWare(),         // 跨域的
		middleware.LoggerMiddleWare(logger), // 日志
		middleware.RecoverMiddleWare(),      // 异常的
	)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET(SwaggerJSONPath, func(c *gin.Context) {
		file, err := SwaggerJSONFile()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.File(file)
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL(SwaggerJSONPath)))

	// 配置全局路径
	ApiGroup := r.Group("/api/v1/admin")
	// 注册路由
	adminRouter.Register(ApiGroup, redis)

	return r
}
