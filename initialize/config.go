package initialize

import (
	"fmt"
	"gin_admin_api/global"
	"github.com/spf13/viper"
	"os"
	"path"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	workDir, _ := os.Getwd()
	isDev := GetEnvInfo("IS_DEV")
	fmt.Println(workDir, "目录")
	configFileName := path.Join(workDir, "application.prod.yml")
	fmt.Println(configFileName, "文件")
	if isDev {
		configFileName = path.Join(workDir, "application.dev.yml")
	}
	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	err := v.Unmarshal(&global.ServerConfig)
	if err != nil {
		fmt.Println("读取配置失败")
	}
	fmt.Println(&global.ServerConfig)
}
