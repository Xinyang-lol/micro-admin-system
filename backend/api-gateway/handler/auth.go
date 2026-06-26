package handler

import (
	"time"

	"micro-admin-system/backend/api-gateway/middleware"
	"micro-admin-system/backend/common/auth"
	redisx "micro-admin-system/backend/common/redis"
	"micro-admin-system/backend/common/response"
	pb "micro-admin-system/backend/proto/gen"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Login(c *gin.Context) {
	req := &pb.LoginRequest{}
	if !bindJSON(c, req) {
		return
	}
	ctx, cancel, cli, conn, ok := h.userClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	resp, err := cli.Login(ctx, req)
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, resp)
}

func (h *Handler) Logout(c *gin.Context) {
	token, _ := c.Get("token")
	ttl := 24 * time.Hour
	if value, ok := c.Get(middleware.ClaimsKey); ok {
		if claims, ok := value.(*auth.Claims); ok && claims.ExpiresAt != nil {
			if remain := time.Until(claims.ExpiresAt.Time); remain > 0 {
				ttl = remain
			}
		}
	}
	if h.redis != nil && token != nil {
		if err := h.redis.Set(c.Request.Context(), redisx.TokenBlacklistKey(token.(string)), "1", ttl).Err(); err != nil {
			response.ServerError(c, "退出登录失败: "+err.Error())
			return
		}
	}
	response.Success(c, gin.H{"message": "已退出登录"})
}

func (h *Handler) Profile(c *gin.Context) {
	userID, _ := c.Get("userID")
	ctx, cancel, cli, conn, ok := h.userClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	resp, err := cli.GetUserInfo(ctx, &pb.UserIDRequest{UserId: userID.(int64)})
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, resp)
}
