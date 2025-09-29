package configs

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Iam      IamConfig      `mapstructure:"iam"`
	Workflow WorkflowConfig `mapstructure:"workflow"`
	Kb       KbConfig       `mapstructure:"kb"`
	Gin      GinConfig      `mapstructure:"gin"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type IamConfig struct {
	Url string `mapstructure:"url"`
}

type WorkflowConfig struct {
	Url string `mapstructure:"url"`
}

type KbConfig struct {
	Url string `mapstructure:"url"`
}

type GinConfig struct {
	Mode string `mapstructure:"mode"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

func Load() (*Config, error) {
	v := viper.New()

	// 设置配置文件名称和路径
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./internal/gateway/configs")
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
	v.BindEnv("iam.url", "KBASE_IAM_URL", "IAM_URL")
	v.BindEnv("workflow.url", "KBASE_WORKFLOW_URL", "WORKFLOW_URL")
	v.BindEnv("kb.url", "KBASE_KB_URL", "KB_URL")
	v.BindEnv("gin.mode", "KBASE_GIN_MODE", "GIN_MODE")
	v.BindEnv("log.level", "KBASE_LOG_LEVEL", "LOG_LEVEL")
	v.BindEnv("log.db_log_level", "KBASE_LOG_DB_LOG_LEVEL", "LOG_DB_LOG_LEVEL")
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", "8080")
	v.SetDefault("iam.url", "http://iam-service:8081")
	v.SetDefault("workflow.url", "http://workflow-service:8082")
	v.SetDefault("kb.url", "http://kb-service:8083")
	v.SetDefault("gin.mode", "debug")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.db_log_level", "warn")
}
