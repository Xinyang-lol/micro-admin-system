package client

import (
	"context"

	"micro-admin-system/backend/common/grpcjson"
	"micro-admin-system/backend/common/registry"
	pb "micro-admin-system/backend/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	discovery *registry.Client
}

func New(discovery *registry.Client) *Clients {
	return &Clients{discovery: discovery}
}

func (c *Clients) dial(ctx context.Context, service string) (*grpc.ClientConn, error) {
	address, err := c.discovery.NextAddress(service)
	if err != nil {
		return nil, err
	}
	return grpc.DialContext(
		ctx,
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(grpcjson.Codec{})),
		grpc.WithBlock(),
	)
}

func (c *Clients) User(ctx context.Context) (pb.UserServiceClient, *grpc.ClientConn, error) {
	conn, err := c.dial(ctx, pb.UserServiceName)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewUserServiceClient(conn), conn, nil
}

func (c *Clients) Device(ctx context.Context) (pb.DeviceServiceClient, *grpc.ClientConn, error) {
	conn, err := c.dial(ctx, pb.DeviceServiceName)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewDeviceServiceClient(conn), conn, nil
}

func (c *Clients) File(ctx context.Context) (pb.FileServiceClient, *grpc.ClientConn, error) {
	conn, err := c.dial(ctx, pb.FileServiceName)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewFileServiceClient(conn), conn, nil
}
