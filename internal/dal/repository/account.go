package repository

import (
	"context"
	"gin-admin-api/pkg/enum"
	"time"

	"gin-admin-api/internal/dal/dao"
	"gin-admin-api/internal/dal/model/entity"
)

type AccountRepository struct{}

type IAccountRepository interface {
	Create(ctx context.Context, username, password, salt string) error
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

func (r *AccountRepository) Create(ctx context.Context, username, password, salt string) error {
	return dao.AccountEntity.WithContext(ctx).
		Create(&entity.AccountEntity{
			Username: username,
			Password: password,
			Status:   enum.StatusNormalEnum,
			IsAdmin:  0,
		})
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
