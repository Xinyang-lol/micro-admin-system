package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"micro-admin-system/backend/common/response"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "ID 参数不正确")
		return 0, false
	}
	return id, true
}

func pageQuery(c *gin.Context) (int32, int32, string) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", c.DefaultQuery("page_size", "10")))
	keyword := c.DefaultQuery("keyword", "")
	return int32(page), int32(pageSize), keyword
}

func timeoutContext(c *gin.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), timeout)
}

func grpcError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	st, ok := status.FromError(err)
	if !ok {
		response.ServerError(c, "服务异常: "+err.Error())
		return
	}
	switch st.Code() {
	case codes.InvalidArgument:
		response.Error(c, http.StatusBadRequest, 400, st.Message())
	case codes.Unauthenticated:
		response.Error(c, http.StatusUnauthorized, 401, st.Message())
	case codes.PermissionDenied:
		response.Error(c, http.StatusForbidden, 403, st.Message())
	case codes.NotFound:
		response.Error(c, http.StatusNotFound, 404, st.Message())
	case codes.DeadlineExceeded, codes.Unavailable:
		response.Error(c, http.StatusServiceUnavailable, 503, st.Message())
	case codes.FailedPrecondition:
		response.Error(c, http.StatusBadRequest, 409, st.Message())
	default:
		response.ServerError(c, st.Message())
	}
}

func bindJSON(c *gin.Context, target any) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return false
	}
	return true
}

func internalOrGRPC(c *gin.Context, err error) {
	if err == nil {
		return
	}
	if _, ok := status.FromError(err); ok {
		grpcError(c, err)
		return
	}
	if errors.Is(err, context.DeadlineExceeded) {
		response.Error(c, http.StatusServiceUnavailable, 503, "服务调用超时")
		return
	}
	response.ServerError(c, err.Error())
}
