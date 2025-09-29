package database

import (
	"fmt"

	model "gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/models"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(cfg config.DatabaseConfig, logCfg config.LogConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	// 根据配置设置数据库日志级别
	var logLevel logger.LogLevel
	switch logCfg.DBLogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Warn // 默认只显示警告和错误
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logLevel),
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动外键约束创建
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 自动迁移数据库表结构
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	// 自动迁移所有模型
	// 注意：Space表由IAM服务管理，这里不创建
	return db.AutoMigrate(
		&model.Document{},
		&model.DocumentChunk{},
	)
}
