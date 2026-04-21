package generator

import (
	"fmt"
)

// ExampleGenerate 演示如何使用代码生成器
func ExampleGenerate() {
	// 直接传递配置
	cfg := &Config{
		DBType:       "mysql",
		Host:         "localhost",
		Port:         3306,
		Username:     "root",
		Password:     "123456",
		Database:     "gin-admin-api",
		OutPath:      "./query/dao",
		ModelPkgPath: "./query/model",
		RepoPath:     "./query/repository",
		ApiPath:      "./apps/admin/desc",
		Package:      "esim-api",
	}

	err := Generate(cfg)
	if err != nil {
		fmt.Printf("生成失败: %v\n", err)
	}

	// 或者使用LoadConfig方法加载配置
	cfg1, err1 := LoadConfig("./config/generator.yaml")
	if err1 != nil {
		panic(err)
	}
	Generate(cfg1)
}
