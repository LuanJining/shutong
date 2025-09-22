package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

// Server holds HTTP server configuration.
type Server struct {
	Host string
	Port int
}

// Database holds SQL database configuration.
type Database struct {
	DSN string
}

// Cache holds cache configuration (e.g. Redis DSN).
type Cache struct {
	Addr string
}

// Auth holds authentication configuration.
type Auth struct {
	JWTSigningKey string
	JWTTTLSeconds int
}

// Workflow holds workflow specific configuration values.
type Workflow struct {
	DefaultApproverRole string
}

// App aggregates configuration sections used by services.
type App struct {
	Name     string
	Env      string
	Server   Server
	Database Database
	Cache    Cache
	Auth     Auth
	Workflow Workflow
}

// YAMLConfig represents the structure of the YAML configuration file
type YAMLConfig struct {
	App      YAMLApp      `yaml:"app"`
	Server   YAMLServer   `yaml:"server"`
	Database YAMLDatabase `yaml:"database"`
	Cache    YAMLCache    `yaml:"cache"`
	Auth     YAMLAuth     `yaml:"auth"`
	Workflow YAMLWorkflow `yaml:"workflow"`
}

type YAMLApp struct {
	Name string `yaml:"name"`
	Env  string `yaml:"env"`
}

type YAMLServer struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type YAMLDatabase struct {
	DSN string `yaml:"dsn"`
}

type YAMLCache struct {
	Addr string `yaml:"addr"`
}

type YAMLAuth struct {
	JWTSigningKey string `yaml:"jwt_signing_key"`
	JWTTTLSeconds int    `yaml:"jwt_ttl_seconds"`
}

type YAMLWorkflow struct {
	DefaultApproverRole string `yaml:"default_approver_role"`
}

// Load returns the application configuration resolved from YAML file (local) or environment variables (production).
func Load(prefix string) (App, error) {
	// Check if we're in local development mode
	env := getEnvOrDefault(prefix+"_APP_ENV", "development")
	
	if env == "development" || env == "local" {
		return loadFromYAML(prefix)
	}
	
	return loadFromEnv(prefix)
}

// loadFromYAML loads configuration from YAML file for local development
func loadFromYAML(prefix string) (App, error) {
	yamlPath := "configs/local_config.yaml"
	
	// Check if YAML file exists
	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		// Fallback to environment variables if YAML file doesn't exist
		return loadFromEnv(prefix)
	}
	
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return App{}, fmt.Errorf("failed to read YAML config file: %w", err)
	}
	
	var yamlConfig YAMLConfig
	if err := yaml.Unmarshal(data, &yamlConfig); err != nil {
		return App{}, fmt.Errorf("failed to parse YAML config: %w", err)
	}
	
	// Convert YAML config to App struct
	app := App{
		Name: yamlConfig.App.Name,
		Env:  yamlConfig.App.Env,
		Server: Server{
			Host: yamlConfig.Server.Host,
			Port: yamlConfig.Server.Port,
		},
		Database: Database{
			DSN: yamlConfig.Database.DSN,
		},
		Cache: Cache{
			Addr: yamlConfig.Cache.Addr,
		},
		Auth: Auth{
			JWTSigningKey: yamlConfig.Auth.JWTSigningKey,
			JWTTTLSeconds: yamlConfig.Auth.JWTTTLSeconds,
		},
		Workflow: Workflow{
			DefaultApproverRole: yamlConfig.Workflow.DefaultApproverRole,
		},
	}
	
	if app.Auth.JWTSigningKey == "" {
		return App{}, fmt.Errorf("jwt_signing_key must be provided in YAML config")
	}
	
	return app, nil
}

// loadFromEnv loads configuration from environment variables for production
func loadFromEnv(prefix string) (App, error) {
	app := App{
		Name: getEnvOrDefault(prefix+"_APP_NAME", "knowledge-base"),
		Env:  getEnvOrDefault(prefix+"_APP_ENV", "development"),
		Server: Server{
			Host: getEnvOrDefault(prefix+"_SERVER_HOST", "0.0.0.0"),
			Port: getEnvAsInt(prefix+"_SERVER_PORT", 8080),
		},
		Database: Database{
			DSN: getEnvOrDefault(prefix+"_DATABASE_DSN", ""),
		},
		Cache: Cache{
			Addr: getEnvOrDefault(prefix+"_CACHE_ADDR", ""),
		},
		Auth: Auth{
			JWTSigningKey: getEnvOrDefault(prefix+"_AUTH_JWT_SIGNING_KEY", "dev-secret"),
			JWTTTLSeconds: getEnvAsInt(prefix+"_AUTH_JWT_TTL_SECONDS", 3600),
		},
		Workflow: Workflow{
			DefaultApproverRole: getEnvOrDefault(prefix+"_WORKFLOW_DEFAULT_ROLE", "content_reviewer"),
		},
	}

	if app.Auth.JWTSigningKey == "" {
		return App{}, fmt.Errorf("%s_AUTH_JWT_SIGNING_KEY must be provided", prefix)
	}

	return app, nil
}

func getEnvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}
