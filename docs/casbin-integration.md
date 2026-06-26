# Casbin 权限控制接入手册

本文档基于当前项目的数据库表、路由结构和已引入的 Casbin 依赖，说明如何在项目中一步一步接入 Casbin 实现后台接口权限控制。

当前项目使用 `sys_account`、`sys_role`、`sys_resources`、`sys_account_role`、`sys_role_resources` 维护业务权限数据，使用 Casbin 的 `casbin_rule` 表维护运行时鉴权策略。

## 1. 权限设计

本项目采用 RBAC 模型：

```text
账号 sys_account
  -> 账号角色 sys_account_role
  -> 角色 sys_role
  -> 角色资源 sys_role_resources
  -> 资源 sys_resources
```

Casbin 只负责接口权限，不负责菜单展示和数据权限。

`sys_resources.resources_type` 的含义：

| 值 | 枚举 | 含义 | 是否写入 Casbin |
| --- | --- | --- | --- |
| `1` | `enum.ResourcesTypeCatalogEnum` | 目录 | 否 |
| `2` | `enum.ResourcesTypeMenuEnum` | 菜单 | 否 |
| `3` | `enum.ResourcesTypeApiEnum` | 接口 | 是 |

最终写入 Casbin 的策略分两类：

```text
g 策略：user_{accountId} -> role_{roleId}
p 策略：role_{roleId} -> 接口路由模板 -> HTTP 方法
```

示例：

```text
g, user_5, role_2
p, role_2, /api/v1/admin/account/:id, DELETE
```

## 2. 确认依赖

当前项目已经在 `go.mod` 中引入：

```go
github.com/casbin/casbin/v2
github.com/casbin/gorm-adapter/v3
```

如果是从未接入 Casbin 的分支，需要先安装：

```bash
go get github.com/casbin/casbin/v2
go get github.com/casbin/gorm-adapter/v3
go mod tidy
```

`gorm-adapter` 会使用当前项目的 GORM 数据库连接，并自动创建 Casbin 的 `casbin_rule` 表。

## 3. 确认业务表

当前项目已有权限相关业务表：

| 表名 | 作用 |
| --- | --- |
| `sys_account` | 后台账号表，`is_admin = 1` 表示超级管理员 |
| `sys_role` | 角色表 |
| `sys_account_role` | 账号和角色中间表 |
| `sys_resources` | 目录、菜单、接口资源表 |
| `sys_role_resources` | 角色和资源中间表 |

接口资源必须录入到 `sys_resources`：

```text
resources_type = 3
status = 1
url = Gin 路由模板
method = HTTP 方法
```

示例：

```text
title: 删除账号
url: /api/v1/admin/account/:id
method: DELETE
resources_type: 3
status: 1
```

注意：`url` 必须和 Gin 路由模板一致，不能保存真实请求路径。例如应保存 `/api/v1/admin/account/:id`，不要保存 `/api/v1/admin/account/123`。

## 4. 初始化 Casbin

初始化入口是 `initialize/casbin.go`。

核心步骤：

1. 使用当前 GORM `db` 创建 `gorm-adapter`。
2. 使用内嵌字符串定义 RBAC 模型。
3. 创建 `casbin.Enforcer`。
4. 从 `casbin_rule` 加载策略。
5. 开启 `AutoSave`，让增删策略时自动落库。

当前模型：

```ini
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
```

含义：

| 字段 | 当前项目含义 |
| --- | --- |
| `sub` | 请求主体，格式为 `user_{accountId}` |
| `obj` | 接口路径，使用 Gin 路由模板 |
| `act` | HTTP 方法，如 `GET`、`POST`、`PUT`、`DELETE` |
| `g` | 用户和角色关系 |
| `p` | 角色和接口权限关系 |

`keyMatch2` 支持 `:id` 这种路径参数匹配，因此 `/api/v1/admin/account/:id` 可以匹配实际请求 `/api/v1/admin/account/123`。

当前初始化代码关键点：

```go
adapter, err := gormadapter.NewAdapterByDB(db)
enforcer, err := casbin.NewEnforcer(m, adapter)
err = enforcer.LoadPolicy()
enforcer.EnableAutoSave(true)
```

## 5. 注入 Enforcer

当前项目使用 Wire 组装依赖，入口在 `wire.go`：

```go
initialize.NewCasbin,
account.NewSysAccount,
role.NewRole,
resources.NewResources,
router.NewAdminRouter,
initialize.NewRouter,
```

`NewCasbin(db)` 依赖数据库连接，返回 `*casbin.Enforcer`。

需要使用 Casbin 的模块通过构造函数接收：

```go
func NewRole(baseApi *base.BaseApi, enforcer *casbin.Enforcer) IRole
func NewSysAccount(baseApi *base.BaseApi, enforcer *casbin.Enforcer) ISysAccount
func NewResources(baseApi *base.BaseApi, enforcer *casbin.Enforcer) IResources
```

如果调整了 `wire.go`，需要重新生成：

```bash
make wire
```

## 6. 添加鉴权中间件

中间件位置：`internal/middleware/CasbinMiddleWare.go`。

挂载顺序必须在认证中间件之后，因为 Casbin 鉴权需要从 Gin context 读取 `accountId` 和 `isAdmin`：

```go
middleware.AuthMiddleWare(redis),
middleware.OperatorMiddleware(),
middleware.CasbinMiddleWare(enforcer),
```

当前已经在账号、角色、资源路由组中挂载：

```go
Router.Group(
    "account",
    middleware.AuthMiddleWare(redis),
    middleware.OperatorMiddleware(),
    middleware.CasbinMiddleWare(enforcer),
)
```

运行时鉴权流程：

1. `AuthMiddleWare` 解析 token，写入 `accountId`、`isAdmin`。
2. `OperatorMiddleware` 写入操作人上下文。
3. `CasbinMiddleWare` 判断是否超级管理员。
4. `isAdmin == 1` 直接放行。
5. 普通账号构造 `sub = user_{accountId}`。
6. 使用 `ctx.FullPath()` 获取路由模板。
7. 使用 `ctx.Request.Method` 获取 HTTP 方法。
8. 调用 `enforcer.Enforce(sub, path, method)`。

核心代码：

```go
path := ctx.FullPath()
method := strings.ToUpper(ctx.Request.Method)
sub := fmt.Sprintf("user_%v", accountId)

allowed, err := enforcer.Enforce(sub, path, method)
```

如果新增后台接口并需要权限控制，应把 `CasbinMiddleWare(enforcer)` 挂到对应路由组。

## 7. 同步角色接口权限

角色和资源关系保存在 `sys_role_resources`，但 Casbin 运行时需要的是 `p` 策略。

当前项目在 `internal/api/role/handler.go` 中通过 `syncRoleResourcesCasbin` 同步：

```go
syncRoleResourcesCasbin(ctx, roleId, resourcesIdList)
```

调用场景：

| 场景 | 处理 |
| --- | --- |
| 创建角色 | 创建 `sys_role` 和 `sys_role_resources` 后，同步 `p` 策略 |
| 修改角色 | 重建 `sys_role_resources` 后，同步 `p` 策略 |
| 删除角色 | 删除业务数据后，清空该角色 `p` 策略 |

同步逻辑：

1. 构造角色主体：

```go
sub := fmt.Sprintf("role_%d", roleId)
```

2. 清理该角色旧接口权限：

```go
enforcer.RemoveFilteredPolicy(0, sub)
```

3. 根据 `resourcesIdList` 查询接口资源，只保留：

```text
resources_type = 3
status = 1
url != ""
method != ""
```

4. 写入新的 `p` 策略：

```go
enforcer.AddPolicies([][]string{
    {"role_2", "/api/v1/admin/account/:id", "DELETE"},
})
```

由于已经执行 `enforcer.EnableAutoSave(true)`，`AddPolicies` 和 `RemoveFilteredPolicy` 会自动同步到 `casbin_rule`。

## 8. 同步账号角色关系

账号和角色关系保存在 `sys_account_role`，但 Casbin 运行时需要的是 `g` 策略。

当前项目在 `internal/api/account/handler.go` 中通过 `syncAccountRolesCasbin` 同步：

```go
syncAccountRolesCasbin(accountId, roleIdList)
```

调用场景：

| 场景 | 处理 |
| --- | --- |
| 创建账号 | 创建 `sys_account` 和 `sys_account_role` 后，同步 `g` 策略 |
| 修改账号 | 重建账号角色关系后，同步 `g` 策略 |
| 删除账号 | 删除业务数据后，清空该用户 `g` 策略 |

同步逻辑：

1. 构造用户主体：

```go
sub := fmt.Sprintf("user_%d", accountId)
```

2. 清理该用户旧角色：

```go
enforcer.DeleteRolesForUser(sub)
```

3. 写入新的 `g` 策略：

```go
enforcer.AddGroupingPolicies([][]string{
    {"user_5", "role_2"},
    {"user_5", "role_3"},
})
```

## 9. 资源变更后刷新策略

接口资源的 `url`、`method`、`status`、`resources_type` 发生变化时，已经授权该资源的角色策略也必须刷新，否则 `casbin_rule` 中可能残留旧接口权限。

当前项目在 `internal/api/resources/handler.go` 中处理：

```go
syncRelatedRoleResourcesCasbin(ctx, resourcesID)
```

处理流程：

1. 根据 `resourcesID` 查询 `sys_role_resources`，找到绑定该资源的角色。
2. 对每个角色调用 `syncRoleResourcesCasbin(ctx, roleID)`。
3. 重新查询该角色当前拥有的全部正常接口资源。
4. 清空并重建该角色的 `p` 策略。

资源删除时也需要兜底清理相关角色策略。

## 10. 历史数据一次性同步

如果是在已有业务数据的项目中后期接入 Casbin，需要把 `sys_account_role` 和 `sys_role_resources` 的历史数据反向同步到 `casbin_rule`。

当前项目提供了同步工具：

```go
initialize.SyncCasbinFromDB(db, enforcer)
```

命令：

```bash
make sync-casbin
```

生产配置：

```bash
make sync-casbin-prod
```

也可以直接执行：

```bash
go run ./cmd/sync_casbin.go application.dev.yml
```

同步过程：

1. 清空 `casbin_rule`。
2. 从 `sys_role_resources INNER JOIN sys_resources` 生成 `p` 策略。
3. 从 `sys_account_role` 生成 `g` 策略。
4. 调用 `SavePolicy()` 兜底持久化。

注意：这是一次性全量重建工具，会删除并重建全部 Casbin 策略，不要在普通业务接口中调用。

## 11. SavePolicy 使用约束

普通业务接口中不要调用：

```go
enforcer.SavePolicy()
```

原因：`gorm-adapter` 的 `SavePolicy()` 是全量保存，可能先清空 `casbin_rule`，再把当前 Enforcer 内存中的策略写回。如果内存策略不完整，会导致其他策略丢失。

业务接口只使用增量方法：

```go
RemoveFilteredPolicy
AddPolicies
DeleteRolesForUser
AddGroupingPolicies
```

当前项目只有一次性同步工具 `initialize/casbinSync.go` 会调用 `SavePolicy()`。

## 12. 新增接口的接入步骤

新增一个需要权限控制的后台接口时，按下面步骤处理：

1. 在 Gin 路由中注册接口，并确保路由组挂载 `CasbinMiddleWare(enforcer)`。
2. 在 `sys_resources` 中新增一条接口资源：

```text
resources_type = 3
url = ctx.FullPath() 对应的路由模板
method = HTTP 方法大写
status = 1
```

3. 在角色管理中把该接口资源授权给角色，写入 `sys_role_resources`。
4. 角色创建或修改成功后，调用 `syncRoleResourcesCasbin` 重建该角色 `p` 策略。
5. 账号绑定角色后，调用 `syncAccountRolesCasbin` 写入该账号 `g` 策略。
6. 使用普通账号请求接口，验证有权限时放行、无权限时拒绝。

## 13. 验证 SQL

查看某用户是否拥有角色：

```sql
SELECT *
FROM casbin_rule
WHERE ptype = 'g'
  AND v0 = 'user_账号ID';
```

查看某角色是否拥有接口权限：

```sql
SELECT *
FROM casbin_rule
WHERE ptype = 'p'
  AND v0 = 'role_角色ID';
```

查看某接口资源是否可同步到 Casbin：

```sql
SELECT id, title, url, method, resources_type, status
FROM sys_resources
WHERE id IN (...);
```

应满足：

```text
resources_type = 3
status = 1
url 不为空
method 不为空
```

## 14. 常见问题

### 普通账号访问接口被拒绝

检查三点：

1. `casbin_rule` 是否存在 `g, user_{accountId}, role_{roleId}`。
2. `casbin_rule` 是否存在 `p, role_{roleId}, 当前接口路由模板, 当前 HTTP 方法`。
3. `sys_resources.url` 是否和 `ctx.FullPath()` 一致。

### 角色修改后 casbin_rule 只有删除没有新增

通常是角色绑定的资源不是接口资源，或接口资源状态不正常。

只会写入满足下面条件的资源：

```text
resources_type = 3
status = 1
url != ""
method != ""
```

### 路径带 ID 的接口匹配失败

确认 `sys_resources.url` 保存的是 Gin 路由模板：

```text
/api/v1/admin/account/:id
```

不要保存真实请求路径：

```text
/api/v1/admin/account/123
```

运行时中间件使用 `ctx.FullPath()`，Casbin 模型使用 `keyMatch2` 匹配路径参数。

### 超级管理员仍然被拦截

确认登录 token 中包含 `isAdmin = 1`，并且 `AuthMiddleWare` 已写入：

```go
ctx.Set("isAdmin", claims.IsAdmin)
```

`CasbinMiddleWare` 中会对 `isAdmin == int64(1)` 直接放行。

### casbin_rule 中其他策略被清空

检查业务接口是否调用了：

```go
SavePolicy()
```

业务接口应使用增量 API，不应调用 `SavePolicy()`。

## 15. 当前项目落点

| 能力 | 文件 |
| --- | --- |
| Casbin 初始化 | `initialize/casbin.go` |
| 历史数据同步 | `initialize/casbinSync.go`、`cmd/sync_casbin.go` |
| 鉴权中间件 | `internal/middleware/CasbinMiddleWare.go` |
| 路由挂载 | `internal/router/account.go`、`internal/router/role.go`、`internal/router/resources.go` |
| 角色资源同步 p 策略 | `internal/api/role/handler.go` |
| 账号角色同步 g 策略 | `internal/api/account/handler.go` |
| 资源变更刷新策略 | `internal/api/resources/handler.go` |
| Wire 依赖注入 | `wire.go`、`wire_gen.go` |
| 同步命令 | `Makefile` 的 `sync-casbin`、`sync-casbin-prod` |
