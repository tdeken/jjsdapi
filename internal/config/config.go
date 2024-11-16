package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

var Conf = &Config{}

// FilePath 配置文件路径信息
type FilePath struct {
	ConfigName string // 设置配置文件名 (不需要后缀)
	ConfigType string // 设置配置文件类型
	ConfigPath string // 设置配置文件路径
}

// LoadConfig 从文件中加载配置
func LoadConfig(filepath FilePath) error {
	viper.SetConfigType(filepath.ConfigType) // 设置配置文件类型
	viper.SetConfigName(filepath.ConfigName) // 设置配置文件名 (不需要后缀)
	viper.AddConfigPath(filepath.ConfigPath) // 设置配置文件路径

	err := viper.ReadInConfig() // 读取配置文件
	if err != nil {
		return fmt.Errorf("failed to read etc file: %s", err)
	}

	err = viper.Unmarshal(Conf) // 将读取的配置映射到 Config 结构体中
	if err != nil {
		return fmt.Errorf("failed to unmarshal etc file: %s", err)
	}

	return nil
}

// Config 本地配置
type Config struct {
	Server   Server   `mapstructure:"server"` //系统配置
	Logger   Logger   `mapstructure:"logger"`
	Database Database `mapstructure:"database"`
	Redis    Redis    `mapstructure:"redis"`
	Alert    Alert    `mapstructure:"alert"`
}

type Server struct {
	Env  string `mapstructure:"env"`  // 动环境，取值为 "dev"、"test" 或 "prod"，默认为 "dev"
	Port int    `mapstructure:"port"` // 服务端口
}

// Logger defines the configuration for logging.
type Logger struct {
	StdOut     bool   `mapstructure:"std_out"`
	FileOut    bool   `mapstructure:"file_out"`
	Path       string `mapstructure:"path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// Database defines the configuration for the database.
type Database struct {
	URL             string        `mapstructure:"url"`               // 数据库连接地址
	MaxIdleConn     int           `mapstructure:"max_idle_conn"`     // 最大空闲连接数，默认为10
	MaxOpenConn     int           `mapstructure:"max_open_conn"`     // 最大活动连接数，默认为100
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"` // 连接最大生命周期(单位:秒) 默认为3600
	SlowThreshold   time.Duration `mapstructure:"slow_threshold"`    //慢sql警告阈值(单位:毫秒)
}

// Redis defines the configuration for Redis.
type Redis struct {
	Addr               string `mapstructure:"addr"`                 // Redis连接地址
	Password           string `mapstructure:"password"`             // Redis连接密码
	DB                 int    `mapstructure:"db"`                   // 连接库
	Sentinel           string `mapstructure:"sentinel"`             // 哨兵模式时 哨兵地址，没有则不填
	SentinelMasterName string `mapstructure:"sentinel_master_name"` // 哨兵模式下master名称，没有则不填
}

// Alert defines the configuration for alerting.
type Alert struct {
	FeiShu FeiShu         `mapstructure:"feishu"` // 飞书告警配置
	Scenes []FeiShuScenes `mapstructure:"scenes"` // 场景配置
}

type FeiShu struct {
	Url    string `mapstructure:"url"`    // 告警地址
	Secret string `mapstructure:"secret"` // 告警地址密码
}

type FeiShuScenes struct {
	Url    string `mapstructure:"url"`    // 告警地址
	Secret string `mapstructure:"secret"` // 告警地址密码
	Scene  string `mapstructure:"scene"`  //场景
}
