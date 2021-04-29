## 一、`gin`项目介绍

`gin`框架对于写`go`语言的人来说入门很简单，有点类似`python`中的`flask`框架，什么都需要自己去找第三方包，然后根据自己的经验来创建目录结构，对于没有项目经验的人来说，这点真的不如同类型的`beego`框架，已经有清晰的目录结构，有时候我们可以说`gin`仅仅是一个包，算不上框架。自由组装度比较灵活，这也提现了我们开发人员的经验重要性。如何更好的搭建`gin-api`项目是一个难事。

以下是本人根据后端经验采用`mvc`的思路搭建一套基本的`gin-api`框架。以供大家参考使用,大家觉得好点个👍

* [github源码地址](https://github.com/kuangshp/gin_admin_api)
* [码云地址](https://gitee.com/k_2021/gin_admin_api)

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


## 五、在`gin`中使用路由分组实现路由管理

* 1、创建一个`route`的文件夹，里面负责收集全部控制器下的路由

  ```go
  package route
  
  import (
  	"gin_admin_api/controller/account"
  	"gin_admin_api/controller/login"
  	"gin_admin_api/controller/register"
  	"gin_admin_api/middleware"
  	"github.com/gin-gonic/gin"
  )
  
  func CollectRoute(router *gin.Engine) {
  	// 创建账号路由分组,先忽视中间件的存在
  	accountGroup := router.Group("/account", middleware.AuthMiddleWare())
  	account.AccountRouter(accountGroup)
  	// 登录的路由
  	loginGroup := router.Group("/login")
  	login.LoginRouter(loginGroup)
   
  	registerGroup := router.Group("/register")
  	register.RegisterRouter(registerGroup)
  }
  ```

* 2、比如登录的路由

  ```go
  package login
  
  import (
  	"github.com/gin-gonic/gin"
  )
  
  func LoginRouter(router *gin.RouterGroup) {
  	router.POST("/", Login)
  }
  ```

* 3、在`main.go`中使用路由组

  ```go
  func main() {
  	router := gin.Default()
  	// 注册路由组
  	route.CollectRoute(router)
    ...
  }
  ```

## 六、使用数据校验实现用户注册

* 1、在控制器下创建一个`dto`的文件,专门用来接收前端传递过来的数据

  ```go
  package dto
  
  import (
  	"fmt"
  	"gin_admin_api/model"
  	"github.com/go-playground/validator"
  	"unicode/utf8"
  )
  var valildate *validator.Validate
  
  func init() {
  	valildate = validator.New()
  	valildate.RegisterValidation("checkName", CheckNameFunc)
  }
  
  //定义注册的结构体(前端需要发送的数据结构)
  type RegisterDto struct {
  	UserName string `validate:"required,checkName" json:"username"`
  	Password string `validate:"required" json:"password"`
  }
  
  // 自定义校验器校验用户名
  func CheckNameFunc(f validator.FieldLevel) bool {
  	count := utf8.RuneCountInString(f.Field().String())
  	if count >= 2 && count <= 12 {
  		return true
  	} else {
  		return false
  	}
  }
  
  // 定义校验数据的方法
  func ValidatorRegister(account RegisterDto) error {
  	err := valildate.Struct(account)
  	if err != nil {
  		// 输出校验错误 .(validator.ValidationErrors)是断言
  		for _, e := range err.(validator.ValidationErrors)[:1] {
  			fmt.Println("错误字段:", e.Field())
  			fmt.Println("错误的值:", e.Value())
  			fmt.Println("错误的tag:", e.Tag())
  		}
  		return err
  	} else {
  		return nil
  	}
  }
  ```

* 2、在控制器中实现将前端传递过来的数据插入到数据库中

  ```go
  // 用户注册账号
  func Register(c *gin.Context) {
  	// 1.获取前端传递过来的数据
  	var registerDto dto.RegisterDto
  	err := c.Bind(&registerDto)
  	if err != nil {
  		response.Fail(c, "解析前端传递的数据错误")
  		return
  	}
  	// 2.对前端传递过来的数据进行校验
  	err = dto.ValidatorRegister(registerDto)
  	if err != nil {
  		response.Fail(c, "数据校验错误")
  		return
  	}
  	// 3.将数据插入到数据库中
  	newPassword, err := utils.GeneratePassword(registerDto.Password)
  	if err != nil {
  		response.Fail(c, "密码加密错误")
  		return
  	}
    // 4.组装成数据模型的数据结构
  	account := model.Account{
  		UserName: registerDto.UserName,
  		Password: newPassword,
  	}
  	tx := common.DB.Create(&account)
  	fmt.Println(tx.RowsAffected, tx.Error)
  	if tx.RowsAffected > 0 {
  		response.Success(c, nil)
  	} else {
  		response.Fail(c, "插入数据错误")
  	}
  }
  ```

* 3、关于密码加密和解密，可以参考`utils`里面的方法
* 4、查看数据库是否插入成功

## 七、关于中间件的使用

* 1、登录中间件可以参考文章[链接地址](https://blog.csdn.net/kuangshp128/article/details/116023080)
* 2、跨域中间件比较固定，可以直接百度或者参考我百度的数据
* 3、日志处理可以参考文档[链接地址](https://blog.csdn.net/kuangshp128/article/details/116230964)



