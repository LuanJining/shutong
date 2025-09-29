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
	Database DatabaseConfig `mapstructure:"database"`
	Minio    MinioConfig    `mapstructure:"minio"`
	Qdrant   QdrantConfig   `mapstructure:"qdrant"`
	OCR      OCRConfig      `mapstructure:"ocr"`
	Gin      GinConfig      `mapstructure:"gin"`
	Log      LogConfig      `mapstructure:"log"`
	Kb       KbConfig       `mapstructure:"kb"`
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

type KbConfig struct {
	Url string `mapstructure:"url"`
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

	// JWT配置
	v.BindEnv("jwt.secret", "KBASE_JWT_SECRET", "JWT_SECRET")
	v.BindEnv("jwt.access_token_expire_time", "KBASE_JWT_ACCESS_TOKEN_EXPIRE_TIME", "JWT_ACCESS_TOKEN_EXPIRE_TIME")
	v.BindEnv("jwt.refresh_token_expire_time", "KBASE_JWT_REFRESH_TOKEN_EXPIRE_TIME", "JWT_REFRESH_TOKEN_EXPIRE_TIME")

	// 日志配置
	v.BindEnv("log.level", "KBASE_LOG_LEVEL", "LOG_LEVEL")
	v.BindEnv("log.db_log_level", "KBASE_DB_LOG_LEVEL", "DB_LOG_LEVEL")

	// KB配置
	v.BindEnv("kb.url", "KBASE_KB_URL", "KB_URL")

	// 工作流配置
	v.BindEnv("workflow.url", "KBASE_WORKFLOW_URL", "WORKFLOW_URL")

	// minio配置
	v.BindEnv("minio.endpoint", "KBASE_MINIO_ENDPOINT", "MINIO_ENDPOINT")
	v.BindEnv("minio.access_key", "KBASE_MINIO_ACCESS_KEY", "MINIO_ACCESS_KEY")
	v.BindEnv("minio.secret_key", "KBASE_MINIO_SECRET_KEY", "MINIO_SECRET_KEY")
	v.BindEnv("minio.bucket", "KBASE_MINIO_BUCKET", "MINIO_BUCKET")
	v.BindEnv("minio.region", "KBASE_MINIO_REGION", "MINIO_REGION")

	// Qdrant配置
	v.BindEnv("qdrant.host", "KBASE_QDRANT_HOST", "QDRANT_HOST")
	v.BindEnv("qdrant.port", "KBASE_QDRANT_PORT", "QDRANT_PORT")
	v.BindEnv("qdrant.user", "KBASE_QDRANT_USER", "QDRANT_USER")
	v.BindEnv("qdrant.password", "KBASE_QDRANT_PASSWORD", "QDRANT_PASSWORD")
	v.BindEnv("qdrant.dbname", "KBASE_QDRANT_DBNAME", "QDRANT_DBNAME")
	v.BindEnv("qdrant.sslmode", "KBASE_QDRANT_SSLMODE", "QDRANT_SSLMODE")

	// OCR配置
	v.BindEnv("ocr.url", "KBASE_OCR_URL", "OCR_URL")

	// Gin配置
	v.BindEnv("gin.mode", "KBASE_GIN_MODE", "GIN_MODE")

	// 日志配置
	v.BindEnv("log.level", "KBASE_LOG_LEVEL", "LOG_LEVEL")
	v.BindEnv("log.db_log_level", "KBASE_DB_LOG_LEVEL", "DB_LOG_LEVEL")
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", "8080")
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
	v.SetDefault("kb.url", "http://localhost:8083")
}
