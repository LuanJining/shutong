package middleware

import (
	"net/http"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/gateway/handler"
	"github.com/gin-gonic/gin"
)

// AuthRequired Gateway 专用的认证中间件
// 通过 IAM 服务验证 token 并将用户 ID 添加到请求头中
func AuthRequired(iamHandler *handler.IamHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少Authorization头"})
			c.Abort()
			return
		}

		user, err := iamHandler.ValidateToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
