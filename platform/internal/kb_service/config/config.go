package config

import (
	"fmt"
	"os"
	"strings"

	commonConfig "gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/config"
	"github.com/spf13/viper"
)

// Config KB Service 配置
type Config struct {
	Server   commonConfig.ServerConfig   `mapstructure:"server"`
	Iam      IamConfig                   `mapstructure:"iam"`
	Workflow WorkflowConfig              `mapstructure:"workflow"`
	Database commonConfig.DatabaseConfig `mapstructure:"database"`
	Minio    MinioConfig                 `mapstructure:"minio"`
	Gin      commonConfig.GinConfig      `mapstructure:"gin"`
	Log      commonConfig.LogConfig      `mapstructure:"log"`
	OpenAI   OpenAIConfig                `mapstructure:"openai"`
}

// IamConfig IAM 服务配置（KB Service 特有）
type IamConfig struct {
	Url string `mapstructure:"url"`
}

// WorkflowConfig Workflow 服务配置（KB Service 特有）
type WorkflowConfig struct {
	Url string `mapstructure:"url"`
}

// MinioConfig Minio 配置（KB Service 特有）
type MinioConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
}

// OpenAIConfig OpenAI 配置（KB Service 特有）
type OpenAIConfig struct {
	ApiKey string `mapstructure:"api_key"`
	Url    string `mapstructure:"url"`
}

func Load() (*Config, error) {
	v := viper.New()

	// 设置配置文件名称和路径
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./internal/kb_service/config")
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

	// Gin配置
	v.BindEnv("gin.mode", "KBASE_GIN_MODE", "GIN_MODE")

	// Minio配置
	v.BindEnv("minio.endpoint", "KBASE_MINIO_ENDPOINT", "MINIO_ENDPOINT")
	v.BindEnv("minio.access_key", "KBASE_MINIO_ACCESS_KEY", "MINIO_ACCESS_KEY")
	v.BindEnv("minio.secret_key", "KBASE_MINIO_SECRET_KEY", "MINIO_SECRET_KEY")
	v.BindEnv("minio.bucket", "KBASE_MINIO_BUCKET", "MINIO_BUCKET")
	v.BindEnv("minio.region", "KBASE_MINIO_REGION", "MINIO_REGION")

	// OpenAI配置
	v.BindEnv("openai.api_key", "KBASE_OPENAI_API_KEY", "OPENAI_API_KEY")
	v.BindEnv("openai.url", "KBASE_OPENAI_URL", "OPENAI_URL")

	// Workflow配置
	v.BindEnv("workflow.url", "KBASE_WORKFLOW_URL", "WORKFLOW_URL")

	// Iam配置
	v.BindEnv("iam.url", "KBASE_IAM_URL", "IAM_URL")

	// Gin配置
	v.BindEnv("gin.mode", "KBASE_GIN_MODE", "GIN_MODE")
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", "8083")
	v.SetDefault("iam.url", "http://localhost:8081")
	v.SetDefault("workflow.url", "http://localhost:8082")
	v.SetDefault("gin.mode", "debug")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "5432")
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "password")
	v.SetDefault("database.dbname", "kb_platform")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("minio.endpoint", "localhost:9000")
	v.SetDefault("minio.access_key", "minioadmin")
	v.SetDefault("minio.secret_key", "minioadmin")
	v.SetDefault("minio.bucket", "kb-platform")
	v.SetDefault("minio.region", "us-east-1")

	v.SetDefault("log.level", "info")
	v.SetDefault("log.db_log_level", "warn")
	v.SetDefault("openai.api_key", "")
	v.SetDefault("openai.url", "https://api.deepseek.com/v1")
}
