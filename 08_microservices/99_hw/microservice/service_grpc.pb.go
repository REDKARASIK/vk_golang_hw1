// This file is a stub. To be replaced by protoc generated file
package main

import (
	context "context"

	grpc "google.golang.org/grpc"
)

type AdminClient interface {
	Logging(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (Admin_LoggingClient, error)
	Statistics(ctx context.Context, in *StatInterval, opts ...grpc.CallOption) (Admin_StatisticsClient, error)
}

func NewAdminClient(cc grpc.ClientConnInterface) AdminClient {
	return nil
}

type Admin_LoggingClient interface {
	Recv() (*Event, error)
	grpc.ClientStream
}

type Admin_StatisticsClient interface {
	Recv() (*Stat, error)
	grpc.ClientStream
}

type AdminServer interface {
	Logging(*Nothing, Admin_LoggingServer) error
	Statistics(*StatInterval, Admin_StatisticsServer) error
}

func RegisterAdminServer(s grpc.ServiceRegistrar, srv AdminServer) {}

type Admin_LoggingServer interface {
	Send(*Event) error
	grpc.ServerStream
}

type Admin_StatisticsServer interface {
	Send(*Stat) error
	grpc.ServerStream
}

type BizClient interface {
	Check(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Nothing, error)
	Add(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Nothing, error)
	Test(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Nothing, error)
}

func NewBizClient(cc grpc.ClientConnInterface) BizClient {
	return nil
}

type BizServer interface {
	Check(context.Context, *Nothing) (*Nothing, error)
	Add(context.Context, *Nothing) (*Nothing, error)
	Test(context.Context, *Nothing) (*Nothing, error)
}

func RegisterBizServer(s grpc.ServiceRegistrar, srv BizServer) {}
