package handler

import (
	"context"
	"time"

	"micro-admin-system/backend/api-gateway/client"
	"micro-admin-system/backend/common/response"
	pb "micro-admin-system/backend/proto/gen"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type closer interface {
	Close() error
}

type Handler struct {
	clients   *client.Clients
	timeout   time.Duration
	redis     *redis.Client
	uploadDir string
}

func New(clients *client.Clients, timeout time.Duration, redisClient *redis.Client, uploadDir string) *Handler {
	return &Handler{clients: clients, timeout: timeout, redis: redisClient, uploadDir: uploadDir}
}

func (h *Handler) userClient(c *gin.Context) (context.Context, context.CancelFunc, pb.UserServiceClient, closer, bool) {
	ctx, cancel := timeoutContext(c, h.timeout)
	cli, conn, err := h.clients.User(ctx)
	if err != nil {
		cancel()
		response.ServerError(c, "发现 user-service 失败: "+err.Error())
		return nil, nil, nil, nil, false
	}
	return ctx, cancel, cli, conn, true
}

func (h *Handler) deviceClient(c *gin.Context) (context.Context, context.CancelFunc, pb.DeviceServiceClient, closer, bool) {
	ctx, cancel := timeoutContext(c, h.timeout)
	cli, conn, err := h.clients.Device(ctx)
	if err != nil {
		cancel()
		response.ServerError(c, "发现 device-service 失败: "+err.Error())
		return nil, nil, nil, nil, false
	}
	return ctx, cancel, cli, conn, true
}

func (h *Handler) fileClient(c *gin.Context) (context.Context, context.CancelFunc, pb.FileServiceClient, closer, bool) {
	ctx, cancel := timeoutContext(c, h.timeout)
	cli, conn, err := h.clients.File(ctx)
	if err != nil {
		cancel()
		response.ServerError(c, "发现 file-service 失败: "+err.Error())
		return nil, nil, nil, nil, false
	}
	return ctx, cancel, cli, conn, true
}
