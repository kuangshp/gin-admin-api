package middleware

import (
	"gin-admin-api/internal/datascope"
	"gin-admin-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/kuangshp/gorm-plus"
	"gorm.io/gorm"
)

func DataPermissionMiddleware(db *gorm.DB) gin.HandlerFunc {
	service := datascope.NewService(db)
	return func(ctx *gin.Context) {
		accountIDValue, exists := ctx.Get("accountId")
		if !exists {
			ctx.Next()
			return
		}
		accountID, ok := accountIDValue.(int64)
		if !ok || accountID <= 0 {
			utils.Fail(ctx, "当前登录用户数据错误")
			ctx.Abort()
			return
		}

		scopeCtx, err := service.BuildContext(ctx.Request.Context(), accountID)
		if err != nil {
			utils.Fail(ctx, "获取数据权限失败")
			ctx.Abort()
			return
		}

		reqCtx := ctx.Request.Context()
		if scopeCtx.IsSuperAdmin() {
			reqCtx = gormplus.SkipDataPermission(reqCtx)
		} else {
			reqCtx = gormplus.WithDataPermission(reqCtx, datascope.Inject(scopeCtx))
		}
		ctx.Request = ctx.Request.WithContext(reqCtx)
		ctx.Set("dataScopeContext", scopeCtx)
		ctx.Next()
	}
}
