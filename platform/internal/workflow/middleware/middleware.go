package middleware

import (
	"errors"
	"net/http"
	"strconv"

	model "gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

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
