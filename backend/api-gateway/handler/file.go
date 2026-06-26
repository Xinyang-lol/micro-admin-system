package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"micro-admin-system/backend/common/response"
	pb "micro-admin-system/backend/proto/gen"

	"github.com/gin-gonic/gin"
)

func (h *Handler) UploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请选择要上传的文件")
		return
	}
	if err := os.MkdirAll(h.uploadDir, 0755); err != nil {
		response.ServerError(c, "创建上传目录失败: "+err.Error())
		return
	}
	storedName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), cleanFileName(fileHeader.Filename))
	targetPath := filepath.Join(h.uploadDir, storedName)
	if err := c.SaveUploadedFile(fileHeader, targetPath); err != nil {
		response.ServerError(c, "保存文件失败: "+err.Error())
		return
	}

	userID, _ := c.Get("userID")
	req := &pb.UploadFileMetaRequest{
		OriginalName: fileHeader.Filename,
		StoredName:   storedName,
		Path:         targetPath,
		Size:         fileHeader.Size,
		ContentType:  fileHeader.Header.Get("Content-Type"),
		UploaderId:   userID.(int64),
	}
	ctx, cancel, cli, conn, ok := h.fileClient(c)
	if !ok {
		_ = os.Remove(targetPath)
		return
	}
	defer cancel()
	defer conn.Close()
	resp, err := cli.UploadFileMeta(ctx, req)
	if err != nil {
		_ = os.Remove(targetPath)
		grpcError(c, err)
		return
	}
	response.Success(c, resp)
}

func (h *Handler) ListFiles(c *gin.Context) {
	page, pageSize, keyword := pageQuery(c)
	ctx, cancel, cli, conn, ok := h.fileClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	resp, err := cli.ListFiles(ctx, &pb.FileListRequest{Page: page, PageSize: pageSize, Keyword: keyword})
	if err != nil {
		grpcError(c, err)
		return
	}
	response.Success(c, resp)
}

func (h *Handler) DownloadFile(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	ctx, cancel, cli, conn, ok := h.fileClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	info, err := cli.GetFile(ctx, &pb.FileIDRequest{Id: id})
	if err != nil {
		grpcError(c, err)
		return
	}
	if _, err := os.Stat(info.Path); err != nil {
		response.Error(c, 404, 404, "物理文件不存在")
		return
	}
	c.FileAttachment(info.Path, info.OriginalName)
}

func (h *Handler) DeleteFile(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	ctx, cancel, cli, conn, ok := h.fileClient(c)
	if !ok {
		return
	}
	defer cancel()
	defer conn.Close()
	info, err := cli.GetFile(ctx, &pb.FileIDRequest{Id: id})
	if err != nil {
		grpcError(c, err)
		return
	}
	if _, err := cli.DeleteFile(ctx, &pb.FileIDRequest{Id: id}); err != nil {
		grpcError(c, err)
		return
	}
	_ = os.Remove(info.Path)
	response.Success(c, gin.H{})
}

func cleanFileName(name string) string {
	name = filepath.Base(name)
	replacer := strings.NewReplacer("\\", "_", "/", "_", ":", "_", "*", "_", "?", "_", "\"", "_", "<", "_", ">", "_", "|", "_", " ", "_")
	return replacer.Replace(name)
}
