package middleware

import (
	"strings"

	"go_blog/internal/pkg/auth"
	"go_blog/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	// AuthorizationHeader Authorization 头部名称
	AuthorizationHeader = "Authorization"
	// BearerPrefix Bearer 前缀
	BearerPrefix = "Bearer "
	// ContextUserIDKey 上下文中用户ID的key
	ContextUserIDKey = "user_id"
	// ContextUsernameKey 上下文中用户名的key
	ContextUsernameKey = "username"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization 头
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		// 检查 Bearer 前缀
		if !strings.HasPrefix(authHeader, BearerPrefix) {
			response.Unauthorized(c, "Token格式错误")
			c.Abort()
			return
		}

		// 提取 token
		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		if tokenString == "" {
			response.Unauthorized(c, "Token为空")
			c.Abort()
			return
		}

		// 解析 token
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			response.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUsernameKey, claims.Username)

		c.Next()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) uint {
	userID, exists := c.Get(ContextUserIDKey)
	if !exists {
		return 0
	}
	return userID.(uint)
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) string {
	username, exists := c.Get(ContextUsernameKey)
	if !exists {
		return ""
	}
	return username.(string)
}
