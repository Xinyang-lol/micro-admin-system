package handler

import (
	"micro-admin-system/backend/common/response"
	pb "micro-admin-system/backend/proto/gen"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListUsers(c *gin.Context) {
	page, pageSize, keyword := pageQuery(c)
	ctx, cancel, cli, conn, ok := h.userClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	resp, err := cli.ListUsers(ctx, &pb.ListRequest{Page: page, PageSize: pageSize, Keyword: keyword})
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, resp)
}

func (h *Handler) CreateUser(c *gin.Context) {
	req := &pb.UserSaveRequest{}
	if !bindJSON(c, req) {
		return
	}
	ctx, cancel, cli, conn, ok := h.userClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	resp, err := cli.CreateUser(ctx, req)
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, resp)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	req := &pb.UserSaveRequest{Id: id}
	if !bindJSON(c, req) {
		return
	}
	req.Id = id
	ctx, cancel, cli, conn, ok := h.userClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	_, err := cli.UpdateUser(ctx, req)
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	ctx, cancel, cli, conn, ok := h.userClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	_, err := cli.DeleteUser(ctx, &pb.IDRequest{Id: id})
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

func (h *Handler) UpdateUserStatus(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	req := &pb.UserStatusRequest{Id: id}
	if !bindJSON(c, req) {
		return
	}
	req.Id = id
	ctx, cancel, cli, conn, ok := h.userClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	_, err := cli.UpdateUserStatus(ctx, req)
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

func (h *Handler) ResetUserPassword(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	req := &pb.UserPasswordRequest{Id: id}
	if !bindJSON(c, req) {
		return
	}
	req.Id = id
	ctx, cancel, cli, conn, ok := h.userClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	_, err := cli.ResetPassword(ctx, req)
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

func (h *Handler) AssignUserRoles(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	req := &pb.UserRolesRequest{Id: id}
	if !bindJSON(c, req) {
		return
	}
	req.Id = id
	ctx, cancel, cli, conn, ok := h.userClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	_, err := cli.AssignUserRoles(ctx, req)
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, gin.H{})
}
