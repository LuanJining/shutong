package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// ServerConfig 服务器配置（通用）
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

func (c *ServerConfig) GetHost() string {
	return c.Host
}

func (c *ServerConfig) GetPort() string {
	return c.Port
}

// DatabaseConfig 数据库配置（通用）
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// GinConfig Gin 配置（通用）
type GinConfig struct {
	Mode string `mapstructure:"mode"`
}

// LogConfig 日志配置（通用）
type LogConfig struct {
	Level      string `mapstructure:"level"`        // 日志级别: debug, info, warn, error
	DBLogLevel string `mapstructure:"db_log_level"` // 数据库日志级别: silent, error, warn, info
}

// LoaderOptions 配置加载选项
type LoaderOptions struct {
	ConfigPaths []string // 配置文件搜索路径
	EnvPrefix   string   // 环境变量前缀，默认 "KBASE"
	DefaultEnv  string   // 默认环境，默认 "localtest"
}

// IamConfig IAM 服务配置（通用）
type IamConfig struct {
	Url string `mapstructure:"url"`
}

// Load 通用配置加载函数
// configPaths: 配置文件搜索路径，例如 ["./internal/iam/config", "./config", "."]
// cfg: 配置结构体指针
func Load(options LoaderOptions, cfg interface{}) error {
	v := viper.New()

	// 设置配置文件名称和路径
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// 添加配置路径
	for _, path := range options.ConfigPaths {
		v.AddConfigPath(path)
	}

	// 设置环境变量前缀
	envPrefix := options.EnvPrefix
	if envPrefix == "" {
		envPrefix = "KBASE"
	}
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 自动绑定环境变量
	v.AutomaticEnv()

	// 根据环境变量选择配置文件
	env := os.Getenv(envPrefix + "_ENV")
	if env == "" {
		env = options.DefaultEnv
		if env == "" {
			env = "localtest" // 最终默认值
		}
	}

	// 设置配置文件名称
	v.SetConfigName(env)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
		// 配置文件不存在时使用默认值
	}

	// 解析配置到结构体
	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// BindCommonEnvVars 绑定通用环境变量
func BindCommonEnvVars(v *viper.Viper) {
	// 服务器配置
	v.BindEnv("server.host", "KBASE_SERVER_HOST", "SERVER_HOST")
	v.BindEnv("server.port", "KBASE_SERVER_PORT", "SERVER_PORT")

	// 数据库配置
	v.BindEnv("database.host", "KBASE_DATABASE_HOST", "DB_HOST")
	v.BindEnv("database.port", "KBASE_DATABASE_PORT", "DB_PORT")
	v.BindEnv("database.user", "KBASE_DATABASE_USER", "DB_USER")
	v.BindEnv("database.password", "KBASE_DATABASE_PASSWORD", "DB_PASSWORD")
	v.BindEnv("database.dbname", "KBASE_DATABASE_DBNAME", "DB_NAME")
	v.BindEnv("database.sslmode", "KBASE_DATABASE_SSLMODE", "DB_SSLMODE")

	// 日志配置
	v.BindEnv("log.level", "KBASE_LOG_LEVEL", "LOG_LEVEL")
	v.BindEnv("log.db_log_level", "KBASE_DB_LOG_LEVEL", "DB_LOG_LEVEL")

	// Gin配置
	v.BindEnv("gin.mode", "KBASE_GIN_MODE", "GIN_MODE")

	// IAM配置
	v.BindEnv("iam.url", "KBASE_IAM_URL", "IAM_URL")
}

// SetCommonDefaults 设置通用默认值
func SetCommonDefaults(v *viper.Viper, serverPort string) {
	// 服务器默认配置
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", serverPort)

	// 数据库默认配置
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "5432")
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "password")
	v.SetDefault("database.dbname", "kb_platform")
	v.SetDefault("database.sslmode", "disable")

	// 日志默认配置
	v.SetDefault("log.level", "info")
	v.SetDefault("log.db_log_level", "warn")

	// Gin默认配置
	v.SetDefault("gin.mode", "debug")

	// IAM默认配置
	v.SetDefault("iam.url", "http://localhost:8081")
}
