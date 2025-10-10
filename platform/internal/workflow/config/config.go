package config

import (
	"fmt"
	"os"
	"strings"

	commonConfig "gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/config"
	"github.com/spf13/viper"
)

// Config Workflow 服务配置
type Config struct {
	Server   commonConfig.ServerConfig   `mapstructure:"server"`
	Gin      commonConfig.GinConfig      `mapstructure:"gin"`
	Database commonConfig.DatabaseConfig `mapstructure:"database"`
	Log      commonConfig.LogConfig      `mapstructure:"log"`
	Iam      commonConfig.IamConfig      `mapstructure:"iam"`
}

func Load() (*Config, error) {
	v := viper.New()

	// 设置配置文件名称和路径
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./internal/workflow/config")
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

	// 日志配置
	v.BindEnv("log.level", "KBASE_LOG_LEVEL", "LOG_LEVEL")
	v.BindEnv("log.db_log_level", "KBASE_DB_LOG_LEVEL", "DB_LOG_LEVEL")
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", "8082")
	v.SetDefault("gin.mode", "debug")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "5432")
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "password")
	v.SetDefault("database.dbname", "kb_platform")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.db_log_level", "warn")
}
