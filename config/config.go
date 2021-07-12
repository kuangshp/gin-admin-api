package config

// DataSourceConfig 数据库的结构体
type DataSourceConfig struct {
	DriverName string `mapstructure:"driverName"`
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Database   string `mapstructure:"database"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	Charset    string `mapstructure:"charset"`
	Loc        string `mapstructure:"loc"`
}

type ServerConfig struct {
	Port       int              `mapstructure:"port"`
	Salt       string           `mapstructure:"salt"`
	DataSource DataSourceConfig `mapstructure:"datasource"`
}
