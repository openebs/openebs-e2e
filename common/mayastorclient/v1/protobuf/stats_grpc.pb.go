// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.8
// source: stats.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// StatsRpcClient is the client API for StatsRpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StatsRpcClient interface {
	// RPC to get Pool IO stats.
	GetPoolIoStats(ctx context.Context, in *ListStatsOption, opts ...grpc.CallOption) (*PoolIoStatsResponse, error)
	// RPC to reset Io Stats for all Pool, Replica and Nexus hosted on the node.
	ResetIoStats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type statsRpcClient struct {
	cc grpc.ClientConnInterface
}

func NewStatsRpcClient(cc grpc.ClientConnInterface) StatsRpcClient {
	return &statsRpcClient{cc}
}

func (c *statsRpcClient) GetPoolIoStats(ctx context.Context, in *ListStatsOption, opts ...grpc.CallOption) (*PoolIoStatsResponse, error) {
	out := new(PoolIoStatsResponse)
	err := c.cc.Invoke(ctx, "/mayastor.v1.StatsRpc/GetPoolIoStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statsRpcClient) ResetIoStats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/mayastor.v1.StatsRpc/ResetIoStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StatsRpcServer is the server API for StatsRpc service.
// All implementations must embed UnimplementedStatsRpcServer
// for forward compatibility
type StatsRpcServer interface {
	// RPC to get Pool IO stats.
	GetPoolIoStats(context.Context, *ListStatsOption) (*PoolIoStatsResponse, error)
	// RPC to reset Io Stats for all Pool, Replica and Nexus hosted on the node.
	ResetIoStats(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	mustEmbedUnimplementedStatsRpcServer()
}

// UnimplementedStatsRpcServer must be embedded to have forward compatible implementations.
type UnimplementedStatsRpcServer struct {
}

func (UnimplementedStatsRpcServer) GetPoolIoStats(context.Context, *ListStatsOption) (*PoolIoStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPoolIoStats not implemented")
}
func (UnimplementedStatsRpcServer) ResetIoStats(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetIoStats not implemented")
}
func (UnimplementedStatsRpcServer) mustEmbedUnimplementedStatsRpcServer() {}

// UnsafeStatsRpcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StatsRpcServer will
// result in compilation errors.
type UnsafeStatsRpcServer interface {
	mustEmbedUnimplementedStatsRpcServer()
}

func RegisterStatsRpcServer(s grpc.ServiceRegistrar, srv StatsRpcServer) {
	s.RegisterService(&StatsRpc_ServiceDesc, srv)
}

func _StatsRpc_GetPoolIoStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListStatsOption)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatsRpcServer).GetPoolIoStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mayastor.v1.StatsRpc/GetPoolIoStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatsRpcServer).GetPoolIoStats(ctx, req.(*ListStatsOption))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatsRpc_ResetIoStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatsRpcServer).ResetIoStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mayastor.v1.StatsRpc/ResetIoStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatsRpcServer).ResetIoStats(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// StatsRpc_ServiceDesc is the grpc.ServiceDesc for StatsRpc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StatsRpc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mayastor.v1.StatsRpc",
	HandlerType: (*StatsRpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPoolIoStats",
			Handler:    _StatsRpc_GetPoolIoStats_Handler,
		},
		{
			MethodName: "ResetIoStats",
			Handler:    _StatsRpc_ResetIoStats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "stats.proto",
}
