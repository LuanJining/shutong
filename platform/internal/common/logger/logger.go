package logger

import (
	"log"
	"strings"
	"sync"
)

// LogLevel 日志级别
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var (
	currentLevel LogLevel = INFO
	mu           sync.RWMutex
)

// SetLevel 设置日志级别
func SetLevel(level string) {
	mu.Lock()
	defer mu.Unlock()

	switch strings.ToLower(level) {
	case "debug":
		currentLevel = DEBUG
	case "info":
		currentLevel = INFO
	case "warn", "warning":
		currentLevel = WARN
	case "error":
		currentLevel = ERROR
	default:
		currentLevel = INFO
	}
}

// GetLevel 获取当前日志级别
func GetLevel() LogLevel {
	mu.RLock()
	defer mu.RUnlock()
	return currentLevel
}

// Debug 输出 DEBUG 级别日志
func Debug(v ...any) {
	if GetLevel() <= DEBUG {
		log.Println(v...)
	}
}

// Debugf 格式化输出 DEBUG 级别日志
func Debugf(format string, v ...any) {
	if GetLevel() <= DEBUG {
		log.Printf(format, v...)
	}
}

// Info 输出 INFO 级别日志
func Info(v ...any) {
	if GetLevel() <= INFO {
		log.Println(v...)
	}
}

// Infof 格式化输出 INFO 级别日志
func Infof(format string, v ...any) {
	if GetLevel() <= INFO {
		log.Printf(format, v...)
	}
}

// Warn 输出 WARN 级别日志
func Warn(v ...any) {
	if GetLevel() <= WARN {
		log.Println(v...)
	}
}

// Warnf 格式化输出 WARN 级别日志
func Warnf(format string, v ...any) {
	if GetLevel() <= WARN {
		log.Printf(format, v...)
	}
}

// Error 输出 ERROR 级别日志
func Error(v ...any) {
	if GetLevel() <= ERROR {
		log.Println(v...)
	}
}

// Errorf 格式化输出 ERROR 级别日志
func Errorf(format string, v ...any) {
	if GetLevel() <= ERROR {
		log.Printf(format, v...)
	}
}

// Printf 兼容标准 log.Printf（使用 INFO 级别）
func Printf(format string, v ...any) {
	Infof(format, v...)
}

// Println 兼容标准 log.Println（使用 INFO 级别）
func Println(v ...any) {
	Info(v...)
}
