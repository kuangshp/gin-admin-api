package generator

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 代码生成器配置
type Config struct {
	// 数据库配置
	DBType     string `yaml:"db_type"`      // 数据库类型，如 "mysql"
	Host       string `yaml:"host"`         // 数据库地址
	Port       int    `yaml:"port"`         // 数据库端口
	Username   string `yaml:"username"`     // 数据库账号
	Password   string `yaml:"password"`    // 数据库密码
	Database   string `yaml:"database"`    // 数据库名

	// 生成路径配置
	OutPath      string `yaml:"out_path"`       // dao输出路径，如 "./query/dao"
	ModelPkgPath string `yaml:"model_pkg_path"` // model包路径，如 "./query/model"
	RepoPath     string `yaml:"repo_path"`      // repository输出路径，如 "./query/repository"
	ApiPath      string `yaml:"api_path"`       // api desc路径，如 "./apps/admin/desc"

	// 项目包名
	Package string `yaml:"package"` // 项目包名，如 "esim-api"
}

// LoadConfig 从YAML文件加载配置
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &cfg, nil
}