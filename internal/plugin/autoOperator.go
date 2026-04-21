package plugin

import (
	"context"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type contextKey string

const (
	CtxOperatorKey   contextKey = "operator"
	CtxGinContextKey contextKey = "ginContext"
)

type AutoOperatorPlugin struct{}

func (p *AutoOperatorPlugin) Name() string { return "AutoOperatorPlugin" }

func (p *AutoOperatorPlugin) Initialize(db *gorm.DB) error {
	if err := db.Callback().Create().Before("gorm:create").
		Register("auto_operator:before_create", autoBeforeCreate); err != nil {
		return err
	}
	if err := db.Callback().Update().Before("gorm:update").
		Register("auto_operator:before_update", autoBeforeUpdate); err != nil {
		return err
	}
	if err := db.Callback().Update().Before("gorm:update").
		Register("auto_operator:before_update_column", autoBeforeUpdateColumn); err != nil {
		return err
	}
	return nil
}

func autoBeforeCreate(tx *gorm.DB) {
	if !hasSchema(tx) {
		return
	}
	op := getOperator(tx.Statement.Context)
	setIfExists(tx, "CreatedBy", op)
	setIfExists(tx, "UpdatedBy", op)
}

func autoBeforeUpdate(tx *gorm.DB) {
	if !hasSchema(tx) || tx.Statement.SkipHooks {
		return
	}
	if tx.Statement.Schema.LookUpField("UpdatedBy") == nil {
		return
	}
	op := getOperator(tx.Statement.Context)

	// 普通 Updates/Update 路径用 SetColumn
	tx.Statement.SetColumn("UpdatedBy", op, true)

	// UpdateSimple 路径：clause.Set 已存在时直接追加
	injectIntoClauseSet(tx, "updated_by", op)
}

func autoBeforeUpdateColumn(tx *gorm.DB) {
	if !hasSchema(tx) || !tx.Statement.SkipHooks {
		return
	}
	if tx.Statement.Schema.LookUpField("UpdatedBy") == nil {
		return
	}
	op := getOperator(tx.Statement.Context)
	// SkipHooks 路径直接操作 Dest map
	if dest, ok := tx.Statement.Dest.(map[string]interface{}); ok {
		dest["updated_by"] = op
	}
}

// injectIntoClauseSet 在 clause.Set 已存在时把字段追加进去
// 专门处理 UpdateSimple 这种表达式更新路径
func injectIntoClauseSet(tx *gorm.DB, col string, value any) {
	setClause, ok := tx.Statement.Clauses["SET"]
	if !ok {
		return // 不是 UpdateSimple 路径，不需要处理
	}
	set, ok := setClause.Expression.(clause.Set)
	if !ok {
		return
	}
	// 避免重复注入
	for _, a := range set {
		if a.Column.Name == col {
			return
		}
	}
	set = append(set, clause.Assignment{
		Column: clause.Column{Name: col},
		Value:  value,
	})
	setClause.Expression = set
	tx.Statement.Clauses["SET"] = setClause
}

func hasSchema(tx *gorm.DB) bool {
	return tx.Statement != nil && tx.Statement.Schema != nil
}

func setIfExists(tx *gorm.DB, field string, value any) {
	if tx.Statement.Schema.LookUpField(field) != nil {
		tx.Statement.SetColumn(field, value, true)
	}
}

func getOperator(ctx context.Context) int64 {
	if ctx == nil {
		return 0
	}
	return cast.ToInt64(ctx.Value(CtxOperatorKey))
}
