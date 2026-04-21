package repository

import (
	"context"
	"fmt"
	"gin-admin-api/internal/dal/dao"
	"gin-admin-api/internal/dal/model/entity"
	"gin-admin-api/internal/plugin"
	"gin-admin-api/pkg/enum"
	"time"

	"github.com/gin-gonic/gin"
)

type AccountRepository struct{}

type IAccountRepository interface {
	Create(ctx context.Context, username, password string) error
	GetByUsername(ctx context.Context, username string) (*entity.AccountEntity, error)
	GetByID(ctx context.Context, id int64) (*entity.AccountEntity, error)
	Delete(ctx context.Context, id int64) error
	UpdatePassword(ctx context.Context, id int64, password string) error
	UpdateStatus(ctx context.Context, id int64, status int64) error
	UpdateLoginInfo(ctx context.Context, id int64, ip string) error
	GetPage(ctx context.Context, username string, status int64, offset, limit int) ([]*entity.AccountEntity, int64, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}

func NewAccountRepository() IAccountRepository {
	return &AccountRepository{}
}

func (r *AccountRepository) Create(ctx context.Context, username, password string) error {
	fmt.Printf("Repository Create ctx type=%T, CtxOperatorKey=%v, CtxGinContextKey=%v\n",
		ctx, ctx.Value(plugin.CtxOperatorKey), ctx.Value(plugin.CtxGinContextKey))
	// 获取当前操作人
	operator := getOperatorFromContext(ctx)
	fmt.Println("Repository Create operator:", operator)
	return dao.AccountEntity.WithContext(ctx).
		Omit(dao.AccountEntity.LastLoginDate, dao.AccountEntity.LastLoginIP).
		Create(&entity.AccountEntity{
			Username: username,
			Password: password,
			Status:   enum.StatusNormalEnum,
			IsAdmin:  0,
		})
}

// getOperatorFromContext 从 context 中获取 operator
func getOperatorFromContext(ctx context.Context) int64 {
	if ctx == nil {
		return 0
	}
	// 尝试从 context value 获取
	if op, ok := ctx.Value(plugin.CtxOperatorKey).(int64); ok {
		return op
	}
	// 尝试从 gin context 获取
	if ginCtx, ok := ctx.Value(plugin.CtxGinContextKey).(*gin.Context); ok {
		if op, exists := ginCtx.Get("operator"); exists {
			if v, ok := op.(int64); ok {
				return v
			}
		}
	}
	return 0
}

func (r *AccountRepository) GetByUsername(ctx context.Context, username string) (*entity.AccountEntity, error) {
	return dao.AccountEntity.WithContext(ctx).
		Where(dao.AccountEntity.Username.Eq(username)).
		First()
}

func (r *AccountRepository) GetByID(ctx context.Context, id int64) (*entity.AccountEntity, error) {
	return dao.AccountEntity.WithContext(ctx).
		Where(dao.AccountEntity.ID.Eq(id)).
		First()
}

func (r *AccountRepository) Delete(ctx context.Context, id int64) error {
	if _, err := dao.AccountEntity.WithContext(ctx).
		Where(dao.AccountEntity.ID.Eq(id)).
		Delete(); err != nil {
		return err
	}
	return nil
}

func (r *AccountRepository) UpdatePassword(ctx context.Context, id int64, password string) error {
	if _, err := dao.AccountEntity.WithContext(ctx).
		Where(dao.AccountEntity.ID.Eq(id)).
		UpdateSimple(dao.AccountEntity.Password.Value(password)); err != nil {
		return err
	}
	return nil
}

func (r *AccountRepository) UpdateStatus(ctx context.Context, id int64, status int64) error {
	if _, err := dao.AccountEntity.WithContext(ctx).
		Where(dao.AccountEntity.ID.Eq(id)).
		UpdateSimple(dao.AccountEntity.Status.Value(status)); err != nil {
		return err
	}
	return nil
}

func (r *AccountRepository) UpdateLoginInfo(ctx context.Context, id int64, ip string) error {
	if _, err := dao.AccountEntity.WithContext(ctx).
		Where(dao.AccountEntity.ID.Eq(id)).
		Select(dao.AccountEntity.LastLoginDate, dao.AccountEntity.LastLoginIP).
		UpdateSimple(
			dao.AccountEntity.LastLoginDate.Value(time.Now()),
			dao.AccountEntity.LastLoginIP.Value(ip),
		); err != nil {
		return err
	}
	return nil
}

func (r *AccountRepository) GetPage(ctx context.Context, username string, status int64, offset, limit int) ([]*entity.AccountEntity, int64, error) {
	query := dao.AccountEntity.WithContext(ctx)
	if username != "" {
		query = query.Where(dao.AccountEntity.Username.Like("%" + username + "%"))
	}
	if status > 0 {
		query = query.Where(dao.AccountEntity.Status.Eq(status))
	}
	countQuery := query
	total, err := countQuery.Count()
	if err != nil {
		return nil, 0, err
	}
	list, err := query.Offset(offset).Limit(limit).Find()
	return list, total, err
}

func (r *AccountRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	result, err := dao.AccountEntity.WithContext(ctx).
		Where(dao.AccountEntity.Username.Eq(username)).
		Select(dao.AccountEntity.ID, dao.AccountEntity.Username).
		First()
	if err != nil {
		return false, err
	}
	return result != nil, nil
}
