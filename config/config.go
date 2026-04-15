package config

// ServerConfig 顶级配置结构体，与 yml 文件对应
type ServerConfig struct {
	Port       int         `mapstructure:"port"`
	Mode       string      `mapstructure:"mode"`
	DataSource DataSource  `mapstructure:"dataSource"`
	Redis      RedisConfig `mapstructure:"redis"`
	Log        LogConfig   `mapstructure:"log"`
	JWT        JWTConfig   `mapstructure:"jwt"`
}

type DataSource struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type LogConfig struct {
	Level     string `mapstructure:"level"`
	Filename  string `mapstructure:"filename"`
	MaxSize   int    `mapstructure:"maxSize"`
	MaxAge    int    `mapstructure:"maxAge"`
	MaxBackup int    `mapstructure:"maxBackup"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"signingKey"`
	ExpireTime int    `mapstructure:"expireTime"`
}
