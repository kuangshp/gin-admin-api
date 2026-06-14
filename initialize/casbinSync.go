package initialize

import (
	"fmt"
	"strings"

	"gin-admin-api/internal/dal/model"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

// SyncCasbinFromDB 从业务表全量重建 casbin_rule。
//
// 适用场景：后期接入 Casbin，需要把已有的 sys_role_resources 和 sys_account_role
// 关联关系反向同步到 casbin_rule 表，使既有账号和角色立即获得对应的接口权限。
//
// 同步过程：
//  1. 直接清空 casbin_rule 表，并重载 enforcer 让其内存与空表对齐
//  2. 从 sys_role_resources INNER JOIN sys_resources(type=3, status=1) 生成 p 规则
//     形如 (role_2, /api/v1/admin/account, GET)
//  3. 从 sys_account_role 生成 g 规则
//     形如 (user_5, role_2)
//  4. 调用 SavePolicy 兜底持久化
//
// 这是一次性迁移操作，正常情况下只需运行一次。后续业务通过角色/账号管理 API
// 持续维护 casbin_rule，不再需要此函数。
func SyncCasbinFromDB(db *gorm.DB, enforcer *casbin.Enforcer) error {
	if enforcer == nil {
		return fmt.Errorf("enforcer 未初始化")
	}
	if db == nil {
		return fmt.Errorf("db 未初始化")
	}

	// 1. 清空 casbin_rule（Casbin 的 RemoveFilteredPolicy 不允许空过滤条件，
	//    所以直接用 raw SQL 清表，再 LoadPolicy 让 enforcer 内存与 DB 对齐）
	if err := db.Exec("DELETE FROM casbin_rule").Error; err != nil {
		return fmt.Errorf("清空 casbin_rule 表失败: %w", err)
	}
	if err := enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("重载 enforcer 失败: %w", err)
	}

	// 2. 重建 p 规则（role -> 接口）
	pCount, err := rebuildPolicyRules(db, enforcer)
	if err != nil {
		return err
	}

	// 3. 重建 g 规则（user -> role）
	gCount, err := rebuildGroupingRules(db, enforcer)
	if err != nil {
		return err
	}

	// 4. 兜底持久化
	if err := enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("SavePolicy 失败: %w", err)
	}

	fmt.Printf("✔ Casbin 反向同步完成: %d 条 p 规则, %d 条 g 规则\n", pCount, gCount)
	return nil
}

// rebuildPolicyRules 从 sys_role_resources + sys_resources(type=3) 生成 p 规则。
// 返回成功写入的规则条数。
func rebuildPolicyRules(db *gorm.DB, enforcer *casbin.Enforcer) (int, error) {
	type row struct {
		RoleID int64  `gorm:"column:role_id"`
		URL    string `gorm:"column:url"`
		Method string `gorm:"column:method"`
	}

	var rows []row
	err := db.
		Table("sys_role_resources AS rr").
		Select("rr.role_id AS role_id, sr.url AS url, sr.method AS method").
		Joins("INNER JOIN sys_resources AS sr ON rr.resources_id = sr.id").
		Where("sr.resources_type = ? AND sr.status = ?", 3, 1).
		Where("sr.url <> '' AND sr.method <> ''").
		Where("rr.deleted_at IS NULL AND sr.deleted_at IS NULL").
		Scan(&rows).Error
	if err != nil {
		return 0, fmt.Errorf("查询角色资源关联失败: %w", err)
	}
	if len(rows) == 0 {
		return 0, nil
	}

	// 去重：method 大小写不一致或重复关联时，避免写入重复规则
	seen := make(map[string]struct{}, len(rows))
	rules := make([][]string, 0, len(rows))
	for _, r := range rows {
		url := strings.TrimSpace(r.URL)
		method := strings.ToUpper(strings.TrimSpace(r.Method))
		if url == "" || method == "" || r.RoleID <= 0 {
			continue
		}
		sub := fmt.Sprintf("role_%d", r.RoleID)
		key := sub + "|" + url + "|" + method
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		rules = append(rules, []string{sub, url, method})
	}

	if len(rules) == 0 {
		return 0, nil
	}
	if _, err := enforcer.AddPolicies(rules); err != nil {
		return 0, fmt.Errorf("写入 p 规则失败: %w", err)
	}
	return len(rules), nil
}

// rebuildGroupingRules 从 sys_account_role 生成 g 规则。
// 返回成功写入的规则条数。
func rebuildGroupingRules(db *gorm.DB, enforcer *casbin.Enforcer) (int, error) {
	var list []*model.SysAccountRoleEntity
	if err := db.Where("deleted_at IS NULL").Find(&list).Error; err != nil {
		return 0, fmt.Errorf("查询账号角色关联失败: %w", err)
	}
	if len(list) == 0 {
		return 0, nil
	}

	// 去重：同一账号同一角色可能因脏数据出现多次
	seen := make(map[string]struct{}, len(list))
	rules := make([][]string, 0, len(list))
	for _, item := range list {
		if item.AccountID <= 0 || item.RoleID <= 0 {
			continue
		}
		user := fmt.Sprintf("user_%d", item.AccountID)
		role := fmt.Sprintf("role_%d", item.RoleID)
		key := user + "|" + role
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		rules = append(rules, []string{user, role})
	}

	if len(rules) == 0 {
		return 0, nil
	}
	if _, err := enforcer.AddGroupingPolicies(rules); err != nil {
		return 0, fmt.Errorf("写入 g 规则失败: %w", err)
	}
	return len(rules), nil
}
