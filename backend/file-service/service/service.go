package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	pb "micro-admin-system/backend/proto/gen"
	"micro-admin-system/backend/file-service/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FileService struct {
	repo    *repository.Repository
	auditCh chan string
}

func New(repo *repository.Repository) *FileService {
	s := &FileService{
		repo:    repo,
		auditCh: make(chan string, 128),
	}
	go s.consumeAudit()
	return s
}

func (s *FileService) consumeAudit() {
	for msg := range s.auditCh {
		log.Printf("[file-audit] %s", msg)
	}
}

func (s *FileService) UploadFileMeta(ctx context.Context, req *pb.UploadFileMetaRequest) (*pb.FileIDResponse, error) {
	if strings.TrimSpace(req.OriginalName) == "" || strings.TrimSpace(req.Path) == "" {
		return nil, status.Error(codes.InvalidArgument, "文件名和路径不能为空")
	}
	id, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "保存文件元数据失败: %v", err)
	}
	go func() {
		select {
		case s.auditCh <- time.Now().Format(time.RFC3339) + " uploaded " + req.OriginalName:
		default:
			log.Printf("[file-audit] channel full, drop log for %s", req.OriginalName)
		}
	}()
	return &pb.FileIDResponse{Id: id}, nil
}

func (s *FileService) ListFiles(ctx context.Context, req *pb.FileListRequest) (*pb.FileListResponse, error) {
	resp, err := s.repo.List(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询文件列表失败: %v", err)
	}
	return resp, nil
}

func (s *FileService) GetFile(ctx context.Context, req *pb.FileIDRequest) (*pb.FileInfo, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "文件 ID 不正确")
	}
	item, err := s.repo.Get(ctx, req.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "文件不存在")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询文件失败: %v", err)
	}
	return item, nil
}

func (s *FileService) DeleteFile(ctx context.Context, req *pb.FileIDRequest) (*pb.FileEmpty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "文件 ID 不正确")
	}
	if err := s.repo.Delete(ctx, req.Id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "文件不存在")
		}
		return nil, status.Errorf(codes.Internal, "删除文件失败: %v", err)
	}
	return &pb.FileEmpty{}, nil
}
