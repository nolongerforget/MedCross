package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"medcross/utils"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证格式无效"})
			c.Abort()
			return
		}

		// 获取令牌部分
		tokenString := parts[1]

		// 解析JWT令牌
		token, err := parseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("无效的令牌: %v", err)})
			c.Abort()
			return
		}

		// 验证令牌
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 将用户信息存储在上下文中
			c.Set("userID", claims["sub"])
			c.Set("userName", claims["name"])
			c.Set("userRole", claims["role"])
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌声明"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RoleMiddleware 角色验证中间件
func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户角色
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
			c.Abort()
			return
		}

		// 检查用户是否有所需角色
		hasRole := false
		for _, role := range roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// 使用utils包中的ParseJWT函数
func parseJWT(tokenString string) (*jwt.Token, error) {
	return utils.ParseJWT(tokenString)
}
