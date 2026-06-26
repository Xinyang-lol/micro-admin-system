package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	pb "micro-admin-system/backend/proto/gen"
	"micro-admin-system/backend/device-service/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeviceService struct {
	repo *repository.Repository
}

type statResult struct {
	name      string
	value     int64
	typeStats map[string]int64
	err       error
}

func New(repo *repository.Repository) *DeviceService {
	return &DeviceService{repo: repo}
}

func (s *DeviceService) ListDevices(ctx context.Context, req *pb.DeviceListRequest) (*pb.DeviceListResponse, error) {
	resp, err := s.repo.ListDevices(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询设备列表失败: %v", err)
	}
	return resp, nil
}

func (s *DeviceService) CreateDevice(ctx context.Context, req *pb.DeviceSaveRequest) (*pb.DeviceIDResponse, error) {
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Code) == "" {
		return nil, status.Error(codes.InvalidArgument, "设备名称和编码不能为空")
	}
	if req.Status == "" {
		req.Status = "offline"
	}
	id, err := s.repo.CreateDevice(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建设备失败: %v", err)
	}
	return &pb.DeviceIDResponse{Id: id}, nil
}

func (s *DeviceService) UpdateDevice(ctx context.Context, req *pb.DeviceSaveRequest) (*pb.DeviceEmpty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "设备 ID 不正确")
	}
	if err := s.repo.UpdateDevice(ctx, req); err != nil {
		return nil, mapRepoError(err, "更新设备失败")
	}
	return &pb.DeviceEmpty{}, nil
}

func (s *DeviceService) DeleteDevice(ctx context.Context, req *pb.DeviceIDRequest) (*pb.DeviceEmpty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "设备 ID 不正确")
	}
	if err := s.repo.DeleteDevice(ctx, req.Id); err != nil {
		return nil, mapRepoError(err, "删除设备失败")
	}
	return &pb.DeviceEmpty{}, nil
}

func (s *DeviceService) GetDeviceStatistics(ctx context.Context, _ *pb.DeviceEmpty) (*pb.DeviceStatisticsResponse, error) {
	resultCh := make(chan statResult, 5)
	run := func(name string, fn func(context.Context) (int64, error)) {
		go func() {
			value, err := fn(ctx)
			resultCh <- statResult{name: name, value: value, err: err}
		}()
	}
	run("total", s.repo.CountAll)
	run("online", func(ctx context.Context) (int64, error) { return s.repo.CountByStatus(ctx, "online") })
	run("offline", func(ctx context.Context) (int64, error) { return s.repo.CountByStatus(ctx, "offline") })
	run("repair", func(ctx context.Context) (int64, error) { return s.repo.CountByStatus(ctx, "repair") })
	go func() {
		stats, err := s.repo.CountByType(ctx)
		resultCh <- statResult{name: "type", typeStats: stats, err: err}
	}()

	resp := &pb.DeviceStatisticsResponse{TypeStats: map[string]int64{}}
	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done():
			return nil, status.Error(codes.DeadlineExceeded, "设备统计超时")
		case result := <-resultCh:
			if result.err != nil {
				return nil, status.Errorf(codes.Internal, "设备统计失败: %v", result.err)
			}
			switch result.name {
			case "total":
				resp.Total = result.value
			case "online":
				resp.Online = result.value
			case "offline":
				resp.Offline = result.value
			case "repair":
				resp.Repair = result.value
			case "type":
				resp.TypeStats = result.typeStats
			}
		}
	}
	return resp, nil
}

func (s *DeviceService) ListDeviceTypes(ctx context.Context, _ *pb.DeviceEmpty) (*pb.DeviceTypeListResponse, error) {
	resp, err := s.repo.ListDeviceTypes(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询设备类型失败: %v", err)
	}
	return resp, nil
}

func (s *DeviceService) CreateDeviceType(ctx context.Context, req *pb.DeviceTypeSaveRequest) (*pb.DeviceIDResponse, error) {
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Code) == "" {
		return nil, status.Error(codes.InvalidArgument, "设备类型名称和编码不能为空")
	}
	id, err := s.repo.CreateDeviceType(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建设备类型失败: %v", err)
	}
	return &pb.DeviceIDResponse{Id: id}, nil
}

func (s *DeviceService) UpdateDeviceType(ctx context.Context, req *pb.DeviceTypeSaveRequest) (*pb.DeviceEmpty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "设备类型 ID 不正确")
	}
	if err := s.repo.UpdateDeviceType(ctx, req); err != nil {
		return nil, mapRepoError(err, "更新设备类型失败")
	}
	return &pb.DeviceEmpty{}, nil
}

func (s *DeviceService) DeleteDeviceType(ctx context.Context, req *pb.DeviceIDRequest) (*pb.DeviceEmpty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "设备类型 ID 不正确")
	}
	if err := s.repo.DeleteDeviceType(ctx, req.Id); err != nil {
		return nil, mapRepoError(err, "删除设备类型失败")
	}
	return &pb.DeviceEmpty{}, nil
}

func mapRepoError(err error, prefix string) error {
	if errors.Is(err, sql.ErrNoRows) {
		return status.Error(codes.NotFound, "数据不存在")
	}
	return status.Errorf(codes.Internal, "%s: %v", prefix, err)
}
