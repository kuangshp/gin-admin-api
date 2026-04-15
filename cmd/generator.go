package main

import (
	"fmt"
	"os"

	"gin-admin-api/internal/config"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"strings"
)

// Case2Camel 下划线转驼峰(大驼峰)
func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

// LowerCamelCase 转换为小驼峰
func LowerCamelCase(name string) string {
	name = Case2Camel(name)
	return strings.ToLower(name[:1]) + name[1:]
}

func main() {
	// 读取配置文件
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

	username := serverConfig.DataSource.Username
	password := serverConfig.DataSource.Password
	hostname := serverConfig.DataSource.Host
	port := serverConfig.DataSource.Port
	database := serverConfig.DataSource.Database
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, hostname, port, database)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:      "./internal/dal/dao",          // curd代码的输出路径
		ModelPkgPath: "./internal/dal/model/entity", // model代码的输出路径

		Mode: gen.WithDefaultQuery | gen.WithoutContext | gen.WithQueryInterface,

		FieldNullable:     false,
		FieldCoverable:    false,
		FieldSignable:     false,
		FieldWithIndexTag: false,
		FieldWithTypeTag:  true,
	})

	g.UseDB(db)

	dataMap := map[string]func(detailType gorm.ColumnType) (dataType string){
		"tinyint":   func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"smallint":  func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"mediumint": func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"bigint": func(detailType gorm.ColumnType) (dataType string) {
			if detailType.Name() == "version" {
				return "Version"
			}
			return "int64"
		},
		"int": func(detailType gorm.ColumnType) (dataType string) {
			if detailType.Name() == "version" {
				return "Version"
			}
			return "int64"
		},
		"int2":    func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"int4":    func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"integer": func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"json":    func(detailType gorm.ColumnType) (dataType string) { return "datatypes.JSON" },
		"numeric": func(detailType gorm.ColumnType) (dataType string) { return "decimal.Decimal" },
	}
	g.WithDataTypeMap(dataMap)

	g.WithModelNameStrategy(func(tableName string) (modelName string) {
		return Case2Camel(strings.ToUpper(tableName[:1]) + tableName[1:] + "Entity")
	})

	jsonField := gen.FieldJSONTagWithNS(func(columnName string) (tagContent string) {
		toStringField := `id, `
		if strings.Contains(toStringField, columnName) {
			return columnName + ",string"
		} else if strings.Contains(`deleted_at`, columnName) {
			return "-"
		}
		return LowerCamelCase(columnName)
	})

	autoUpdateTimeField := gen.FieldGORMTag("updated_at", func(tag field.GormTag) field.GormTag {
		return map[string][]string{
			"column":  {"updated_at"},
			"comment": {"更新时间"},
		}
	})

	autoCreateTimeField := gen.FieldGORMTag("created_at", func(tag field.GormTag) field.GormTag {
		return map[string][]string{
			"column":  {"created_at"},
			"comment": {"创建时间"},
		}
	})

	softDeleteField := gen.FieldType("deleted_at", "gorm.DeletedAt")
	fieldOpts := []gen.ModelOpt{jsonField, autoCreateTimeField, autoUpdateTimeField, softDeleteField}

	allModel := g.GenerateAllTable(fieldOpts...)
	g.ApplyBasic(allModel...)

	g.Execute()
}
