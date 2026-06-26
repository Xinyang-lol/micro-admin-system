package pb

import (
	"context"

	"google.golang.org/grpc"
)

const DeviceServiceName = "device-service"

type DeviceServiceClient interface {
	ListDevices(ctx context.Context, in *DeviceListRequest, opts ...grpc.CallOption) (*DeviceListResponse, error)
	CreateDevice(ctx context.Context, in *DeviceSaveRequest, opts ...grpc.CallOption) (*DeviceIDResponse, error)
	UpdateDevice(ctx context.Context, in *DeviceSaveRequest, opts ...grpc.CallOption) (*DeviceEmpty, error)
	DeleteDevice(ctx context.Context, in *DeviceIDRequest, opts ...grpc.CallOption) (*DeviceEmpty, error)
	GetDeviceStatistics(ctx context.Context, in *DeviceEmpty, opts ...grpc.CallOption) (*DeviceStatisticsResponse, error)
	ListDeviceTypes(ctx context.Context, in *DeviceEmpty, opts ...grpc.CallOption) (*DeviceTypeListResponse, error)
	CreateDeviceType(ctx context.Context, in *DeviceTypeSaveRequest, opts ...grpc.CallOption) (*DeviceIDResponse, error)
	UpdateDeviceType(ctx context.Context, in *DeviceTypeSaveRequest, opts ...grpc.CallOption) (*DeviceEmpty, error)
	DeleteDeviceType(ctx context.Context, in *DeviceIDRequest, opts ...grpc.CallOption) (*DeviceEmpty, error)
}

type deviceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDeviceServiceClient(cc grpc.ClientConnInterface) DeviceServiceClient {
	return &deviceServiceClient{cc}
}

func (c *deviceServiceClient) ListDevices(ctx context.Context, in *DeviceListRequest, opts ...grpc.CallOption) (*DeviceListResponse, error) {
	out := new(DeviceListResponse)
	err := c.cc.Invoke(ctx, "/microadmin.device.DeviceService/ListDevices", in, out, opts...)
	return out, err
}

func (c *deviceServiceClient) CreateDevice(ctx context.Context, in *DeviceSaveRequest, opts ...grpc.CallOption) (*DeviceIDResponse, error) {
	out := new(DeviceIDResponse)
	err := c.cc.Invoke(ctx, "/microadmin.device.DeviceService/CreateDevice", in, out, opts...)
	return out, err
}

func (c *deviceServiceClient) UpdateDevice(ctx context.Context, in *DeviceSaveRequest, opts ...grpc.CallOption) (*DeviceEmpty, error) {
	out := new(DeviceEmpty)
	err := c.cc.Invoke(ctx, "/microadmin.device.DeviceService/UpdateDevice", in, out, opts...)
	return out, err
}

func (c *deviceServiceClient) DeleteDevice(ctx context.Context, in *DeviceIDRequest, opts ...grpc.CallOption) (*DeviceEmpty, error) {
	out := new(DeviceEmpty)
	err := c.cc.Invoke(ctx, "/microadmin.device.DeviceService/DeleteDevice", in, out, opts...)
	return out, err
}

func (c *deviceServiceClient) GetDeviceStatistics(ctx context.Context, in *DeviceEmpty, opts ...grpc.CallOption) (*DeviceStatisticsResponse, error) {
	out := new(DeviceStatisticsResponse)
	err := c.cc.Invoke(ctx, "/microadmin.device.DeviceService/GetDeviceStatistics", in, out, opts...)
	return out, err
}

func (c *deviceServiceClient) ListDeviceTypes(ctx context.Context, in *DeviceEmpty, opts ...grpc.CallOption) (*DeviceTypeListResponse, error) {
	out := new(DeviceTypeListResponse)
	err := c.cc.Invoke(ctx, "/microadmin.device.DeviceService/ListDeviceTypes", in, out, opts...)
	return out, err
}

func (c *deviceServiceClient) CreateDeviceType(ctx context.Context, in *DeviceTypeSaveRequest, opts ...grpc.CallOption) (*DeviceIDResponse, error) {
	out := new(DeviceIDResponse)
	err := c.cc.Invoke(ctx, "/microadmin.device.DeviceService/CreateDeviceType", in, out, opts...)
	return out, err
}

func (c *deviceServiceClient) UpdateDeviceType(ctx context.Context, in *DeviceTypeSaveRequest, opts ...grpc.CallOption) (*DeviceEmpty, error) {
	out := new(DeviceEmpty)
	err := c.cc.Invoke(ctx, "/microadmin.device.DeviceService/UpdateDeviceType", in, out, opts...)
	return out, err
}

func (c *deviceServiceClient) DeleteDeviceType(ctx context.Context, in *DeviceIDRequest, opts ...grpc.CallOption) (*DeviceEmpty, error) {
	out := new(DeviceEmpty)
	err := c.cc.Invoke(ctx, "/microadmin.device.DeviceService/DeleteDeviceType", in, out, opts...)
	return out, err
}

type DeviceServiceServer interface {
	ListDevices(context.Context, *DeviceListRequest) (*DeviceListResponse, error)
	CreateDevice(context.Context, *DeviceSaveRequest) (*DeviceIDResponse, error)
	UpdateDevice(context.Context, *DeviceSaveRequest) (*DeviceEmpty, error)
	DeleteDevice(context.Context, *DeviceIDRequest) (*DeviceEmpty, error)
	GetDeviceStatistics(context.Context, *DeviceEmpty) (*DeviceStatisticsResponse, error)
	ListDeviceTypes(context.Context, *DeviceEmpty) (*DeviceTypeListResponse, error)
	CreateDeviceType(context.Context, *DeviceTypeSaveRequest) (*DeviceIDResponse, error)
	UpdateDeviceType(context.Context, *DeviceTypeSaveRequest) (*DeviceEmpty, error)
	DeleteDeviceType(context.Context, *DeviceIDRequest) (*DeviceEmpty, error)
}

func RegisterDeviceServiceServer(s grpc.ServiceRegistrar, srv DeviceServiceServer) {
	s.RegisterService(&DeviceService_ServiceDesc, srv)
}

func deviceUnaryHandler[T any](srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor, fullMethod string, call func(context.Context, DeviceServiceServer, *T) (any, error)) (any, error) {
	in := new(T)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return call(ctx, srv.(DeviceServiceServer), in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: fullMethod}
	handler := func(ctx context.Context, req any) (any, error) {
		return call(ctx, srv.(DeviceServiceServer), req.(*T))
	}
	return interceptor(ctx, in, info, handler)
}

var DeviceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "microadmin.device.DeviceService",
	HandlerType: (*DeviceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "ListDevices", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return deviceUnaryHandler[DeviceListRequest](s, c, d, i, "/microadmin.device.DeviceService/ListDevices", func(ctx context.Context, srv DeviceServiceServer, req *DeviceListRequest) (any, error) { return srv.ListDevices(ctx, req) })
		}},
		{MethodName: "CreateDevice", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return deviceUnaryHandler[DeviceSaveRequest](s, c, d, i, "/microadmin.device.DeviceService/CreateDevice", func(ctx context.Context, srv DeviceServiceServer, req *DeviceSaveRequest) (any, error) { return srv.CreateDevice(ctx, req) })
		}},
		{MethodName: "UpdateDevice", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return deviceUnaryHandler[DeviceSaveRequest](s, c, d, i, "/microadmin.device.DeviceService/UpdateDevice", func(ctx context.Context, srv DeviceServiceServer, req *DeviceSaveRequest) (any, error) { return srv.UpdateDevice(ctx, req) })
		}},
		{MethodName: "DeleteDevice", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return deviceUnaryHandler[DeviceIDRequest](s, c, d, i, "/microadmin.device.DeviceService/DeleteDevice", func(ctx context.Context, srv DeviceServiceServer, req *DeviceIDRequest) (any, error) { return srv.DeleteDevice(ctx, req) })
		}},
		{MethodName: "GetDeviceStatistics", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return deviceUnaryHandler[DeviceEmpty](s, c, d, i, "/microadmin.device.DeviceService/GetDeviceStatistics", func(ctx context.Context, srv DeviceServiceServer, req *DeviceEmpty) (any, error) { return srv.GetDeviceStatistics(ctx, req) })
		}},
		{MethodName: "ListDeviceTypes", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return deviceUnaryHandler[DeviceEmpty](s, c, d, i, "/microadmin.device.DeviceService/ListDeviceTypes", func(ctx context.Context, srv DeviceServiceServer, req *DeviceEmpty) (any, error) { return srv.ListDeviceTypes(ctx, req) })
		}},
		{MethodName: "CreateDeviceType", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return deviceUnaryHandler[DeviceTypeSaveRequest](s, c, d, i, "/microadmin.device.DeviceService/CreateDeviceType", func(ctx context.Context, srv DeviceServiceServer, req *DeviceTypeSaveRequest) (any, error) { return srv.CreateDeviceType(ctx, req) })
		}},
		{MethodName: "UpdateDeviceType", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return deviceUnaryHandler[DeviceTypeSaveRequest](s, c, d, i, "/microadmin.device.DeviceService/UpdateDeviceType", func(ctx context.Context, srv DeviceServiceServer, req *DeviceTypeSaveRequest) (any, error) { return srv.UpdateDeviceType(ctx, req) })
		}},
		{MethodName: "DeleteDeviceType", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return deviceUnaryHandler[DeviceIDRequest](s, c, d, i, "/microadmin.device.DeviceService/DeleteDeviceType", func(ctx context.Context, srv DeviceServiceServer, req *DeviceIDRequest) (any, error) { return srv.DeleteDeviceType(ctx, req) })
		}},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/device.proto",
}
