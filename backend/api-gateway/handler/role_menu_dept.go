package handler

import (
	"micro-admin-system/backend/common/response"
	pb "micro-admin-system/backend/proto/gen"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListRoles(c *gin.Context) {
	page, pageSize, keyword := pageQuery(c)
	ctx, cancel, cli, conn, ok := h.userClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	resp, err := cli.ListRoles(ctx, &pb.ListRequest{Page: page, PageSize: pageSize, Keyword: keyword})
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, resp)
}

func (h *Handler) CreateRole(c *gin.Context) {
	req := &pb.RoleSaveRequest{}
	if !bindJSON(c, req) {
		return
	}
	ctx, cancel, cli, conn, ok := h.userClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	resp, err := cli.CreateRole(ctx, req)
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, resp)
}

func (h *Handler) UpdateRole(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	req := &pb.RoleSaveRequest{Id: id}
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
	_, err := cli.UpdateRole(ctx, req)
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

func (h *Handler) DeleteRole(c *gin.Context) {
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
	_, err := cli.DeleteRole(ctx, &pb.IDRequest{Id: id})
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

func (h *Handler) AssignRoleMenus(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	req := &pb.RoleMenusRequest{Id: id}
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
	_, err := cli.AssignRoleMenus(ctx, req)
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

func (h *Handler) MenuTree(c *gin.Context) {
	ctx, cancel, cli, conn, ok := h.userClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	resp, err := cli.ListMenus(ctx, &pb.Empty{})
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, resp.Items)
}

func (h *Handler) CreateMenu(c *gin.Context) {
	req := &pb.MenuSaveRequest{}
	if !bindJSON(c, req) {
		return
	}
	ctx, cancel, cli, conn, ok := h.userClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	resp, err := cli.CreateMenu(ctx, req)
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, resp)
}

func (h *Handler) UpdateMenu(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	req := &pb.MenuSaveRequest{Id: id}
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
	_, err := cli.UpdateMenu(ctx, req)
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

func (h *Handler) DeleteMenu(c *gin.Context) {
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
	_, err := cli.DeleteMenu(ctx, &pb.IDRequest{Id: id})
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

func (h *Handler) DeptTree(c *gin.Context) {
	ctx, cancel, cli, conn, ok := h.userClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	resp, err := cli.ListDepts(ctx, &pb.Empty{})
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, resp.Items)
}

func (h *Handler) CreateDept(c *gin.Context) {
	req := &pb.DeptSaveRequest{}
	if !bindJSON(c, req) {
		return
	}
	ctx, cancel, cli, conn, ok := h.userClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	resp, err := cli.CreateDept(ctx, req)
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, resp)
}

func (h *Handler) UpdateDept(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	req := &pb.DeptSaveRequest{Id: id}
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
	_, err := cli.UpdateDept(ctx, req)
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

func (h *Handler) DeleteDept(c *gin.Context) {
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
	_, err := cli.DeleteDept(ctx, &pb.IDRequest{Id: id})
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, gin.H{})
}
