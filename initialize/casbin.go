package initialize

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// NewCasbin 初始化 Casbin Enforcer
// 使用内嵌模型字符串，无需外部配置文件
// gorm-adapter 会自动创建 casbin_rule 表
func NewCasbin(db *gorm.DB) (*casbin.Enforcer, error) {
	// 第一步：创建 GORM 适配器（自动建表 casbin_rule）
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("casbin adapter 创建失败: %w", err)
	}

	// 第二步：定义 RBAC 权限模型
	// keyMatch2 支持路径参数匹配：/account/:id 可以匹配 /account/123
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")
`)
	if err != nil {
		return nil, fmt.Errorf("casbin model 解析失败: %w", err)
	}

	// 第三步：创建 Enforcer 并从数据库加载策略
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("casbin enforcer 创建失败: %w", err)
	}

	if err = enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("casbin 策略加载失败: %w", err)
	}
	enforcer.EnableAutoSave(true)

	fmt.Println("✔ Casbin 初始化成功")
	return enforcer, nil
}
