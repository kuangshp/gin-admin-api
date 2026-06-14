package initialize

import (
	"fmt"
	"gin-admin-api/internal/dal/dao"
	"gin-admin-api/internal/dal/model"
	"gin-admin-api/pkg/enum"
	"github.com/kuangshp/go-utils/k"
	"time"
)

func InitAccountDataWithDao() error {
	total, err := dao.SysAccountEntity.Count()
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

	admin := &model.SysAccountEntity{
		Username:      "admin",
		Password:      password,
		IsAdmin:       enum.AdminAccount,
		Status:        enum.StatusNormalEnum,
		LastLoginDate: time.Now(),
	}

	if err := dao.SysAccountEntity.Create(admin); err != nil {
		return fmt.Errorf("创建默认账号失败: %w", err)
	}

	fmt.Println("默认超级管理员账号创建成功: admin/123456")
	return nil
}

func InitResourcesDataWithDao() error {
	total, err := dao.SysResourcesEntity.Count()
	if err != nil {
		return fmt.Errorf("查询资源数量失败: %w", err)
	}

	if total > 0 {
		fmt.Println("系统资源已存在，跳过初始化")
		return nil
	}

	resources := []*model.SysResourcesEntity{
		{
			Title:         "系统管理",
			URL:           "system",
			Method:        "",
			Icon:          "system",
			ResourcesType: 1,
			IsCache:       1,
			IsHidden:      1,
			IsLink:        1,
			ParentID:      0,
			Sort:          100,
			Status:        enum.StatusNormalEnum,
			Description:   "",
			IsAdminHave:   0,
		},
		{
			Title:         "账号管理",
			URL:           "account",
			Method:        "",
			Icon:          "personnel-manage",
			ResourcesType: 2,
			IsCache:       1,
			IsHidden:      1,
			IsLink:        1,
			ParentID:      1,
			Sort:          1,
			Status:        enum.StatusNormalEnum,
			Description:   "",
			IsAdminHave:   0,
		},
		{
			Title:         "角色管理",
			URL:           "role",
			Method:        "",
			Icon:          "role",
			ResourcesType: 2,
			IsCache:       1,
			IsHidden:      1,
			IsLink:        1,
			ParentID:      1,
			Sort:          2,
			Status:        enum.StatusNormalEnum,
			Description:   "",
			IsAdminHave:   0,
		},
		{
			Title:         "资源管理",
			URL:           "resources",
			Method:        "",
			Icon:          "permission",
			ResourcesType: 2,
			IsCache:       1,
			IsHidden:      1,
			IsLink:        1,
			ParentID:      1,
			Sort:          3,
			Status:        enum.StatusNormalEnum,
			Description:   "",
			IsAdminHave:   0,
		},
	}

	if err := dao.SysResourcesEntity.Create(resources...); err != nil {
		return fmt.Errorf("创建默认系统资源失败: %w", err)
	}

	fmt.Println("默认系统资源创建成功")
	return nil
}
