package middleware

import (
	"github.com/gin-gonic/gin"
	"my-go-blog/utils"
	"net/http"
	"strings"
)

// GetCurrentUserID 获取已经通过鉴权的用户 ID
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	value, exists := c.Get("userID")
	if !exists {
		return 0, false
	}

	userID, ok := value.(uint)

	if !ok {
		return 0, false
	}

	return userID, true
}
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"code":    401,
					"message": "缺少 Authorization 请求头",
				},
			)
			return
		}
		// 2. 检查 Bearer Token 格式
		parts := strings.Fields(authHeader)
		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"code":    401,
					"message": "Token 格式错误",
				},
			)
			return
		}
		userID, err := utils.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"code":    401,
					"message": "Token 无效或已过期",
				},
			)
			return
		}
		// 4. 将用户 ID 保存到 Gin Context
		c.Set("userID", userID)

		// 5. 验证通过，继续执行后面的 Handler
		c.Next()

	}
}
