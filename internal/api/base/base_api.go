package base

import (
	"context"
	"gin-admin-api/internal/config"
	"gin-admin-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BaseApi struct {
	Logger *zap.Logger
	Db     *gorm.DB
	Cfg    *config.ServerConfig
	Redis  *redis.Client
}

func NewBaseApi(logger *zap.Logger, db *gorm.DB, cfg *config.ServerConfig, redis *redis.Client) *BaseApi {
	return &BaseApi{
		Logger: logger,
		Db:     db,
		Cfg:    cfg,
		Redis:  redis,
	}
}

// BindAndValidateJSON 绑定 JSON 并验证，失败直接返回错误响应
func (b *BaseApi) BindAndValidateJSON(ctx *gin.Context, dto interface{}) error {
	if err := ctx.ShouldBindJSON(dto); err != nil {
		utils.Fail(ctx, utils.ShowErrorMessage(err))
		return err
	}
	return nil
}

// Fail 输出日志并返回失败响应
func (b *BaseApi) Fail(ctx *gin.Context, message string, err error) {
	if err != nil {
		b.Logger.Error(message + err.Error())
	}
	utils.Fail(ctx, message)
}

// Success 返回成功响应
func (b *BaseApi) Success(ctx *gin.Context, data interface{}) {
	utils.Success(ctx, data)
}

// Ctx 从 gin.Context 中取出携带中间件数据的 context
func (b *BaseApi) Ctx(c *gin.Context) context.Context {
	return c.Request.Context()
}

// BuildPageData 构造分页响应
func (b *BaseApi) BuildPageData(ctx *gin.Context, data interface{}, total int64) {
	utils.BuildPageData(ctx, data, total)
}

// DbPage 获取分页参数
func (b *BaseApi) DbPage(pageNumber, pageSize int) (offset int, limit int) {
	return utils.DbPage(pageNumber, pageSize)
}
