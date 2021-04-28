package common

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path"
)

func init() {
	fmt.Println("配置文件")
	InitConfig()
}

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
