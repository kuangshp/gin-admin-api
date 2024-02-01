package initialize

import (
	"fmt"
	"gin-admin-api/dao"
	"gin-admin-api/global"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"net/url"
	"strconv"
	"sync"
)

var (
	once sync.Once
)

func InitDB() {
	once.Do(func() {
		// 从配置文件中获取参数
		host := global.ServerConfig.DataSource.Host
		port := strconv.Itoa(global.ServerConfig.DataSource.Port)
		database := global.ServerConfig.DataSource.Database
		username := global.ServerConfig.DataSource.Username
		password := global.ServerConfig.DataSource.Password
		charset := global.ServerConfig.DataSource.Charset
		loc := global.ServerConfig.DataSource.Loc
		fmt.Println(host, port, database, username, password, "===aaa===>")
		// 字符串拼接
		sqlStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
			username,
			password,
			host,
			port,
			database,
			charset,
			url.QueryEscape(loc),
		)
		fmt.Println("数据库连接:", sqlStr)

		db, err := gorm.Open(mysql.Open(sqlStr), &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Info),
			DisableForeignKeyConstraintWhenMigrating: true, // 自动创建表的时候不创建外键
			NamingStrategy: schema.NamingStrategy{ // 自动创建表时候表名的配置
				SingularTable: true,
				// 全部的表名前面加前缀
				//TablePrefix: "mall_",
			},
		})
		if err != nil {
			fmt.Println("打开数据库失败", err)
			panic("打开数据库失败" + err.Error())
		}
		dao.SetDefault(db)
		//global.DB = dao.Use(db)
		global.DB = db
	})

}

// TODO 文档地址: https://gorm.io/zh_CN/docs/
