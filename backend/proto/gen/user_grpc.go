package pb

import (
	"context"

	"google.golang.org/grpc"
)

const UserServiceName = "user-service"

type UserServiceClient interface {
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	GetUserInfo(ctx context.Context, in *UserIDRequest, opts ...grpc.CallOption) (*UserInfo, error)
	ListUsers(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*UserListResponse, error)
	CreateUser(ctx context.Context, in *UserSaveRequest, opts ...grpc.CallOption) (*IDResponse, error)
	UpdateUser(ctx context.Context, in *UserSaveRequest, opts ...grpc.CallOption) (*Empty, error)
	DeleteUser(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*Empty, error)
	UpdateUserStatus(ctx context.Context, in *UserStatusRequest, opts ...grpc.CallOption) (*Empty, error)
	ResetPassword(ctx context.Context, in *UserPasswordRequest, opts ...grpc.CallOption) (*Empty, error)
	AssignUserRoles(ctx context.Context, in *UserRolesRequest, opts ...grpc.CallOption) (*Empty, error)
	ListRoles(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*RoleListResponse, error)
	CreateRole(ctx context.Context, in *RoleSaveRequest, opts ...grpc.CallOption) (*IDResponse, error)
	UpdateRole(ctx context.Context, in *RoleSaveRequest, opts ...grpc.CallOption) (*Empty, error)
	DeleteRole(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*Empty, error)
	AssignRoleMenus(ctx context.Context, in *RoleMenusRequest, opts ...grpc.CallOption) (*Empty, error)
	ListMenus(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*MenuTreeResponse, error)
	CreateMenu(ctx context.Context, in *MenuSaveRequest, opts ...grpc.CallOption) (*IDResponse, error)
	UpdateMenu(ctx context.Context, in *MenuSaveRequest, opts ...grpc.CallOption) (*Empty, error)
	DeleteMenu(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*Empty, error)
	ListDepts(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*DeptTreeResponse, error)
	CreateDept(ctx context.Context, in *DeptSaveRequest, opts ...grpc.CallOption) (*IDResponse, error)
	UpdateDept(ctx context.Context, in *DeptSaveRequest, opts ...grpc.CallOption) (*Empty, error)
	DeleteDept(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*Empty, error)
	CheckPermission(ctx context.Context, in *PermissionRequest, opts ...grpc.CallOption) (*PermissionResponse, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/Login", in, out, opts...)
	return out, err
}

func (c *userServiceClient) GetUserInfo(ctx context.Context, in *UserIDRequest, opts ...grpc.CallOption) (*UserInfo, error) {
	out := new(UserInfo)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/GetUserInfo", in, out, opts...)
	return out, err
}

func (c *userServiceClient) ListUsers(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*UserListResponse, error) {
	out := new(UserListResponse)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/ListUsers", in, out, opts...)
	return out, err
}

func (c *userServiceClient) CreateUser(ctx context.Context, in *UserSaveRequest, opts ...grpc.CallOption) (*IDResponse, error) {
	out := new(IDResponse)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/CreateUser", in, out, opts...)
	return out, err
}

func (c *userServiceClient) UpdateUser(ctx context.Context, in *UserSaveRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/UpdateUser", in, out, opts...)
	return out, err
}

func (c *userServiceClient) DeleteUser(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/DeleteUser", in, out, opts...)
	return out, err
}

func (c *userServiceClient) UpdateUserStatus(ctx context.Context, in *UserStatusRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/UpdateUserStatus", in, out, opts...)
	return out, err
}

func (c *userServiceClient) ResetPassword(ctx context.Context, in *UserPasswordRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/ResetPassword", in, out, opts...)
	return out, err
}

func (c *userServiceClient) AssignUserRoles(ctx context.Context, in *UserRolesRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/AssignUserRoles", in, out, opts...)
	return out, err
}

func (c *userServiceClient) ListRoles(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*RoleListResponse, error) {
	out := new(RoleListResponse)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/ListRoles", in, out, opts...)
	return out, err
}

func (c *userServiceClient) CreateRole(ctx context.Context, in *RoleSaveRequest, opts ...grpc.CallOption) (*IDResponse, error) {
	out := new(IDResponse)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/CreateRole", in, out, opts...)
	return out, err
}

func (c *userServiceClient) UpdateRole(ctx context.Context, in *RoleSaveRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/UpdateRole", in, out, opts...)
	return out, err
}

func (c *userServiceClient) DeleteRole(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/DeleteRole", in, out, opts...)
	return out, err
}

func (c *userServiceClient) AssignRoleMenus(ctx context.Context, in *RoleMenusRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/AssignRoleMenus", in, out, opts...)
	return out, err
}

func (c *userServiceClient) ListMenus(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*MenuTreeResponse, error) {
	out := new(MenuTreeResponse)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/ListMenus", in, out, opts...)
	return out, err
}

func (c *userServiceClient) CreateMenu(ctx context.Context, in *MenuSaveRequest, opts ...grpc.CallOption) (*IDResponse, error) {
	out := new(IDResponse)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/CreateMenu", in, out, opts...)
	return out, err
}

func (c *userServiceClient) UpdateMenu(ctx context.Context, in *MenuSaveRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/UpdateMenu", in, out, opts...)
	return out, err
}

func (c *userServiceClient) DeleteMenu(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/DeleteMenu", in, out, opts...)
	return out, err
}

func (c *userServiceClient) ListDepts(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*DeptTreeResponse, error) {
	out := new(DeptTreeResponse)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/ListDepts", in, out, opts...)
	return out, err
}

func (c *userServiceClient) CreateDept(ctx context.Context, in *DeptSaveRequest, opts ...grpc.CallOption) (*IDResponse, error) {
	out := new(IDResponse)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/CreateDept", in, out, opts...)
	return out, err
}

func (c *userServiceClient) UpdateDept(ctx context.Context, in *DeptSaveRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/UpdateDept", in, out, opts...)
	return out, err
}

func (c *userServiceClient) DeleteDept(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/DeleteDept", in, out, opts...)
	return out, err
}

func (c *userServiceClient) CheckPermission(ctx context.Context, in *PermissionRequest, opts ...grpc.CallOption) (*PermissionResponse, error) {
	out := new(PermissionResponse)
	err := c.cc.Invoke(ctx, "/microadmin.user.UserService/CheckPermission", in, out, opts...)
	return out, err
}

type UserServiceServer interface {
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	GetUserInfo(context.Context, *UserIDRequest) (*UserInfo, error)
	ListUsers(context.Context, *ListRequest) (*UserListResponse, error)
	CreateUser(context.Context, *UserSaveRequest) (*IDResponse, error)
	UpdateUser(context.Context, *UserSaveRequest) (*Empty, error)
	DeleteUser(context.Context, *IDRequest) (*Empty, error)
	UpdateUserStatus(context.Context, *UserStatusRequest) (*Empty, error)
	ResetPassword(context.Context, *UserPasswordRequest) (*Empty, error)
	AssignUserRoles(context.Context, *UserRolesRequest) (*Empty, error)
	ListRoles(context.Context, *ListRequest) (*RoleListResponse, error)
	CreateRole(context.Context, *RoleSaveRequest) (*IDResponse, error)
	UpdateRole(context.Context, *RoleSaveRequest) (*Empty, error)
	DeleteRole(context.Context, *IDRequest) (*Empty, error)
	AssignRoleMenus(context.Context, *RoleMenusRequest) (*Empty, error)
	ListMenus(context.Context, *Empty) (*MenuTreeResponse, error)
	CreateMenu(context.Context, *MenuSaveRequest) (*IDResponse, error)
	UpdateMenu(context.Context, *MenuSaveRequest) (*Empty, error)
	DeleteMenu(context.Context, *IDRequest) (*Empty, error)
	ListDepts(context.Context, *Empty) (*DeptTreeResponse, error)
	CreateDept(context.Context, *DeptSaveRequest) (*IDResponse, error)
	UpdateDept(context.Context, *DeptSaveRequest) (*Empty, error)
	DeleteDept(context.Context, *IDRequest) (*Empty, error)
	CheckPermission(context.Context, *PermissionRequest) (*PermissionResponse, error)
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func userUnaryHandler[T any](srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor, fullMethod string, call func(context.Context, UserServiceServer, *T) (any, error)) (any, error) {
	in := new(T)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return call(ctx, srv.(UserServiceServer), in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: fullMethod}
	handler := func(ctx context.Context, req any) (any, error) {
		return call(ctx, srv.(UserServiceServer), req.(*T))
	}
	return interceptor(ctx, in, info, handler)
}

var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "microadmin.user.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Login", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[LoginRequest](s, c, d, i, "/microadmin.user.UserService/Login", func(ctx context.Context, srv UserServiceServer, req *LoginRequest) (any, error) { return srv.Login(ctx, req) })
		}},
		{MethodName: "GetUserInfo", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[UserIDRequest](s, c, d, i, "/microadmin.user.UserService/GetUserInfo", func(ctx context.Context, srv UserServiceServer, req *UserIDRequest) (any, error) { return srv.GetUserInfo(ctx, req) })
		}},
		{MethodName: "ListUsers", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[ListRequest](s, c, d, i, "/microadmin.user.UserService/ListUsers", func(ctx context.Context, srv UserServiceServer, req *ListRequest) (any, error) { return srv.ListUsers(ctx, req) })
		}},
		{MethodName: "CreateUser", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[UserSaveRequest](s, c, d, i, "/microadmin.user.UserService/CreateUser", func(ctx context.Context, srv UserServiceServer, req *UserSaveRequest) (any, error) { return srv.CreateUser(ctx, req) })
		}},
		{MethodName: "UpdateUser", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[UserSaveRequest](s, c, d, i, "/microadmin.user.UserService/UpdateUser", func(ctx context.Context, srv UserServiceServer, req *UserSaveRequest) (any, error) { return srv.UpdateUser(ctx, req) })
		}},
		{MethodName: "DeleteUser", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[IDRequest](s, c, d, i, "/microadmin.user.UserService/DeleteUser", func(ctx context.Context, srv UserServiceServer, req *IDRequest) (any, error) { return srv.DeleteUser(ctx, req) })
		}},
		{MethodName: "UpdateUserStatus", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[UserStatusRequest](s, c, d, i, "/microadmin.user.UserService/UpdateUserStatus", func(ctx context.Context, srv UserServiceServer, req *UserStatusRequest) (any, error) { return srv.UpdateUserStatus(ctx, req) })
		}},
		{MethodName: "ResetPassword", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[UserPasswordRequest](s, c, d, i, "/microadmin.user.UserService/ResetPassword", func(ctx context.Context, srv UserServiceServer, req *UserPasswordRequest) (any, error) { return srv.ResetPassword(ctx, req) })
		}},
		{MethodName: "AssignUserRoles", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[UserRolesRequest](s, c, d, i, "/microadmin.user.UserService/AssignUserRoles", func(ctx context.Context, srv UserServiceServer, req *UserRolesRequest) (any, error) { return srv.AssignUserRoles(ctx, req) })
		}},
		{MethodName: "ListRoles", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[ListRequest](s, c, d, i, "/microadmin.user.UserService/ListRoles", func(ctx context.Context, srv UserServiceServer, req *ListRequest) (any, error) { return srv.ListRoles(ctx, req) })
		}},
		{MethodName: "CreateRole", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[RoleSaveRequest](s, c, d, i, "/microadmin.user.UserService/CreateRole", func(ctx context.Context, srv UserServiceServer, req *RoleSaveRequest) (any, error) { return srv.CreateRole(ctx, req) })
		}},
		{MethodName: "UpdateRole", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[RoleSaveRequest](s, c, d, i, "/microadmin.user.UserService/UpdateRole", func(ctx context.Context, srv UserServiceServer, req *RoleSaveRequest) (any, error) { return srv.UpdateRole(ctx, req) })
		}},
		{MethodName: "DeleteRole", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[IDRequest](s, c, d, i, "/microadmin.user.UserService/DeleteRole", func(ctx context.Context, srv UserServiceServer, req *IDRequest) (any, error) { return srv.DeleteRole(ctx, req) })
		}},
		{MethodName: "AssignRoleMenus", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[RoleMenusRequest](s, c, d, i, "/microadmin.user.UserService/AssignRoleMenus", func(ctx context.Context, srv UserServiceServer, req *RoleMenusRequest) (any, error) { return srv.AssignRoleMenus(ctx, req) })
		}},
		{MethodName: "ListMenus", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[Empty](s, c, d, i, "/microadmin.user.UserService/ListMenus", func(ctx context.Context, srv UserServiceServer, req *Empty) (any, error) { return srv.ListMenus(ctx, req) })
		}},
		{MethodName: "CreateMenu", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[MenuSaveRequest](s, c, d, i, "/microadmin.user.UserService/CreateMenu", func(ctx context.Context, srv UserServiceServer, req *MenuSaveRequest) (any, error) { return srv.CreateMenu(ctx, req) })
		}},
		{MethodName: "UpdateMenu", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[MenuSaveRequest](s, c, d, i, "/microadmin.user.UserService/UpdateMenu", func(ctx context.Context, srv UserServiceServer, req *MenuSaveRequest) (any, error) { return srv.UpdateMenu(ctx, req) })
		}},
		{MethodName: "DeleteMenu", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[IDRequest](s, c, d, i, "/microadmin.user.UserService/DeleteMenu", func(ctx context.Context, srv UserServiceServer, req *IDRequest) (any, error) { return srv.DeleteMenu(ctx, req) })
		}},
		{MethodName: "ListDepts", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[Empty](s, c, d, i, "/microadmin.user.UserService/ListDepts", func(ctx context.Context, srv UserServiceServer, req *Empty) (any, error) { return srv.ListDepts(ctx, req) })
		}},
		{MethodName: "CreateDept", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[DeptSaveRequest](s, c, d, i, "/microadmin.user.UserService/CreateDept", func(ctx context.Context, srv UserServiceServer, req *DeptSaveRequest) (any, error) { return srv.CreateDept(ctx, req) })
		}},
		{MethodName: "UpdateDept", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[DeptSaveRequest](s, c, d, i, "/microadmin.user.UserService/UpdateDept", func(ctx context.Context, srv UserServiceServer, req *DeptSaveRequest) (any, error) { return srv.UpdateDept(ctx, req) })
		}},
		{MethodName: "DeleteDept", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[IDRequest](s, c, d, i, "/microadmin.user.UserService/DeleteDept", func(ctx context.Context, srv UserServiceServer, req *IDRequest) (any, error) { return srv.DeleteDept(ctx, req) })
		}},
		{MethodName: "CheckPermission", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return userUnaryHandler[PermissionRequest](s, c, d, i, "/microadmin.user.UserService/CheckPermission", func(ctx context.Context, srv UserServiceServer, req *PermissionRequest) (any, error) { return srv.CheckPermission(ctx, req) })
		}},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/user.proto",
}
