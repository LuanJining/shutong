package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	model "gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 状态码颜色
		statusColor := param.StatusCodeColor()
		methodColor := param.MethodColor()
		resetColor := param.ResetColor()

		// 根据状态码选择日志级别
		var level string
		if param.StatusCode >= 500 {
			level = "ERROR"
		} else if param.StatusCode >= 400 {
			level = "WARN"
		} else {
			level = "INFO"
		}

		// 格式化日志输出
		return fmt.Sprintf("[%s] %s %s %s%s%s %s %s%d%s %s %s\n",
			level,
			param.TimeStamp.Format("2006/01/02 15:04:05"),
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.Request.UserAgent(),
		)
	})
}

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-ID")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// FetchUserFromHeader 从请求头中获取用户信息
// Gateway 会在请求头中设置 X-User-ID，下游服务从中获取用户信息
func FetchUserFromHeader(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
			c.Abort()
			return
		}

		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			c.Abort()
			return
		}

		var user model.User
		if err := db.Preload("Roles").First(&user, uint(userID)).Error; err != nil {
			// 区分记录不存在和其他数据库错误
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询失败"})
			}
			c.Abort()
			return
		}

		// 检查用户状态
		if user.Status != 1 {
			c.JSON(http.StatusForbidden, gin.H{"error": "用户已被禁用"})
			c.Abort()
			return
		}

		c.Set("user", &user)
		c.Next()
	}
}
