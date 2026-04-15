package main

import (
	"fmt"
	"gin-admin-api/internal/dal/model"
	"os"

	"gin-admin-api/internal/config"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	configPath := "application.dev.yml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Errorf("cannot read config file: %w", err))
	}

	var serverConfig config.ServerConfig
	if err := yaml.Unmarshal(data, &serverConfig); err != nil {
		panic(fmt.Errorf("cannot parse config file: %w", err))
	}

	ds := serverConfig.DataSource
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		ds.Username, ds.Password, ds.Host, ds.Port, ds.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}

	fmt.Println("▶ 开始数据库迁移...")
	if err := db.AutoMigrate(model.GetAllModels()...); err != nil {
		panic(fmt.Errorf("数据库迁移失败: %w", err))
	}
	fmt.Println("✔ 数据库迁移完成")
}
