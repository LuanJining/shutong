package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	Gin      GinConfig      `mapstructure:"gin"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type JWTConfig struct {
	Secret                 string `mapstructure:"secret"`
	AccessTokenExpireTime  int    `mapstructure:"access_token_expire_time"`  // 小时
	RefreshTokenExpireTime int    `mapstructure:"refresh_token_expire_time"` // 小时
}

type LogConfig struct {
	Level      string `mapstructure:"level"`        // 日志级别: debug, info, warn, error
	DBLogLevel string `mapstructure:"db_log_level"` // 数据库日志级别: silent, error, warn, info
}

type GinConfig struct {
	Mode string `mapstructure:"mode"`
}

func Load() (*Config, error) {
	v := viper.New()

	// 设置配置文件名称和路径
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./internal/iam/config")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

	// 设置环境变量前缀
	v.SetEnvPrefix("KBASE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 自动绑定环境变量
	v.AutomaticEnv()

	// 根据环境变量选择配置文件
	env := os.Getenv("KBASE_ENV")
	if env == "" {
		env = "localtest" // 默认使用本地测试配置
	}

	// 设置配置文件名称
	v.SetConfigName(env)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// 配置文件不存在时使用默认值
	}

	// 绑定环境变量到配置
	bindEnvVars(v)

	// 设置默认值
	setDefaults(v)

	// 解析配置到结构体
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func bindEnvVars(v *viper.Viper) {
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

	// JWT配置
	v.BindEnv("jwt.secret", "KBASE_JWT_SECRET", "JWT_SECRET")
	v.BindEnv("jwt.access_token_expire_time", "KBASE_JWT_ACCESS_TOKEN_EXPIRE_TIME", "JWT_ACCESS_TOKEN_EXPIRE_TIME")
	v.BindEnv("jwt.refresh_token_expire_time", "KBASE_JWT_REFRESH_TOKEN_EXPIRE_TIME", "JWT_REFRESH_TOKEN_EXPIRE_TIME")

	// 日志配置
	v.BindEnv("log.level", "KBASE_LOG_LEVEL", "LOG_LEVEL")
	v.BindEnv("log.db_log_level", "KBASE_DB_LOG_LEVEL", "DB_LOG_LEVEL")

	// Gin配置
	v.BindEnv("gin.mode", "KBASE_GIN_MODE", "GIN_MODE")
}

func setDefaults(v *viper.Viper) {
	// 服务器默认配置
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.port", "8081")

	// 数据库默认配置
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "5432")
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "password")
	v.SetDefault("database.dbname", "kb_platform")
	v.SetDefault("database.sslmode", "disable")

	// JWT默认配置
	v.SetDefault("jwt.secret", "your-secret-key")
	v.SetDefault("jwt.access_token_expire_time", 24)
	v.SetDefault("jwt.refresh_token_expire_time", 720)

	// 日志默认配置
	v.SetDefault("log.level", "info")
	v.SetDefault("log.db_log_level", "warn")

	// Gin默认配置
	v.SetDefault("gin.mode", "debug")
}
