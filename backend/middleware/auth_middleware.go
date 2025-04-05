package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"d.com/MedCross/backend/config"
)

// AuthMiddleware 认证中间件，验证JWT令牌
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 检查Authorization格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证格式无效"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析JWT令牌
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
			}

			// 返回用于验证的密钥
			return []byte(config.GetConfig().JWT.Secret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌: " + err.Error()})
			c.Abort()
			return
		}

		// 验证令牌有效性
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 将用户信息添加到上下文
			c.Set("userID", claims["id"])
			c.Set("username", claims["username"])
			c.Set("role", claims["role"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
			c.Abort()
			return
		}
	}
}