package middleware

import (
	"context"
	"strings"
	"time"

	"micro-admin-system/backend/api-gateway/client"
	"micro-admin-system/backend/common/auth"
	redisx "micro-admin-system/backend/common/redis"
	"micro-admin-system/backend/common/response"
	pb "micro-admin-system/backend/proto/gen"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const ClaimsKey = "claims"

func Auth(secret string, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenText := bearerToken(c.GetHeader("Authorization"))
		if tokenText == "" {
			response.Unauthorized(c, "未登录或 Token 为空")
			c.Abort()
			return
		}
		if redisClient != nil {
			exists, err := redisClient.Exists(c.Request.Context(), redisx.TokenBlacklistKey(tokenText)).Result()
			if err == nil && exists > 0 {
				response.Unauthorized(c, "Token 已退出登录")
				c.Abort()
				return
			}
		}
		claims, err := auth.ParseToken(secret, tokenText)
		if err != nil {
			response.Unauthorized(c, "Token 无效或已过期")
			c.Abort()
			return
		}
		c.Set(ClaimsKey, claims)
		c.Set("token", tokenText)
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func Permission(permission string, clients *client.Clients, timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, ok := c.Get(ClaimsKey)
		if !ok {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}
		claims := value.(*auth.Claims)
		if auth.HasPermission(claims, permission) {
			c.Next()
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		userClient, conn, err := clients.User(ctx)
		if err == nil {
			defer conn.Close()
			result, callErr := userClient.CheckPermission(ctx, &pb.PermissionRequest{UserId: claims.UserID, Permission: permission})
			if callErr == nil && result.Allowed {
				c.Next()
				return
			}
		}
		response.Forbidden(c, "权限不足: "+permission)
		c.Abort()
	}
}

func bearerToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.Fields(header)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return header
}
