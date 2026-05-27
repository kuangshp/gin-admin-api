// Package repository SysRoleResourcesRepository
// 此文件由用户自定义，不会被代码生成器覆盖
// 如需扩展或覆盖基础方法，请在 SysRoleResourcesRepository 中定义
package repository

// SysRoleResourcesRepository SysRoleResources仓储层
// 嵌入 defaultSysRoleResourcesRepository
type customerSysRoleResourcesRepository struct {
	*defaultSysRoleResourcesRepository
}

// SysRoleResourcesRepository ===== 用户自定义扩展方法请在 SysRoleResourcesRepository 中添加 =====
// 如需覆盖基础方法，实现相同的方法签名即可
type SysRoleResourcesRepository interface {
	iDefaultSysRoleResourcesRepository
}

var _ SysRoleResourcesRepository = (*customerSysRoleResourcesRepository)(nil)

// NewSysRoleResourcesRepository 创建SysRoleResourcesRepository实例
func NewSysRoleResourcesRepository() SysRoleResourcesRepository {
	return &customerSysRoleResourcesRepository{
		defaultSysRoleResourcesRepository: newDefaultSysRoleResourcesRepository(),
	}
}

// ═══════════════════════════════════════════════════════════════════════════
//  自定义方法示例:如何用 SF + 缓存包装自定义查询(DAL / 原生 SQL / 复合查询)
// ═══════════════════════════════════════════════════════════════════════════
//
// 默认 Repository 里的 FindById / FindList / FindPage 等方法已经自动接入了
// SF + 缓存,无需手动包装。但下面这些场景需要自己写 SF 调用:
//
//   1. 跨表 / 复合查询(JOIN、子查询)
//   2. DAL 文件化 SQL 查询(gormplus.DALQuery / gormplus.DALQueryOne)
//   3. 原生 gorm 链式查询
//   4. 任何非生成代码里的自定义查询方法
//
// 模板:任何自定义查询都按下面这个套路包装即可
//
// ─────────────────────────────────────────────────────────────────────────
//  示例 1:POST 分页查询 + DAL + SF 缓存(推荐写法,DTO 直接灌入)
// ─────────────────────────────────────────────────────────────────────────
//
// // 1) 定义请求 DTO,json tag 决定 cache key 字段名
// type LoginActivityReq struct {
//     Days    int    `json:"days"`
//     UserId  int64  `json:"user_id,omitempty"`
//     Source  string `json:"source,omitempty"`
//     Page    int    `json:"page"`
//     Size    int    `json:"size"`
// }
//
// // 2) Repository 方法接收 DTO 整体,BuildArgsFromStruct 一行展开
// func (r *customerSysRoleResourcesRepository) FindLoginActivity(
//     ctx context.Context, dto LoginActivityReq,
// ) ([]*LoginActivityRow, error) {
//     return gormplus.SF(func() ([]*LoginActivityRow, error) {
//         return gormplus.DALQuery[*LoginActivityRow](ctx, "login_activity.sql",
//             dto.Days, dto.UserId, dto.Source, dto.Page, dto.Size)
//     },
//         "sys_role_resources.FindLoginActivity",       // fnName:表名.方法名
//         gormplus.BuildArgsFromStruct(dto),         // ← 整个 DTO 自动展开为 cache args
//         5*time.Minute,                             // TTL,0 等价于 SFNoCache
//     )
// }
//
// // 好处:DTO 新增字段(比如加 Status)只需加 json tag,cache key 自动变化、
// // 缓存自动隔离,不用动 BuildArgsFromStruct 那一行。
//
// ─────────────────────────────────────────────────────────────────────────
//  示例 2:少量参数手写 BuildArgs(没有 DTO 的简单场景)
// ─────────────────────────────────────────────────────────────────────────
//
// func (r *customerSysRoleResourcesRepository) CountByStatus(
//     ctx context.Context, status int,
// ) (int64, error) {
//     return gormplus.SF(func() (int64, error) {
//         return gormplus.DALQueryOne[int64](ctx, "count_by_status.sql", status)
//     },
//         "sys_role_resources.CountByStatus",
//         gormplus.BuildArgs("status", status),     // 参数少,手写更直白
//         30*time.Second,
//     )
// }
//
// ─────────────────────────────────────────────────────────────────────────
//  示例 3:用 gorm 原生链式 + SF 缓存
// ─────────────────────────────────────────────────────────────────────────
//
// func (r *customerSysRoleResourcesRepository) FindTopActive(
//     ctx context.Context, limit int,
// ) ([]*SysRoleResources, error) {
//     return gormplus.SF(func() ([]*SysRoleResources, error) {
//         var list []*SysRoleResources
//         err := gormplus.Query[*SysRoleResources](r.GetDB(), ctx).
//             Where("status = ? AND last_active_at > ?", 1, time.Now().Add(-7*24*time.Hour)).
//             OrderBy("last_active_at DESC").
//             Limit(limit).
//             Build().Find(&list)
//         return list, err
//     },
//         "sys_role_resources.FindTopActive",
//         gormplus.BuildArgs("limit", limit),
//         5*time.Minute,
//     )
// }
//
// ─────────────────────────────────────────────────────────────────────────
//  示例 4:不要缓存,仅合并并发(防止热点击穿)
// ─────────────────────────────────────────────────────────────────────────
//
// func (r *customerSysRoleResourcesRepository) GetUserBalance(
//     ctx context.Context, userId int64,
// ) (*BalanceVO, error) {
//     return gormplus.SFNoCache(func() (*BalanceVO, error) {
//         // 实时余额不缓存,但同一秒来 100 个请求只打 1 次 DB
//         // 注意:DALQueryOne 返回 *T,泛型参数填 BalanceVO 即可
//         return gormplus.DALQueryOne[BalanceVO](ctx, "balance.sql", userId)
//     },
//         "sys_role_resources.GetUserBalance",
//         gormplus.BuildArgs("user_id", userId),
//     )
// }
//
// ─────────────────────────────────────────────────────────────────────────
//  示例 5:写操作后主动失效自定义缓存
// ─────────────────────────────────────────────────────────────────────────
//
// func (r *customerSysRoleResourcesRepository) UpdateActiveStatus(
//     ctx context.Context, userId int64, active bool,
// ) error {
//     err := gormplus.DALExec(ctx, "update_active.sql", active, userId)
//     if err != nil {
//         return err
//     }
//     // 精确失效:args 必须和查询时完全一致
//     gormplus.SFInvalidate("sys_role_resources.GetUserBalance",
//         gormplus.BuildArgs("user_id", userId))
//
//     // 批量前缀失效:清掉所有列表/聚合类的自定义缓存
//     gormplus.SFInvalidatePrefixes([]string{
//         "sys_role_resources.FindLoginActivity",
//         "sys_role_resources.FindTopActive",
//         "sys_role_resources.CountByStatus",
//     })
//     return nil
// }
//
// ─────────────────────────────────────────────────────────────────────────
//  关键约定
// ─────────────────────────────────────────────────────────────────────────
//
// - fnName 用 "sys_role_resources.方法名" 格式,方便后续按表批量失效
// - args 必须包含所有影响查询结果的参数(WHERE 条件值、limit、orderBy 字段等),
//   否则不同查询会命中同一 cache key 串数据
// - DTO 字段务必标 json tag(规则:json tag 第一段做 key,"-" 跳过),
//   字段名一旦确定就别随便改,改了会导致老缓存失效
// - 写操作记得主动调 SFInvalidate / SFInvalidatePrefix(es) 失效缓存
