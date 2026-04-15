package initialize

import (
	"fmt"
	"gin-admin-api/config"

	"github.com/spf13/viper"
)

// NewConfig Wire Provider：读取配置文件，返回 *config.ServerConfig
// envString 来自 generator.go 的 flag 参数：dev | prod
func NewConfig(envString string) (*config.ServerConfig, error) {
	v := viper.New()
	v.SetConfigName(fmt.Sprintf("application.%s", envString))
	v.SetConfigType("yml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件 application.%s.yml 失败: %w", envString, err)
	}

	var cfg config.ServerConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 同步写入全局变量，保持与原项目兼容
	fmt.Printf("✔ 配置加载成功 [%s] port=%d\n", envString, cfg.Port)
	return &cfg, nil
}
