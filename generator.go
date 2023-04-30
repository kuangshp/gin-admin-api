package main

import (
	"fmt"
	"github.com/kuangshp/generate-model/converter"
	"github.com/kuangshp/generate-model/utils"
	"strings"
)

func main() {
	var tableName string
	fmt.Print("请输入表名:")
	if _, err := fmt.Scanln(&tableName); err != nil {
		return
	}
	var fileName = utils.Case2Camel(strings.ToUpper(tableName[0:1]) + tableName[1:]) // 转为首字目大写
	var err = converter.NewTable2Struct().
		SavePath(fmt.Sprintf("./model/%sEntity.go", fileName)).
		Dsn("root:123456@tcp(localhost:3306)/gin-admin-api?charset=utf8mb4").
		TagKey("gorm").              // orm
		EnableJsonTag(true).         // json
		RealNameMethod("TableName"). // 生成表名
		Table(tableName).            // 表名
		DateToTime(true).            // 对时间字段进行转换
		Config(&converter.T2tConfig{
			JsonTagToHump: true,
			StructEnd:     "Entity",
		}). // 配置tag转驼峰
		Run()
	if err != nil {
		fmt.Println("创建数据模型失败")
		return
	}
	fmt.Println("生成数据模型成功")
}
