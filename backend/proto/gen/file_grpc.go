package pb

import (
	"context"

	"google.golang.org/grpc"
)

const FileServiceName = "file-service"

type FileServiceClient interface {
	UploadFileMeta(ctx context.Context, in *UploadFileMetaRequest, opts ...grpc.CallOption) (*FileIDResponse, error)
	ListFiles(ctx context.Context, in *FileListRequest, opts ...grpc.CallOption) (*FileListResponse, error)
	GetFile(ctx context.Context, in *FileIDRequest, opts ...grpc.CallOption) (*FileInfo, error)
	DeleteFile(ctx context.Context, in *FileIDRequest, opts ...grpc.CallOption) (*FileEmpty, error)
}

type fileServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFileServiceClient(cc grpc.ClientConnInterface) FileServiceClient {
	return &fileServiceClient{cc}
}

func (c *fileServiceClient) UploadFileMeta(ctx context.Context, in *UploadFileMetaRequest, opts ...grpc.CallOption) (*FileIDResponse, error) {
	out := new(FileIDResponse)
	err := c.cc.Invoke(ctx, "/microadmin.file.FileService/UploadFileMeta", in, out, opts...)
	return out, err
}

func (c *fileServiceClient) ListFiles(ctx context.Context, in *FileListRequest, opts ...grpc.CallOption) (*FileListResponse, error) {
	out := new(FileListResponse)
	err := c.cc.Invoke(ctx, "/microadmin.file.FileService/ListFiles", in, out, opts...)
	return out, err
}

func (c *fileServiceClient) GetFile(ctx context.Context, in *FileIDRequest, opts ...grpc.CallOption) (*FileInfo, error) {
	out := new(FileInfo)
	err := c.cc.Invoke(ctx, "/microadmin.file.FileService/GetFile", in, out, opts...)
	return out, err
}

func (c *fileServiceClient) DeleteFile(ctx context.Context, in *FileIDRequest, opts ...grpc.CallOption) (*FileEmpty, error) {
	out := new(FileEmpty)
	err := c.cc.Invoke(ctx, "/microadmin.file.FileService/DeleteFile", in, out, opts...)
	return out, err
}

type FileServiceServer interface {
	UploadFileMeta(context.Context, *UploadFileMetaRequest) (*FileIDResponse, error)
	ListFiles(context.Context, *FileListRequest) (*FileListResponse, error)
	GetFile(context.Context, *FileIDRequest) (*FileInfo, error)
	DeleteFile(context.Context, *FileIDRequest) (*FileEmpty, error)
}

func RegisterFileServiceServer(s grpc.ServiceRegistrar, srv FileServiceServer) {
	s.RegisterService(&FileService_ServiceDesc, srv)
}

func fileUnaryHandler[T any](srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor, fullMethod string, call func(context.Context, FileServiceServer, *T) (any, error)) (any, error) {
	in := new(T)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return call(ctx, srv.(FileServiceServer), in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: fullMethod}
	handler := func(ctx context.Context, req any) (any, error) {
		return call(ctx, srv.(FileServiceServer), req.(*T))
	}
	return interceptor(ctx, in, info, handler)
}

var FileService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "microadmin.file.FileService",
	HandlerType: (*FileServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "UploadFileMeta", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return fileUnaryHandler[UploadFileMetaRequest](s, c, d, i, "/microadmin.file.FileService/UploadFileMeta", func(ctx context.Context, srv FileServiceServer, req *UploadFileMetaRequest) (any, error) { return srv.UploadFileMeta(ctx, req) })
		}},
		{MethodName: "ListFiles", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return fileUnaryHandler[FileListRequest](s, c, d, i, "/microadmin.file.FileService/ListFiles", func(ctx context.Context, srv FileServiceServer, req *FileListRequest) (any, error) { return srv.ListFiles(ctx, req) })
		}},
		{MethodName: "GetFile", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return fileUnaryHandler[FileIDRequest](s, c, d, i, "/microadmin.file.FileService/GetFile", func(ctx context.Context, srv FileServiceServer, req *FileIDRequest) (any, error) { return srv.GetFile(ctx, req) })
		}},
		{MethodName: "DeleteFile", Handler: func(s any, c context.Context, d func(any) error, i grpc.UnaryServerInterceptor) (any, error) {
			return fileUnaryHandler[FileIDRequest](s, c, d, i, "/microadmin.file.FileService/DeleteFile", func(ctx context.Context, srv FileServiceServer, req *FileIDRequest) (any, error) { return srv.DeleteFile(ctx, req) })
		}},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/file.proto",
}
