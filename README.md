## 一、项目基本结构

## 二、需要安装的依赖包

* `gin`框架包

  ```properties
  go get -u github.com/gin-gonic/gin
  ```

* `gorm`数据库包

  ```properties
  go get -u gorm.io/gorm
  go get -u gorm.io/driver/mysql
  ```

* 数据校验的包

  ```properties
  go get github.com/go-playground/validator
  ```

* `token`认证的包

  ```properties
  go get -u github.com/dgrijalva/jwt-go
  ```

* 日志管理包

  ```properties
  go get -u github.com/sirupsen/logrus
  go get -u github.com/lestrrat-go/file-rotatelogs
  go get -u github.com/rifflock/lfshook
  ```

* 配置文件的包

  ```properties
  go get -u github.com/spf13/viper
  ```

## 三、项目配置文件

* 1、在`config/application.yml`文件中创建项目需要的配置参数

  ```yaml
  server:
    port: 9000
  # 数据库配置
  datasource:
    driverName: mysql
    host: localhost
    port: "3306"
    database: gin_admin_api
    username: root
    password: 123456
    charset: utf8mb4
    loc: Asia/Shanghai
  ```

* 2、在`main.go`中定义一个初始化配置的文件

  ```go
  // 初始化配置
  func InitConfig() {
  	workDir, _ := os.Getwd()
  	viper.SetConfigName("application")
  	viper.SetConfigType("yml")
  	viper.AddConfigPath(path.Join(workDir, "config"))
  	// 或者使用全路径
  	//viper.AddConfigPath(path.Join(workDir, "config/application.yml"))
  	err := viper.ReadInConfig()
  	if err != nil {
  		fmt.Print("获取配置文件错误")
  		panic(err)
  	}
  }
  ```

* 3、在`init`函数中调用初始化配置的文件

  ```go
  func init() {
  	InitConfig()
  }
  ```

* 4、测试配置文件是否成功

  ```go
  func main() {
  	router := gin.Default()
  	router.GET("/", func(c *gin.Context) {
  		c.JSON(http.StatusOK, gin.H{
  			"code": 1,
  		})
  	})
  	port := viper.GetString("server.port")
  	fmt.Println("当前端口", port)
  	if port != "" {
  		router.Run(":" + port)
  	} else {
  		router.Run()
  	}
  }
  ```

* 5、或者可以单独到`common/config`文件中

  ```go
  package common
  
  import (
  	"fmt"
  	"github.com/spf13/viper"
  	"os"
  	"path"
  )
  
  // 初始化配置
  func InitConfig() {
  	workDir, _ := os.Getwd()
  	viper.SetConfigName("application")
  	viper.SetConfigType("yml")
  	viper.AddConfigPath(path.Join(workDir, "config"))
  	// 或者使用全路径
  	//viper.AddConfigPath(path.Join(workDir, "config/application.yml"))
  	err := viper.ReadInConfig()
  	if err != nil {
  		fmt.Print("获取配置文件错误")
  		panic(err)
  	}
  }
  
  func init() {
  	InitConfig()
  }
  ```

  借用在`main.go`中引入的文件,那么初始化就会先执行`init`函数
  ```go
  import (
  	...
    // 这里表示编译的时候不需要,但是运行的时候需要，不加这行下面的mian函数中是不能获取到参数的
  	_ "gin_admin_api/common" // gin_admin_api是在go.mod里面配置的module gin_admin_api，一般与项目名称一致
  	...
  )
  
  func main() {
  	...
  	port := viper.GetString("server.port")
  	fmt.Println("当前端口", port)
  	...
  }
  ```
  

## 四、初始化`gorm`数据库连接工具

* 1、在`common/database`下配置数据库连接

  ```go
  package common
  
  import (
  	"fmt"
  	_ "github.com/go-sql-driver/mysql"
  	"github.com/spf13/viper"
  	"gorm.io/driver/mysql"
  	"gorm.io/gorm"
  	"gorm.io/gorm/logger"
  	"log"
  	"net/url"
  	"os"
  	"time"
  )
  
  var DB *gorm.DB
  
  func init() {
  	fmt.Println("数据库连接")
  	InitDB()
  }
  
  func InitDB() *gorm.DB {
  	// 从配置文件中获取参数
  	host := viper.GetString("datasource.host")
  	port := viper.GetString("datasource.port")
  	database := viper.GetString("datasource.database")
  	username := viper.GetString("datasource.username")
  	password := viper.GetString("datasource.password")
  	charset := viper.GetString("datasource.charset")
  	loc := viper.GetString("datasource.loc")
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
  	// 配置日志输出
  	newLogger := logger.New(
  		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
  		logger.Config{
  			SlowThreshold:             time.Second,   // 缓存日志时间
  			LogLevel:                  logger.Silent, // 日志级别
  			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
  			Colorful:                  false,         // Disable color
  		},
  	)
  	db, err := gorm.Open(mysql.Open(sqlStr), &gorm.Config{
  		Logger: newLogger,
  	})
  	if err != nil {
  		fmt.Println("打开数据库失败", err)
  		panic("打开数据库失败" + err.Error())
  	}
  	DB = db
  	return DB
  }
  
  // TODO 文档地址: https://gorm.io/zh_CN/docs/
  ```

* 2、在`model/Account.go`的数据模型

  ```go
  package model
  
  import (
  	"gorm.io/gorm"
  )
  
  type Account struct {
  	gorm.Model
  	UserName string `gorm:"type:varchar(50);column(username);not null;unique;comment:账号"`
  	Password string `gorm:"type:varchar(200);not null;comment:账号密码"`
  	Mobile   string `gorm:"varchar(11);not null;unique;comment:手机号码"`
  }
  ```

* 3、在`main.go`中测试创建的数据模型及数据库连接工具

  ```go
  func init()  {
    // 自动同步数据模型到数据表
  	common.DB.AutoMigrate(&model.Account{})
  }
  ```

* 4、查看数据库的数据表**这里默认会加上一个s上去，表示复数**，如果要重命名表名可以参考下面代码

  ```go
  // 在数据模型的实体类文件中
  
  // 自定义表名
  func (Account) TableName() string {
  	return "account"
  }
  ```

  