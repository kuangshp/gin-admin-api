package main

import (
	"fmt"
	gormplus "github.com/kuangshp/gorm-plus"
)

func main() {
	if err := gormplus.Generate(&gormplus.GeneratorConfig{
		DBType:       "mysql",                   // 数据库类型,详见上方说明
		Host:         "localhost",               // 数据库地址(sqlite 忽略)
		Port:         3306,                      // 数据库端口(postgres 默认 5432,sqlserver 默认 1433,sqlite 忽略)
		Username:     "root",                    // 数据库账号(sqlite 忽略)
		Password:     "123456",                  // 数据库密码(sqlite 忽略)
		Database:     "gin-rbac",                // 数据库名(sqlite 时为文件路径,如 "./data.db")
		OutPath:      "internal/dal/dao",        // dao输出路径，如 "./query/dao"
		ModelPkgPath: "internal/dal/model",      // model包路径，如 "./query/model"
		RepoPath:     "internal/dal/repository", // repository输出路径，如 "./query/repository"
		VoPath:       "internal/dal/vo",         // vo输出路径，如 "./query/vo"
		DtoPath:      "internal/dal/dto",        // dto输出路径，如 "./apps/admin/dto"
		MapperPath:   "internal/dal/mapper",     // mapper输出路径，如 "./query/mapper"
		Package:      "gin-admin-api",           // 项目包名，如 "esim-api"
	}); err != nil {
		fmt.Println("生成模板错误")
	}
}
