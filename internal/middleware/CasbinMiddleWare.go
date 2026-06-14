package middleware

import (
	"fmt"
	"strings"

	"gin-admin-api/pkg/utils"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// CasbinMiddleWare 接口权限拦截中间件
// 必须在 AuthMiddleWare 之后挂载，因为需要从 context 读取 accountId 和 isAdmin
func CasbinMiddleWare(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if enforcer == nil {
			utils.Fail(ctx, "权限服务未初始化")
			ctx.Abort()
			return
		}

		// ① 超级管理员直接放行，不做任何权限检查
		if isAdmin, exists := ctx.Get("isAdmin"); exists && isAdmin == int64(1) {
			ctx.Next()
			return
		}

		// ② 获取当前账号 ID
		accountId, exists := ctx.Get("accountId")
		if !exists {
			utils.Fail(ctx, "无效的认证信息")
			ctx.Abort()
			return
		}

		// ③ 获取请求路径和方法
		//    必须用 ctx.FullPath() 而不是 ctx.Request.URL.Path
		//    FullPath 返回路由模板如 /api/v1/admin/account/:id
		//    这与 sys_resources.url 存储的格式一致
		//    keyMatch2 会自动处理 :id 占位符的匹配
		path := ctx.FullPath()
		if path == "" {
			path = ctx.Request.URL.Path
		}
		method := strings.ToUpper(ctx.Request.Method)

		// ④ 构造 casbin subject：user_{账号ID}
		sub := fmt.Sprintf("user_%v", accountId)

		// ⑤ 执行权限检查
		//    casbin 内部执行顺序：
		//    1. 查 g 规则：user_5 -> role_2
		//    2. 查 p 规则：role_2 能否访问 path + method
		allowed, err := enforcer.Enforce(sub, path, method)
		if err != nil {
			utils.Fail(ctx, "权限校验异常: "+err.Error())
			ctx.Abort()
			return
		}

		if !allowed {
			utils.Fail(ctx, "没有操作权限，请联系管理员授权")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
