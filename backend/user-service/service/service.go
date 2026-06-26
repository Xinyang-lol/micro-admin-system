package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"micro-admin-system/backend/common/auth"
	pb "micro-admin-system/backend/proto/gen"
	"micro-admin-system/backend/user-service/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	repo      *repository.Repository
	jwtSecret string
	tokenTTL  time.Duration
}

func New(repo *repository.Repository, jwtSecret string, tokenTTL time.Duration) *UserService {
	return &UserService{repo: repo, jwtSecret: jwtSecret, tokenTTL: tokenTTL}
}

func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		return nil, status.Error(codes.InvalidArgument, "用户名和密码不能为空")
	}
	user, err := s.repo.GetUserByUsername(ctx, req.Username)
	if errors.Is(err, sql.ErrNoRows) || (err == nil && !auth.CheckPassword(user.PasswordHash, req.Password)) {
		return nil, status.Error(codes.Unauthenticated, "用户名或密码错误")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询用户失败: %v", err)
	}
	if user.Status != 1 {
		return nil, status.Error(codes.PermissionDenied, "用户已被禁用")
	}
	info, err := s.repo.GetUserInfo(ctx, user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询用户权限失败: %v", err)
	}
	token, err := auth.GenerateToken(s.jwtSecret, info.Id, info.Username, info.Roles, info.Permissions, s.tokenTTL)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "生成 Token 失败: %v", err)
	}
	return &pb.LoginResponse{Token: token, User: info}, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, req *pb.UserIDRequest) (*pb.UserInfo, error) {
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "用户 ID 不正确")
	}
	info, err := s.repo.GetUserInfo(ctx, req.UserId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询用户失败: %v", err)
	}
	return info, nil
}

func (s *UserService) ListUsers(ctx context.Context, req *pb.ListRequest) (*pb.UserListResponse, error) {
	resp, err := s.repo.ListUsers(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询用户列表失败: %v", err)
	}
	return resp, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.UserSaveRequest) (*pb.IDResponse, error) {
	if strings.TrimSpace(req.Username) == "" {
		return nil, status.Error(codes.InvalidArgument, "用户名不能为空")
	}
	password := req.Password
	if password == "" {
		password = "123456"
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "密码加密失败: %v", err)
	}
	if req.Status == 0 {
		req.Status = 1
	}
	id, err := s.repo.CreateUser(ctx, req, hash)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建用户失败: %v", err)
	}
	return &pb.IDResponse{Id: id}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UserSaveRequest) (*pb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "用户 ID 不正确")
	}
	if err := s.repo.UpdateUser(ctx, req); err != nil {
		return nil, mapRepoError(err, "更新用户失败")
	}
	return &pb.Empty{}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.IDRequest) (*pb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "用户 ID 不正确")
	}
	if err := s.repo.DeleteUser(ctx, req.Id); err != nil {
		return nil, mapRepoError(err, "删除用户失败")
	}
	return &pb.Empty{}, nil
}

func (s *UserService) UpdateUserStatus(ctx context.Context, req *pb.UserStatusRequest) (*pb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "用户 ID 不正确")
	}
	if req.Status != 0 && req.Status != 1 {
		return nil, status.Error(codes.InvalidArgument, "用户状态只能是 0 或 1")
	}
	if err := s.repo.UpdateUserStatus(ctx, req.Id, req.Status); err != nil {
		return nil, mapRepoError(err, "更新用户状态失败")
	}
	return &pb.Empty{}, nil
}

func (s *UserService) ResetPassword(ctx context.Context, req *pb.UserPasswordRequest) (*pb.Empty, error) {
	if req.Id <= 0 || len(req.Password) < 6 {
		return nil, status.Error(codes.InvalidArgument, "用户 ID 不正确或密码长度不足 6 位")
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "密码加密失败: %v", err)
	}
	if err := s.repo.ResetPassword(ctx, req.Id, hash); err != nil {
		return nil, mapRepoError(err, "重置密码失败")
	}
	return &pb.Empty{}, nil
}

func (s *UserService) AssignUserRoles(ctx context.Context, req *pb.UserRolesRequest) (*pb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "用户 ID 不正确")
	}
	if err := s.repo.AssignUserRoles(ctx, req.Id, req.RoleIds); err != nil {
		return nil, status.Errorf(codes.Internal, "分配用户角色失败: %v", err)
	}
	return &pb.Empty{}, nil
}

func (s *UserService) ListRoles(ctx context.Context, req *pb.ListRequest) (*pb.RoleListResponse, error) {
	resp, err := s.repo.ListRoles(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询角色列表失败: %v", err)
	}
	return resp, nil
}

func (s *UserService) CreateRole(ctx context.Context, req *pb.RoleSaveRequest) (*pb.IDResponse, error) {
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Code) == "" {
		return nil, status.Error(codes.InvalidArgument, "角色名称和编码不能为空")
	}
	if req.Status == 0 {
		req.Status = 1
	}
	id, err := s.repo.CreateRole(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建角色失败: %v", err)
	}
	return &pb.IDResponse{Id: id}, nil
}

func (s *UserService) UpdateRole(ctx context.Context, req *pb.RoleSaveRequest) (*pb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "角色 ID 不正确")
	}
	if err := s.repo.UpdateRole(ctx, req); err != nil {
		return nil, mapRepoError(err, "更新角色失败")
	}
	return &pb.Empty{}, nil
}

func (s *UserService) DeleteRole(ctx context.Context, req *pb.IDRequest) (*pb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "角色 ID 不正确")
	}
	if err := s.repo.DeleteRole(ctx, req.Id); err != nil {
		return nil, mapRepoError(err, "删除角色失败")
	}
	return &pb.Empty{}, nil
}

func (s *UserService) AssignRoleMenus(ctx context.Context, req *pb.RoleMenusRequest) (*pb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "角色 ID 不正确")
	}
	if err := s.repo.AssignRoleMenus(ctx, req.Id, req.MenuIds); err != nil {
		return nil, status.Errorf(codes.Internal, "分配角色菜单失败: %v", err)
	}
	return &pb.Empty{}, nil
}

func (s *UserService) ListMenus(ctx context.Context, _ *pb.Empty) (*pb.MenuTreeResponse, error) {
	resp, err := s.repo.ListMenus(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询菜单树失败: %v", err)
	}
	return resp, nil
}

func (s *UserService) CreateMenu(ctx context.Context, req *pb.MenuSaveRequest) (*pb.IDResponse, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, status.Error(codes.InvalidArgument, "菜单名称不能为空")
	}
	id, err := s.repo.CreateMenu(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建菜单失败: %v", err)
	}
	return &pb.IDResponse{Id: id}, nil
}

func (s *UserService) UpdateMenu(ctx context.Context, req *pb.MenuSaveRequest) (*pb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "菜单 ID 不正确")
	}
	if err := s.repo.UpdateMenu(ctx, req); err != nil {
		return nil, mapRepoError(err, "更新菜单失败")
	}
	return &pb.Empty{}, nil
}

func (s *UserService) DeleteMenu(ctx context.Context, req *pb.IDRequest) (*pb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "菜单 ID 不正确")
	}
	if err := s.repo.DeleteMenu(ctx, req.Id); err != nil {
		return nil, mapRepoError(err, "删除菜单失败")
	}
	return &pb.Empty{}, nil
}

func (s *UserService) ListDepts(ctx context.Context, _ *pb.Empty) (*pb.DeptTreeResponse, error) {
	resp, err := s.repo.ListDepts(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询部门树失败: %v", err)
	}
	return resp, nil
}

func (s *UserService) CreateDept(ctx context.Context, req *pb.DeptSaveRequest) (*pb.IDResponse, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, status.Error(codes.InvalidArgument, "部门名称不能为空")
	}
	if req.Status == 0 {
		req.Status = 1
	}
	id, err := s.repo.CreateDept(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建部门失败: %v", err)
	}
	return &pb.IDResponse{Id: id}, nil
}

func (s *UserService) UpdateDept(ctx context.Context, req *pb.DeptSaveRequest) (*pb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "部门 ID 不正确")
	}
	if err := s.repo.UpdateDept(ctx, req); err != nil {
		return nil, mapRepoError(err, "更新部门失败")
	}
	return &pb.Empty{}, nil
}

func (s *UserService) DeleteDept(ctx context.Context, req *pb.IDRequest) (*pb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "部门 ID 不正确")
	}
	if err := s.repo.DeleteDept(ctx, req.Id); err != nil {
		return nil, mapRepoError(err, "删除部门失败")
	}
	return &pb.Empty{}, nil
}

func (s *UserService) CheckPermission(ctx context.Context, req *pb.PermissionRequest) (*pb.PermissionResponse, error) {
	if req.UserId <= 0 || strings.TrimSpace(req.Permission) == "" {
		return &pb.PermissionResponse{Allowed: false}, nil
	}
	allowed, err := s.repo.HasPermission(ctx, req.UserId, req.Permission)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "校验权限失败: %v", err)
	}
	return &pb.PermissionResponse{Allowed: allowed}, nil
}

func mapRepoError(err error, prefix string) error {
	if errors.Is(err, sql.ErrNoRows) {
		return status.Error(codes.NotFound, "数据不存在")
	}
	if strings.Contains(err.Error(), "cannot") || strings.Contains(err.Error(), "child nodes") {
		return status.Error(codes.FailedPrecondition, err.Error())
	}
	return status.Errorf(codes.Internal, "%s: %v", prefix, err)
}
