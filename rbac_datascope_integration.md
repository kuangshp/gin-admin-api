# Gin Admin API 数据权限集成实施文档

> 本文档基于当前 `gin-admin-api.sql` 表设计整理，目标是在现有 Gin + GORM Gen / gorm-plus Repository + Casbin 功能权限基础上补齐“访问接口后能看到哪些数据”的数据权限能力。

## 1. 当前项目现状

当前项目已经具备 RBAC 功能权限链路：

```text
sys_account
  └── sys_account_role
        └── sys_role
              └── sys_role_resources
                    └── sys_resources
```

相关代码位置：

| 能力 | 当前文件 | 现状 |
|---|---|---|
| 登录认证 | `internal/api/auth/handler.go` | 登录后 JWT 写入 `accountId`、`username`、`isAdmin` |
| JWT 解析 | `internal/middleware/AuthMiddleWare.go` | 将 `accountId`、`userName`、`isAdmin` 写入 `gin.Context` |
| 接口权限 | `internal/middleware/CasbinMiddleWare.go` | 通过 `user_{accountId}`、`role_{roleId}` 校验 URL + Method |
| 角色资源授权 | `internal/api/role/handler.go` | 创建/修改角色时维护 `sys_role_resources` 并同步 Casbin |
| 账号角色授权 | `internal/api/account/handler.go` | 创建/修改账号时维护 `sys_account_role` 并同步 Casbin |
| 数据权限表设计 | `gin-admin-api.sql` | 已有 `sys_account.dept_id`、`sys_role.data_scope`、`sys_dept`、`sys_role_custom_dept` |
| Repository | `internal/dal/repository/*` | 由 gorm-plus 生成，支持自定义扩展文件 |

当前项目已经完成“能不能访问接口”的功能权限，但还没有完成“访问接口后能看到哪些数据”的数据权限。

## 2. 当前 SQL 数据权限表关系

当前 `gin-admin-api.sql` 采用“账号直属部门 + 角色内置数据范围 + 角色自定义部门”的设计：

```text
sys_account.dept_id
  └── sys_dept.id

sys_account
  └── sys_account_role
        └── sys_role.data_scope
              └── sys_role_custom_dept.dept_id
                    └── sys_dept.id
```

完整关系如下：

| 表 / 字段 | 作用 |
|---|---|
| `sys_account` | 后台账号主体 |
| `sys_account.dept_id` | 账号所属部门，直接用于数据权限计算 |
| `sys_account.is_admin` | `1` 为超级管理员，跳过数据权限过滤 |
| `sys_role` | 角色表 |
| `sys_role.data_scope` | 角色数据范围：`1` 全部、`2` 本部门、`3` 本部门及下级、`4` 仅本人、`5` 自定义部门 |
| `sys_account_role` | 账号与角色关系表 |
| `sys_resources` | 菜单、接口等资源表 |
| `sys_role_resources` | 角色与资源关系表，用于功能权限 |
| `sys_dept` | 部门表，用于部门维度数据权限 |
| `sys_dept.full_id` | 部门全层级 ID，例：`1,5,12`，用于查询子部门 |
| `sys_role_custom_dept` | 角色自定义数据权限部门表 |
| `created_by` 字段 | 多张表已有，用于“仅本人数据”过滤 |

注意：当前 SQL 不再使用 `sys_data_scope` / `sys_data_scope_dept`。文档和代码实现都应以 `sys_role.data_scope`、`sys_role_custom_dept` 为准。

## 3. 数据权限目标

集成后，本项目需要支持：

| 范围值 | 名称 | 规则 |
|---|---|---|
| `1` | 全部数据 | 当前账号可查看所有数据 |
| `2` | 本部门数据 | 当前账号所属部门的数据 |
| `3` | 本部门及下级数据 | 当前账号所属部门及其子部门的数据 |
| `4` | 仅本人数据 | `created_by = 当前账号ID` |
| `5` | 自定义部门数据 | 角色绑定的指定部门数据 |

超级管理员 `sys_account.is_admin = 1` 永远跳过数据权限过滤。

## 4. 总体落地顺序

建议按下面顺序集成：

1. 确认当前数据库执行了最新 `gin-admin-api.sql`。
2. 重新生成 `sys_dept`、`sys_post`、`sys_account_post`、`sys_role_custom_dept` 相关 DAO / Model / Repository。
3. 新增数据权限枚举。
4. 在账号创建、修改、详情接口中接入 `deptId`。
5. 在角色创建、修改、详情接口中接入 `dataScope`、`deptIdList`。
6. 新增数据权限上下文和构建服务。
7. 新增数据权限查询构建器。
8. 在需要控制数据范围的列表、详情、修改、删除、导出、统计接口中应用构建器。
9. 增加测试用例和接口验证。

## 5. 数据库表说明与实现约束

### 5.1 `sys_account`

当前 SQL 已有：

```sql
`dept_id` int(11) NOT NULL COMMENT '所属部门id，用于数据权限计算',
KEY `idx_dept_id` (`dept_id`)
```

账号所属部门直接来自 `sys_account.dept_id`，不需要再通过岗位反查部门。

创建或修改账号时必须校验：

1. `dept_id` 必填。
2. 部门存在。
3. 部门未软删除。
4. 部门状态为正常。

如果账号没有合法部门，普通账号的数据权限按空范围处理，不能放开为全部数据。

### 5.2 `sys_role`

当前 SQL 已有：

```sql
`data_scope` tinyint(4) NOT NULL DEFAULT 4 COMMENT '数据范围：1=全部 2=本部门 3=本部门及下级 4=仅本人 5=自定义部门',
KEY `idx_data_scope` (`data_scope`)
```

这表示一个角色只有一个数据范围配置。角色创建和修改时，应直接写入 `sys_role.data_scope`。

推荐默认值继续使用 `4 仅本人`，避免未配置时误放开。

### 5.3 `sys_role_custom_dept`

当前 SQL 已有：

```sql
CREATE TABLE `sys_role_custom_dept` (
    `id` int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '主键id',
    `role_id` int(11) NOT NULL COMMENT '角色id',
    `dept_id` int(11) NOT NULL COMMENT '可查看部门id',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间',
    UNIQUE KEY `uk_role_dept` (`role_id`, `dept_id`),
    KEY `idx_role_id` (`role_id`),
    KEY `idx_dept_id` (`dept_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色自定义数据权限部门表';
```

该表只在 `sys_role.data_scope = 5` 时使用。

由于唯一索引是：

```sql
UNIQUE KEY `uk_role_dept` (`role_id`, `dept_id`)
```

如果使用软删除，同一个角色和部门删除后不能再次创建相同关系。推荐角色自定义部门变更时采用物理删除旧关系后批量创建新关系，和当前角色资源授权的处理方式保持一致。

### 5.4 `sys_dept`

当前 SQL 已有：

```sql
`parent_id` int(11) NOT NULL DEFAULT 0 COMMENT '上级部门id',
`full_id` varchar(255) NOT NULL DEFAULT '' COMMENT '全层级ID，例：1,5,12',
`full_name` varchar(255) NOT NULL DEFAULT '' COMMENT '全层级名称',
KEY `idx_parent_id` (`parent_id`),
KEY `idx_full_id` (`full_id`),
KEY `idx_status` (`status`)
```

部门树查询建议使用 `full_id`：

```sql
SELECT id
FROM sys_dept
WHERE deleted_at IS NULL
  AND status = 1
  AND (id = ? OR FIND_IN_SET(?, full_id));
```

创建、修改、移动部门时必须维护 `full_id` 和 `full_name`。否则“本部门及下级”会计算错误。

建议补充组合索引：

```sql
ALTER TABLE `sys_dept`
ADD KEY `idx_sys_dept_parent_status` (`parent_id`, `status`);
```

如果部门名称要求同一个上级下不重复，建议在业务层校验：

```sql
SELECT id
FROM sys_dept
WHERE parent_id = ?
  AND name = ?
  AND deleted_at IS NULL
LIMIT 1;
```

### 5.5 `sys_post` 和 `sys_account_post`

当前 SQL 仍保留岗位和账号岗位关联：

```text
sys_post
sys_account_post
```

但在当前数据权限设计中，账号部门以 `sys_account.dept_id` 为准，岗位不参与数据权限计算。

岗位可继续用于组织管理、岗位管理、账号扩展属性。若后续业务希望“一个账号多个部门”，再考虑改回通过 `sys_account_post -> sys_post.dept_id` 计算账号部门集合，但这不是当前 SQL 的主设计。

## 6. 重新生成 DAO / Model / Repository

当前项目使用 `cmd/generator.go` 生成：

```text
internal/dal/model
internal/dal/dao
internal/dal/repository
internal/dal/dto
internal/dal/vo
internal/dal/mapper
```

执行前确认本地 MySQL 中已执行最新 `gin-admin-api.sql`，并且 `cmd/generator.go` 指向正确数据库。

生成命令：

```bash
go run cmd/generator.go
```

生成后应出现或更新：

```text
internal/dal/model/sys_dept.gen.go
internal/dal/model/sys_post.gen.go
internal/dal/model/sys_account_post.gen.go
internal/dal/model/sys_role_custom_dept.gen.go

internal/dal/dao/sys_dept.gen.go
internal/dal/dao/sys_post.gen.go
internal/dal/dao/sys_account_post.gen.go
internal/dal/dao/sys_role_custom_dept.gen.go

internal/dal/repository/sysdept_gen.go
internal/dal/repository/syspost_gen.go
internal/dal/repository/sysaccountpost_gen.go
internal/dal/repository/sysrolecustomdept_gen.go
```

同时已有的 `sys_account`、`sys_role` 模型也应重新生成，以包含：

```text
sys_account.dept_id
sys_role.data_scope
```

生成器不会覆盖自定义 repository 文件时，可以继续在非 `_gen.go` 文件里增加扩展方法；生成后仍建议通过 `git diff` 检查。

## 7. 新增数据权限枚举

新增文件：

```text
pkg/enum/data_scope.enum.go
```

建议内容：

```go
package enum

const (
	DataScopeAll          int64 = 1 // 全部数据
	DataScopeDept         int64 = 2 // 本部门数据
	DataScopeDeptAndChild int64 = 3 // 本部门及下级数据
	DataScopeSelf         int64 = 4 // 仅本人数据
	DataScopeCustomDept   int64 = 5 // 自定义部门数据
)
```

## 8. 账号接口接入 `deptId`

### 8.1 DTO 修改

当前 `CreateSysAccountDTO` 增加：

```go
DeptID int64 `json:"deptId" validate:"required,gte=1"` // 所属部门id
```

当前 `ModifySysAccountDTO` 增加：

```go
DeptID int64 `json:"deptId" validate:"required,gte=1"` // 所属部门id
```

### 8.2 VO 修改

当前 `SysAccountVO` 增加：

```go
DeptID int64 `json:"deptId"` // 所属部门id
```

详情结构继续返回角色：

```go
type SysAccountDetailVO struct {
	SysAccountVO
	RoleIdList []int64 `json:"roleIdList"`
}
```

如果前端需要展示部门名称，可追加：

```go
DeptName string `json:"deptName"`
```

### 8.3 创建账号

创建账号事务内建议流程：

```text
1. 校验 deptId 对应部门存在、启用、未删除。
2. 创建 sys_account，写入 dept_id。
3. 批量创建 sys_account_role。
4. 提交事务后同步 Casbin。
```

### 8.4 修改账号

修改账号事务内建议流程：

```text
1. 校验 deptId 对应部门存在、启用、未删除。
2. 更新 sys_account 基础字段和 dept_id。
3. 物理删除旧 sys_account_role。
4. 批量创建新的 sys_account_role。
5. 提交事务后同步 Casbin。
```

账号角色关系表当前唯一索引为：

```sql
UNIQUE KEY `uk_account_role` (`account_id`, `role_id`)
```

推荐继续使用物理删除关系行后重建，避免软删除唯一索引冲突。

## 9. 角色接口接入数据权限配置

### 9.1 DTO 修改

当前 `CreateSysRoleDTO` 建议改为：

```go
type CreateSysRoleDTO struct {
	Name            string  `json:"name" validate:"required"`
	Description     string  `json:"description"`
	Status          int64   `json:"status" validate:"oneof=1 2"`
	Sort            int64   `json:"sort"`
	ResourcesIdList []int64 `json:"resourcesIdList" validate:"required,min=1"`
	DataScope       int64   `json:"dataScope" validate:"required,oneof=1 2 3 4 5"`
	DeptIdList      []int64 `json:"deptIdList"`
}
```

如果后续增加修改角色 DTO，也应包含相同字段。

校验规则：

| 条件 | 规则 |
|---|---|
| `dataScope = 5` | `deptIdList` 必须非空，且所有部门存在、启用、未删除 |
| `dataScope != 5` | 忽略 `deptIdList`，保存前删除该角色旧的自定义部门关系 |
| `dataScope` 为空 | 默认不建议静默补值，接口层应明确要求前端传入 |

### 9.2 VO 修改

当前 `SysRoleVO` 增加：

```go
DataScope int64 `json:"dataScope"` // 数据范围
```

当前 `SysRoleDetailVO` 增加：

```go
type SysRoleDetailVO struct {
	SysRoleVO
	ResourcesIdList []int64 `json:"resourcesIdList"`
	DeptIdList      []int64 `json:"deptIdList"`
}
```

### 9.3 创建角色

在 `internal/api/role/handler.go` 的创建角色事务内增加数据权限处理：

```text
1. 创建 sys_role，写入 data_scope。
2. 创建 sys_role_resources。
3. 如果 dataScope=5，批量创建 sys_role_custom_dept。
4. 提交事务后同步 Casbin。
```

伪代码：

```go
roleEntity := &model.SysRoleEntity{
	Name:        req.Name,
	Description: req.Description,
	Status:      req.Status,
	Sort:        req.Sort,
	DataScope:   req.DataScope,
	CreatedBy:   currentAccountID,
}
```

自定义部门关系：

```go
if req.DataScope == enum.DataScopeCustomDept {
	rows := buildRoleCustomDeptEntityList(roleEntity.ID, req.DeptIdList)
	err = s.SysRoleCustomDeptRepository.CreateBatchTx(ctx, tx, rows)
}
```

### 9.4 修改角色

修改角色事务内建议流程：

```text
1. 更新 sys_role 基础字段和 data_scope。
2. 物理删除旧 sys_role_resources。
3. 批量创建新的 sys_role_resources。
4. 物理删除旧 sys_role_custom_dept。
5. dataScope=5 时批量创建新的 sys_role_custom_dept。
6. 提交事务后同步 Casbin。
```

### 9.5 删除角色

删除角色事务内建议补充：

```text
1. 删除 sys_role_custom_dept by role_id。
2. 删除 sys_role_resources by role_id。
3. 删除 sys_account_role by role_id，或拒绝删除已绑定账号的角色。
4. 删除 sys_role。
5. 提交事务后同步 Casbin。
```

是否允许删除已绑定账号的角色，需要产品规则明确。更稳妥的默认策略是拒绝删除，避免账号权限被静默改变。

### 9.6 角色详情

角色详情返回：

```text
1. 查询 sys_role，返回 dataScope。
2. 查询 sys_role_resources，返回 resourcesIdList。
3. dataScope=5 时查询 sys_role_custom_dept，返回 deptIdList。
4. dataScope!=5 时 deptIdList 返回空数组。
```

## 10. 定义运行时数据权限上下文

新增包：

```text
internal/datascope
```

新增文件：

```text
internal/datascope/context.go
```

建议结构：

```go
package datascope

type Context struct {
	AccountID int64
	IsAdmin   int64
	DeptID    int64
	RoleIDs   []int64
	DataScope int64
	DeptIDs   []int64
}

func (c Context) IsSuperAdmin() bool {
	return c.IsAdmin == 1
}
```

字段说明：

| 字段 | 来源 | 用途 |
|---|---|---|
| `AccountID` | JWT / `sys_account.id` | 仅本人数据使用 |
| `IsAdmin` | JWT / `sys_account.is_admin` | 超管跳过数据权限 |
| `DeptID` | `sys_account.dept_id` | 本部门、本部门及下级使用 |
| `RoleIDs` | `sys_account_role` | 加载角色数据权限 |
| `DataScope` | `sys_role.data_scope` | 最终生效的数据范围 |
| `DeptIDs` | `sys_account.dept_id`、`sys_dept` 子部门、`sys_role_custom_dept` | 最终可访问部门范围 |

## 11. 多角色数据权限合并规则

当前项目账号可以绑定多个角色：

```text
sys_account -> sys_account_role -> sys_role
```

因此一个账号可能同时拥有多个 `data_scope`。必须定义合并规则，否则查询结果不可预测。

推荐规则：取权限最大的范围，同时对部门类权限做并集。

优先级：

```text
全部数据 1
  > 本部门及下级 3
  > 自定义部门 5
  > 本部门 2
  > 仅本人 4
```

合并规则：

1. 任一角色为 `1 全部数据`，最终为全部数据。
2. 任一角色为 `3 本部门及下级`，最终至少包含当前账号部门及其子部门。
3. 多个角色都是 `5 自定义部门`，部门 ID 取并集。
4. `5 自定义部门` 与 `2 本部门` 同时存在，部门 ID 取并集。
5. `3 本部门及下级` 与 `5 自定义部门` 同时存在，部门 ID 取“本部门及下级 + 自定义部门”并集。
6. 只有 `4 仅本人` 时，按本人过滤。
7. 没有任何有效角色时，默认按 `4 仅本人` 或直接返回空范围；推荐返回空范围更严格。

如果最终合并结果是部门集合权限，可将 `DataScope` 规范化为 `5 自定义部门`，并把并集后的部门写入 `DeptIDs`。这样查询构建器只需要处理“全部、仅本人、部门集合”三类核心条件。

## 12. 新增 Repository 自定义方法

建议在自定义 repository 文件中增加方法，不直接改生成文件。

### 12.1 角色数据权限查询

建议新增：

```text
internal/dal/repository/sysrole.go
```

扩展方法：

```go
type RoleDataScope struct {
	RoleID    int64
	DataScope int64
	DeptIDs   []int64
}

func (r *SysRoleRepository) FindRoleDataScopes(ctx context.Context, roleIDs []int64) ([]RoleDataScope, error)
```

查询逻辑：

1. 从 `sys_role` 查询 `id`、`data_scope`。
2. 对 `data_scope=5` 的角色查询 `sys_role_custom_dept`。
3. 返回每个角色对应的数据权限。

### 12.2 自定义部门查询

建议新增：

```text
internal/dal/repository/sysrolecustomdept.go
```

扩展方法：

```go
func (r *SysRoleCustomDeptRepository) FindDeptIDsByRoleIDs(ctx context.Context, roleIDs []int64) (map[int64][]int64, error)
```

SQL 逻辑：

```sql
SELECT role_id, dept_id
FROM sys_role_custom_dept
WHERE role_id IN ?
  AND deleted_at IS NULL;
```

### 12.3 子部门查询

建议新增：

```text
internal/dal/repository/sysdept.go
```

扩展方法：

```go
func (r *SysDeptRepository) FindChildDeptIDs(ctx context.Context, deptID int64) ([]int64, error)
```

推荐 SQL：

```sql
SELECT id
FROM sys_dept
WHERE deleted_at IS NULL
  AND status = 1
  AND (id = ? OR FIND_IN_SET(?, full_id));
```

### 12.4 账号上下文查询

数据权限构建时仅依赖 JWT 不够，因为 JWT 中没有 `dept_id`，并且账号部门可能变更。建议通过账号 ID 查询最新账号：

```go
type AccountDataScopeBase struct {
	AccountID int64
	IsAdmin   int64
	DeptID    int64
}

func (r *SysAccountRepository) FindDataScopeBaseByID(ctx context.Context, accountID int64) (*AccountDataScopeBase, error)
```

查询需过滤：

```sql
deleted_at IS NULL
status = 1
```

## 13. 新增数据权限服务

新增文件：

```text
internal/datascope/service.go
```

建议提供：

```go
type Service interface {
	BuildContext(ctx *gin.Context) (*Context, error)
}
```

`BuildContext` 执行流程：

```text
1. 从 gin.Context 读取 accountId、isAdmin。
2. 查询 sys_account，获取最新 is_admin、dept_id、status。
3. 如果 is_admin=1，返回超管上下文，不查询数据权限。
4. 查询 sys_account_role，拿到 role_id 列表。
5. 查询 sys_role.data_scope 和 sys_role_custom_dept。
6. 按第 11 节合并多角色数据权限。
7. 如果最终范围包含本部门，加入 sys_account.dept_id。
8. 如果最终范围包含本部门及下级，查询 sys_dept 子部门并合并。
9. 如果最终范围包含自定义部门，加入 sys_role_custom_dept.dept_id。
10. 返回 Context。
```

不建议把完整数据权限上下文放进 JWT，原因：

1. 角色数据权限修改后，旧 token 不会自动更新。
2. 自定义部门列表可能很长，JWT 体积不可控。
3. 账号部门变更后，旧 token 无法反映最新 `dept_id`。

可以做请求内缓存：

```go
ctx.Set("dataScopeContext", scopeCtx)
```

避免同一次请求重复查询。

## 14. 新增数据权限查询构建器

新增文件：

```text
internal/datascope/builder.go
```

建议支持两种场景：

1. GORM Gen wrapper 查询。
2. 原生 GORM / DAL SQL 查询。

### 14.1 通用选项

```go
type Options struct {
	HasDept         bool
	HasCreatedBy    bool
	DeptColumn       string
	CreatedByColumn  string
}
```

推荐默认：

```text
DeptColumn = "dept_id"
CreatedByColumn = "created_by"
```

### 14.2 GORM 原生条件构建

建议提供：

```go
func BuildWhere(scope *Context, opt Options) (string, []any)
```

规则：

| 数据范围 | 表有 `dept_id` | 表有 `created_by` | 条件 |
|---|---|---|---|
| 超管 | 任意 | 任意 | 不追加条件 |
| 全部数据 | 任意 | 任意 | 不追加条件 |
| 本部门 | 是 | 任意 | `dept_id IN ?` |
| 本部门 | 否 | 是 | 退化为 `created_by = ?` |
| 本部门及下级 | 是 | 任意 | `dept_id IN ?` |
| 本部门及下级 | 否 | 是 | 退化为 `created_by = ?` |
| 仅本人 | 任意 | 是 | `created_by = ?` |
| 仅本人 | 否 | 否 | `1 = 0` |
| 自定义部门 | 是 | 任意 | `dept_id IN ?` |
| 自定义部门 | 否 | 是 | 退化为 `created_by = ?` |

重要约束：

1. 需要部门列表但 `DeptIDs` 为空时，必须返回 `1 = 0`。
2. 表没有 `dept_id` 且没有 `created_by` 时，普通账号必须返回 `1 = 0`。
3. 不要因为构建失败而省略条件，否则会越权。

### 14.3 当前 Repository 接入方式

当前项目大量使用：

```go
FindPageByWrapper(ctx, req.PageNumber, req.PageSize, func(g gormplus.IGenWrapper[dao.ISysAccountEntityDo]) {
	g.WhereIf(...)
})
```

如果 `gormplus.IGenWrapper` 不支持安全追加 Raw 条件，建议在自定义 repository 里新增专用方法，通过原生 GORM 链式查询实现：

```go
func (r *SysAccountRepository) FindPageWithDataScope(
	ctx context.Context,
	req dto.GetSysAccountPageDTO,
	scope *datascope.Context,
) ([]*model.SysAccountEntity, int64, error)
```

Handler 示例：

```go
scopeCtx, err := s.DataScopeService.BuildContext(ctx)
if err != nil {
	s.Fail(ctx, err, "获取数据权限失败")
	return
}

list, total, err := s.SysAccountRepository.FindPageWithDataScope(ctx, req, scopeCtx)
```

Repository 内部追加：

```sql
WHERE deleted_at IS NULL
  AND (:keyword = '' OR username LIKE ? OR email LIKE ? OR mobile LIKE ?)
  AND (:status = 0 OR status = ?)
  AND <data_scope_condition>
```

## 15. 业务接口如何接入数据权限

数据权限不能只加在列表接口，必须覆盖：

```text
列表
详情
修改
删除
导出
批量操作
统计接口
```

错误示例：

```go
FindById(ctx, id)
```

如果接口面向普通管理员，必须改为：

```text
FindByIdWithDataScope(ctx, id, scopeCtx)
```

### 15.1 列表接口示例：账号列表

账号表本身已有 `dept_id`，因此账号列表可以直接按 `sys_account.dept_id` 过滤。

目标：

```go
scopeCtx, err := s.DataScopeService.BuildContext(ctx)
if err != nil {
	s.Fail(ctx, err, "获取数据权限失败")
	return
}

list, total, err := s.SysAccountRepository.FindPageWithDataScope(ctx, req, scopeCtx)
```

### 15.2 详情接口示例

当前：

```go
accountEntity, err := s.SysAccountRepository.FindById(ctx, id)
```

改为：

```go
scopeCtx, err := s.DataScopeService.BuildContext(ctx)
accountEntity, err := s.SysAccountRepository.FindByIdWithDataScope(ctx, id, scopeCtx)
```

如果查不到，统一返回：

```text
账号不存在或无权限访问
```

避免暴露数据存在性。

### 15.3 修改和删除接口示例

修改前先做数据权限校验：

```go
exists, err := s.SysAccountRepository.ExistsByIdWithDataScope(ctx, id, scopeCtx)
if err != nil {
	...
}
if !exists {
	s.Fail(ctx, errors.New("无权限访问该数据"), "账号不存在或无权限访问")
	return
}
```

通过校验后再执行更新或删除。

## 16. 哪些表需要接数据权限

当前项目核心表建议如下：

| 表 | 是否接数据权限 | 规则 |
|---|---|---|
| `sys_account` | 是 | 按 `dept_id` 或 `created_by` |
| `sys_role` | 谨慎 | 一般只超管或授权管理员维护；可按 `created_by` |
| `sys_resources` | 否 | 系统资源表，通常只超管维护 |
| `sys_dept` | 是 | 部门管理员只能看自己范围内部门 |
| `sys_post` | 是 | 按 `dept_id` |
| `sys_account_post` | 间接 | 通过账号或岗位关联校验 |
| `sys_role_custom_dept` | 否 | 权限配置表，不走普通数据权限；走接口功能权限 |
| `sys_account_role` | 否 | 权限配置关系表，不走普通数据权限；走接口功能权限 |
| `sys_role_resources` | 否 | 权限配置关系表，不走普通数据权限；走接口功能权限 |

原则：

1. 用户业务数据表优先接入。
2. 权限配置表不套普通数据权限，避免管理员无法维护权限。
3. 系统资源表只用功能权限控制。
4. 如果一个表没有 `dept_id`，但有 `created_by`，普通账号默认退化为仅本人。
5. 如果一个表没有 `dept_id` 且没有 `created_by`，普通账号默认无权限访问，除非接口明确只允许超管。

## 17. 缓存与权限变更

当前项目 repository 生成代码里有 SF + 缓存失效逻辑。数据权限接入后要注意：

1. 数据权限相关查询的缓存 key 必须包含 `accountId`、`dataScope`、`deptIds`、请求参数。
2. 角色 `data_scope` 修改后，要失效数据权限上下文缓存。
3. 角色 `sys_role_custom_dept` 修改后，要失效数据权限上下文缓存。
4. 账号 `dept_id` 修改后，要失效数据权限上下文缓存。
5. 如果后续把数据权限上下文缓存到 Redis，key 建议：

```text
datascope:account:{accountId}
```

需要删除缓存的操作：

```text
账号角色变更
账号部门变更
角色 data_scope 变更
角色自定义部门变更
部门层级变更
```

第一阶段可以不做跨请求缓存，只在单次请求 `gin.Context` 内缓存，简单且不会出现权限修改后旧缓存越权。

## 18. Casbin 与数据权限的边界

当前 Casbin 中间件负责：

```text
当前账号能不能访问某个接口
```

数据权限负责：

```text
当前账号访问接口后能操作哪些数据
```

不要把数据权限塞进 Casbin 策略，原因：

1. Casbin 当前模型是 `sub, obj, act`，适合 URL + Method。
2. 数据权限依赖业务表字段，如 `dept_id`、`created_by`，不是静态接口策略。
3. 列表、详情、修改、删除需要不同 SQL 条件，Casbin 中间件无法替业务查询追加 SQL。

正确顺序：

```text
AuthMiddleWare
  -> CasbinMiddleWare
      -> Handler
          -> DataScopeService.BuildContext
              -> Repository 查询追加数据权限条件
```

## 19. 测试清单

### 19.1 数据准备

准备部门树：

```text
1 总部
  2 技术部
    4 后端组
  3 运营部
```

准备角色：

| 角色 | data_scope | 自定义部门 |
|---|---:|---|
| 超管角色 | 1 | 空 |
| 部门主管 | 3 | 空 |
| 运营专员 | 2 | 空 |
| 普通员工 | 4 | 空 |
| 自定义角色 | 5 | 2、3 |

准备账号：

| 账号 | is_admin | dept_id | 角色 |
|---|---:|---:|---|
| admin | 1 | 1 | 超管角色 |
| manager | 2 | 2 | 部门主管 |
| operator | 2 | 3 | 运营专员 |
| staff | 2 | 4 | 普通员工 |
| custom | 2 | 1 | 自定义角色 |

准备业务数据，例如 `sys_account`：

| 数据 | dept_id | created_by |
|---|---:|---:|
| A | 1 | admin |
| B | 2 | manager |
| C | 3 | operator |
| D | 4 | staff |

### 19.2 列表验证

| 登录账号 | 预期 |
|---|---|
| admin | 看到全部 |
| manager | 看到 `dept_id=2、4` |
| operator | 看到 `dept_id=3` |
| staff | 只看到 `created_by=staff.id` |
| custom | 看到 `dept_id=2、3` |

### 19.3 详情验证

逐个访问不属于自己范围的数据 ID，应返回：

```text
数据不存在或无权限访问
```

### 19.4 修改和删除验证

对不属于自己范围的数据执行修改和删除，应返回：

```text
数据不存在或无权限访问
```

并确认数据库没有变化。

### 19.5 多角色验证

| 角色组合 | 预期 |
|---|---|
| 仅本人 + 本部门 | 本部门 |
| 本部门 + 自定义部门 3 | 当前账号部门 + 3 |
| 本部门及下级 + 自定义部门 3 | 当前账号部门及下级 + 3 |
| 任意角色 + 全部数据 | 全部数据 |
| 无有效角色 | 空范围或仅本人，按服务实现保持一致 |

## 20. 分阶段实施建议

### 阶段一：基础表和生成代码

1. 确认数据库已执行最新 `gin-admin-api.sql`。
2. 执行 `go run cmd/generator.go`。
3. 检查 `sys_account.dept_id`、`sys_role.data_scope`、`sys_role_custom_dept` 是否生成。
4. 新增 `pkg/enum/data_scope.enum.go`。

### 阶段二：账号部门接入

1. 账号创建 DTO 增加 `deptId`。
2. 账号修改 DTO 增加 `deptId`。
3. 账号 VO 增加 `deptId`。
4. 创建和修改账号时校验部门有效性。
5. 账号详情返回 `deptId`。

### 阶段三：角色数据权限接入

1. 角色 DTO 增加 `dataScope`、`deptIdList`。
2. 角色 VO 增加 `dataScope`、`deptIdList`。
3. 创建角色时写入 `sys_role.data_scope`。
4. 修改角色时更新 `sys_role.data_scope`。
5. `dataScope=5` 时维护 `sys_role_custom_dept`。
6. 删除角色时清理 `sys_role_custom_dept`。

### 阶段四：数据权限上下文和构建器

1. 新增 `internal/datascope/context.go`。
2. 新增 `internal/datascope/service.go`。
3. 新增 `internal/datascope/builder.go`。
4. 新增 repository 自定义查询方法。
5. 完成多角色合并规则。

### 阶段五：业务接口接入

优先接入：

1. `sys_account` 列表、详情、修改、删除。
2. `sys_dept` 列表、详情、修改、删除。
3. `sys_post` 列表、详情、修改、删除。

再接入后续业务表。

## 21. 关键注意事项

1. 数据权限是查询条件，不是前端过滤。
2. 列表、详情、修改、删除都要校验数据范围。
3. 默认权限必须收紧，不能因为查不到配置而放开。
4. 部门列表为空时必须返回空数据。
5. 权限配置表本身不要套普通数据权限。
6. 超管只跳过数据权限，不代表跳过登录认证。
7. Repository 缓存 key 必须包含数据权限上下文，否则不同用户可能命中同一份列表缓存。
8. 角色或账号权限变更后要考虑缓存失效。
9. Casbin 继续只管接口功能权限，数据权限在业务查询层实现。
10. 当前 SQL 的主设计是 `sys_account.dept_id`，不要再按旧文档从岗位反查账号部门。

## 22. 当前版本结论

本版本按当前 `gin-admin-api.sql` 重新收敛为：

1. 使用 `sys_account.dept_id` 保存账号所属部门。
2. 使用 `sys_role.data_scope` 保存角色数据范围。
3. 使用 `sys_role_custom_dept` 保存自定义部门范围。
4. `sys_dept.full_id` 用于计算本部门及下级。
5. 账号接口支持 `deptId`，不再要求通过岗位计算数据权限部门。
6. 角色接口支持 `dataScope`、`deptIdList`。
7. 登录 JWT 不扩展数据权限，只保留 `accountId` / `isAdmin`，运行时查询最新权限。
8. 数据权限上下文只在请求内缓存，第一阶段不做跨请求缓存。
9. Casbin 继续负责接口权限，数据权限由查询层负责。

这套方案与当前 SQL 表设计保持一致，并保持功能权限和数据权限分离：Casbin 负责接口是否可访问，数据权限在查询层统一追加过滤条件，空授权不放开，详情和写操作也必须校验数据范围。
