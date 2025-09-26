package config

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
	Database DatabaseConfig `mapstructure:"database"`
	Minio    MinioConfig    `mapstructure:"minio"`
	Qdrant   QdrantConfig   `mapstructure:"qdrant"`
	OCR      OCRConfig      `mapstructure:"ocr"`
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

type GinConfig struct {
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type MinioConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
}

type QdrantConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type OCRConfig struct {
	Url string `mapstructure:"url"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`        // 日志级别: debug, info, warn, error
	DBLogLevel string `mapstructure:"db_log_level"` // 数据库日志级别: silent, error, warn, info
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
	v.SetDefault("qdrant.host", "localhost")
	v.SetDefault("qdrant.port", "6333")
	v.SetDefault("qdrant.user", "qdrant")
	v.SetDefault("qdrant.password", "qdrant")
	v.SetDefault("qdrant.dbname", "kb_platform")
	v.SetDefault("qdrant.sslmode", "disable")
	v.SetDefault("ocr.url", "http://localhost:8084")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.db_log_level", "warn")
}
