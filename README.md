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
  ```

  