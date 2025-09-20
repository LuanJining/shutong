package config

import (
	"fmt"
	"os"
	"strconv"
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

// Load returns the application configuration resolved from environment variables.
func Load(prefix string) (App, error) {
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
