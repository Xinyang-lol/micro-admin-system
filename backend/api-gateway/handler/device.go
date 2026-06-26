package handler

import (
	"strconv"

	"micro-admin-system/backend/common/response"
	pb "micro-admin-system/backend/proto/gen"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListDevices(c *gin.Context) {
	page, pageSize, keyword := pageQuery(c)
	typeID, _ := strconv.ParseInt(c.DefaultQuery("typeId", c.DefaultQuery("type_id", "0")), 10, 64)
	status := c.DefaultQuery("status", "")
	ctx, cancel, cli, conn, ok := h.deviceClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	resp, err := cli.ListDevices(ctx, &pb.DeviceListRequest{Page: page, PageSize: pageSize, Keyword: keyword, TypeId: typeID, Status: status})
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, resp)
}

func (h *Handler) CreateDevice(c *gin.Context) {
	req := &pb.DeviceSaveRequest{}
	if !bindJSON(c, req) {
		return
	}
	ctx, cancel, cli, conn, ok := h.deviceClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	resp, err := cli.CreateDevice(ctx, req)
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, resp)
}

func (h *Handler) UpdateDevice(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	req := &pb.DeviceSaveRequest{Id: id}
	if !bindJSON(c, req) {
		return
	}
	req.Id = id
	ctx, cancel, cli, conn, ok := h.deviceClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	_, err := cli.UpdateDevice(ctx, req)
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

func (h *Handler) DeleteDevice(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	ctx, cancel, cli, conn, ok := h.deviceClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	_, err := cli.DeleteDevice(ctx, &pb.DeviceIDRequest{Id: id})
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

func (h *Handler) DeviceStatistics(c *gin.Context) {
	ctx, cancel, cli, conn, ok := h.deviceClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	resp, err := cli.GetDeviceStatistics(ctx, &pb.DeviceEmpty{})
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, resp)
}

func (h *Handler) ListDeviceTypes(c *gin.Context) {
	ctx, cancel, cli, conn, ok := h.deviceClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	resp, err := cli.ListDeviceTypes(ctx, &pb.DeviceEmpty{})
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, resp.Items)
}

func (h *Handler) CreateDeviceType(c *gin.Context) {
	req := &pb.DeviceTypeSaveRequest{}
	if !bindJSON(c, req) {
		return
	}
	ctx, cancel, cli, conn, ok := h.deviceClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	resp, err := cli.CreateDeviceType(ctx, req)
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, resp)
}

func (h *Handler) UpdateDeviceType(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	req := &pb.DeviceTypeSaveRequest{Id: id}
	if !bindJSON(c, req) {
		return
	}
	req.Id = id
	ctx, cancel, cli, conn, ok := h.deviceClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	_, err := cli.UpdateDeviceType(ctx, req)
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, gin.H{})
}

func (h *Handler) DeleteDeviceType(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	ctx, cancel, cli, conn, ok := h.deviceClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	_, err := cli.DeleteDeviceType(ctx, &pb.DeviceIDRequest{Id: id})
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, gin.H{})
}
