package initialize

import (
	"fmt"
	"gin-admin-api/internal/query/dao"
	"gin-admin-api/internal/query/model/entity"
	"gin-admin-api/pkg/enum"
	"github.com/kuangshp/go-utils/k"
	"time"
)

func InitAccountDataWithDao() error {
	total, err := dao.AccountEntity.Count()
	if err != nil {
		return fmt.Errorf("查询账号数量失败: %w", err)
	}

	if total > 0 {
		fmt.Println("管理员账号已存在，跳过初始化")
		return nil
	}

	password, err := k.MakePassword("123456")
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	admin := &entity.AccountEntity{
		Username:      "admin",
		Password:      password,
		Name:          "超级管理员",
		IsAdmin:       enum.AdminAccount,
		Status:        enum.StatusNormalEnum,
		LastLoginDate: time.Now(),
	}

	if err := dao.AccountEntity.Create(admin); err != nil {
		return fmt.Errorf("创建默认账号失败: %w", err)
	}

	fmt.Println("默认超级管理员账号创建成功: admin/123456")
	return nil
}
