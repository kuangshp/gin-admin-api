package main

import (
	"fmt"
	"os"

	"gin-admin-api/initialize"
	"gin-admin-api/internal/config"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 用法：
//
//	go run ./cmd/sync_casbin.go                       # 默认读 application.dev.yml
//	go run ./cmd/sync_casbin.go application.prod.yml  # 指定其它配置
//
// 作用：从 sys_role_resources / sys_account_role 全量重建 casbin_rule，
// 通常在后期接入 Casbin 时执行一次，把历史的角色-资源、账号-角色关联同步过去。
//
// 注意：会先清空 casbin_rule 中所有 p 和 g 规则，再从业务表重新写入。
func main() {
	configPath := "application.dev.yml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Errorf("无法读取配置文件 %s: %w", configPath, err))
	}

	var serverConfig config.ServerConfig
	if err = yaml.Unmarshal(data, &serverConfig); err != nil {
		panic(fmt.Errorf("配置解析失败: %w", err))
	}

	ds := serverConfig.DataSource
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		ds.Username, ds.Password, ds.Host, ds.Port, ds.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		panic(fmt.Errorf("数据库连接失败: %w", err))
	}

	enforcer, err := initialize.NewCasbin(db)
	if err != nil {
		panic(fmt.Errorf("Casbin 初始化失败: %w", err))
	}

	fmt.Println("▶ 开始反向同步 casbin_rule ...")
	if err := initialize.SyncCasbinFromDB(db, enforcer); err != nil {
		panic(fmt.Errorf("同步失败: %w", err))
	}
}
