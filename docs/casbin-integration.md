# Casbin 接口权限集成说明

本文档说明当前项目的 Casbin 集成方式，重点覆盖权限模型、业务表与 `casbin_rule` 的同步关系、运行时鉴权流程，以及常见排查点。

## 设计目标

项目使用 `sys_resources` 统一管理目录、菜单和接口资源：

| 资源类型 | 枚举 | 含义 | 是否写入 Casbin |
| --- | --- | --- | --- |
| `1` | `enum.ResourcesTypeCatalogEnum` | 目录 | 否 |
| `2` | `enum.ResourcesTypeMenuEnum` | 菜单 | 否 |
| `3` | `enum.ResourcesTypeApiEnum` | 接口 | 是 |

Casbin 只负责接口权限拦截。目录和菜单用于前端展示、菜单树、授权勾选，不直接写入 `casbin_rule`。

## Casbin 模型

初始化入口在 `initialize/casbin.go`：

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

- `sub`：请求主体，项目中是 `user_{accountId}`。
- `obj`：接口路径，使用 Gin 路由模板，如 `/api/v1/admin/account/:id`。
- `act`：HTTP 方法，如 `GET`、`POST`。
- `g`：用户和角色关系，如 `user_5 -> role_2`。
- `p`：角色和接口权限关系，如 `role_2 -> /api/v1/admin/account -> POST`。

初始化时显式开启：

```go
enforcer.EnableAutoSave(true)
```

业务接口通过 `Remove/Add` 增量维护策略，不调用 `SavePolicy()`。

## 数据关系

业务表关系：

```text
sys_account
  ↓
sys_account_role       -> Casbin g 策略: user_{accountId} -> role_{roleId}
  ↓
sys_role
  ↓
sys_role_resources     -> 角色拥有的目录/菜单/接口资源
  ↓
sys_resources(type=3)  -> Casbin p 策略: role_{roleId} -> url -> method
```

`sys_resources` 中只有满足以下条件的记录会同步到 Casbin：

```text
resources_type = enum.ResourcesTypeApiEnum
status = enum.StatusNormalEnum
url != ""
method != ""
```

如果角色只绑定目录或菜单资源，`casbin_rule` 不会新增 `p` 记录，这是预期行为。

## 运行时鉴权流程

路由挂载顺序：

```go
middleware.AuthMiddleWare(redis)
middleware.OperatorMiddleware()
middleware.CasbinMiddleWare(enforcer)
```

鉴权中间件位置：`internal/middleware/CasbinMiddleWare.go`。

流程：

1. `AuthMiddleWare` 解析 token，写入 `accountId`、`isAdmin`。
2. `OperatorMiddleware` 将 `accountId` 写入 `Request.Context()`，用于 `CreatedBy/UpdatedBy` 自动填充。
3. `CasbinMiddleWare` 判断是否为超级管理员。
4. 超级管理员 `isAdmin == 1` 直接放行。
5. 普通账号构造：

```go
sub := fmt.Sprintf("user_%v", accountId)
path := ctx.FullPath()
method := strings.ToUpper(ctx.Request.Method)
```

6. 调用：

```go
enforcer.Enforce(sub, path, method)
```

Casbin 内部先查 `g` 策略确认用户有哪些角色，再查 `p` 策略确认角色是否能访问该接口。

## 策略同步

### 角色资源同步 p 策略

方法：`internal/api/role/handler.go`

```go
syncRoleResourcesCasbin(ctx, roleId, resourcesIdList)
```

调用场景：

- 创建角色成功后
- 修改角色成功后
- 删除角色成功后

行为：

1. 删除该角色旧的 `p` 策略：

```go
RemoveFilteredPolicy(0, "role_{roleId}")
```

2. 从 `resourcesIdList` 中去重并过滤无效 ID。
3. 只查询 `resources_type = 3` 且 `status = 1` 的接口资源。
4. 写入新的 `p` 策略：

```go
AddPolicies([][]string{
    {"role_1", "/api/v1/admin/account", "POST"},
})
```

### 账号角色同步 g 策略

方法：`internal/api/account/handler.go`

```go
syncAccountRolesCasbin(accountId, roleIdList)
```

调用场景：

- 创建账号成功后
- 修改账号角色成功后
- 删除账号成功后

行为：

1. 删除该用户旧的角色关系：

```go
DeleteRolesForUser("user_{accountId}")
```

2. 从 `roleIdList` 中去重并过滤无效 ID。
3. 写入新的 `g` 策略：

```go
AddGroupingPolicies([][]string{
    {"user_5", "role_2"},
})
```

### 资源变更刷新角色策略

方法：`internal/api/resources/handler.go`

```go
syncRelatedRoleResourcesCasbin(ctx, resourcesID)
```

当资源被修改时，会查询绑定了该资源的角色，并按每个角色当前授权资源重建 `p` 策略。

这可以避免接口资源的 `url`、`method`、`status`、`resources_type` 修改后，`casbin_rule` 中残留旧权限。

## SavePolicy 使用约束

业务接口中不要调用：

```go
enforcer.SavePolicy()
```

原因：gorm-adapter 的 `SavePolicy()` 是全量保存，会先清空 `casbin_rule`，再把 enforcer 当前内存中的策略写回。如果内存策略不完整，会导致表中其他策略被清空。

当前项目中，`SavePolicy()` 只保留在一次性全量同步工具：

```go
initialize/casbinSync.go
```

业务接口使用以下增量方法即可：

- `RemoveFilteredPolicy`
- `AddPolicies`
- `DeleteRolesForUser`
- `AddGroupingPolicies`

这些方法在 `EnableAutoSave(true)` 下会自动持久化到 `casbin_rule`。

## 一次性全量同步

如果是后期接入 Casbin，已有业务表数据需要反向写入 `casbin_rule`，使用：

```go
initialize.SyncCasbinFromDB(db, enforcer)
```

该方法会从业务表重建全部策略：

- 从 `sys_role_resources + sys_resources(type=3)` 生成 `p` 策略
- 从 `sys_account_role` 生成 `g` 策略

注意：这是全量重建工具，会清空并重建 `casbin_rule`，不要在普通接口请求中调用。

## 数据录入规范

接口资源必须按 Gin 路由模板录入：

```text
url:    /api/v1/admin/account/:id
method: DELETE
```

不要录入真实请求路径：

```text
/api/v1/admin/account/123
```

原因：运行时使用 `ctx.FullPath()`，返回的是路由模板。Casbin 模型使用 `keyMatch2`，可以匹配 `:id` 参数。

## 排查清单

### 角色修改后 casbin_rule 只有 delete 没有 insert

检查角色绑定的资源是否为接口资源：

```sql
SELECT id, title, url, method, resources_type, status
FROM sys_resources
WHERE id IN (...);
```

只有 `resources_type = 3` 且 `status = 1` 且 `url/method` 不为空才会写入 `p` 策略。

### casbin_rule 中其他数据被清空

检查业务接口是否调用了：

```go
SavePolicy()
```

业务接口不应该调用它。应使用增量 `Remove/Add` 方法。

### 普通账号访问接口被拒绝

检查两类策略是否都存在：

```sql
-- 用户是否拥有角色
SELECT * FROM casbin_rule
WHERE ptype = 'g' AND v0 = 'user_账号ID';

-- 角色是否拥有接口
SELECT * FROM casbin_rule
WHERE ptype = 'p' AND v0 = 'role_角色ID';
```

还要确认 `sys_resources.url` 和 Gin 路由模板一致，`method` 与请求方法一致。

### 超级管理员仍被拦截

确认 token 中包含 `isAdmin=1`，并且 `AuthMiddleWare` 写入了：

```go
ctx.Set("isAdmin", claims.IsAdmin)
```

超级管理员在 `CasbinMiddleWare` 中直接放行。
