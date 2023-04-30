package initialize

import (
	"gin-admin-api/global"
	"gin-admin-api/utils"
	"fmt"
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
	isDev := utils.GetEnvInfo("IS_DEV")
	fmt.Println(workDir, "目录", isDev)
	configFileName := path.Join(workDir, "application.prod.yml")
	if isDev {
		configFileName = path.Join(workDir, "application.dev.yml")
	}
	fmt.Println(configFileName, "文件")
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

func GetDefaultEnv(key, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return defaultVal
}
