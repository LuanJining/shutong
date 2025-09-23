package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/iam/model"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/iam/service"
	"github.com/gin-gonic/gin"
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
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// AuthRequired 认证中间件
func AuthRequired(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少Authorization头"})
			c.Abort()
			return
		}

		// 检查Bearer token格式
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token格式"})
			c.Abort()
			return
		}

		token := tokenParts[1]
		user, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user", user)
		c.Next()
	}
}

// RequireRole 角色权限检查中间件
func RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
			c.Abort()
			return
		}

		// 类型断言获取用户信息
		userModel, ok := user.(*model.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "用户信息格式错误"})
			c.Abort()
			return
		}

		// 检查用户是否具有指定角色
		hasRole := false
		for _, role := range userModel.Roles {
			if role.Name == roleName {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "权限不足，需要" + roleName + "角色"})
			c.Abort()
			return
		}

		c.Next()
	}
}
