// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: protos/watcher.proto

package protos

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Watcher_GetInfo_FullMethodName         = "/Watcher/GetInfo"
	Watcher_MoreInfo_FullMethodName        = "/Watcher/MoreInfo"
	Watcher_SubscribeTicker_FullMethodName = "/Watcher/SubscribeTicker"
)

// WatcherClient is the client API for Watcher service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WatcherClient interface {
	GetInfo(ctx context.Context, in *TickerRequest, opts ...grpc.CallOption) (*TickerResponse, error)
	MoreInfo(ctx context.Context, in *TickerRequest, opts ...grpc.CallOption) (*CompanyResponse, error)
	SubscribeTicker(ctx context.Context, opts ...grpc.CallOption) (Watcher_SubscribeTickerClient, error)
}

type watcherClient struct {
	cc grpc.ClientConnInterface
}

func NewWatcherClient(cc grpc.ClientConnInterface) WatcherClient {
	return &watcherClient{cc}
}

func (c *watcherClient) GetInfo(ctx context.Context, in *TickerRequest, opts ...grpc.CallOption) (*TickerResponse, error) {
	out := new(TickerResponse)
	err := c.cc.Invoke(ctx, Watcher_GetInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *watcherClient) MoreInfo(ctx context.Context, in *TickerRequest, opts ...grpc.CallOption) (*CompanyResponse, error) {
	out := new(CompanyResponse)
	err := c.cc.Invoke(ctx, Watcher_MoreInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *watcherClient) SubscribeTicker(ctx context.Context, opts ...grpc.CallOption) (Watcher_SubscribeTickerClient, error) {
	stream, err := c.cc.NewStream(ctx, &Watcher_ServiceDesc.Streams[0], Watcher_SubscribeTicker_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &watcherSubscribeTickerClient{stream}
	return x, nil
}

type Watcher_SubscribeTickerClient interface {
	Send(*TickerRequest) error
	Recv() (*TickerResponse, error)
	grpc.ClientStream
}

type watcherSubscribeTickerClient struct {
	grpc.ClientStream
}

func (x *watcherSubscribeTickerClient) Send(m *TickerRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *watcherSubscribeTickerClient) Recv() (*TickerResponse, error) {
	m := new(TickerResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// WatcherServer is the server API for Watcher service.
// All implementations must embed UnimplementedWatcherServer
// for forward compatibility
type WatcherServer interface {
	GetInfo(context.Context, *TickerRequest) (*TickerResponse, error)
	MoreInfo(context.Context, *TickerRequest) (*CompanyResponse, error)
	SubscribeTicker(Watcher_SubscribeTickerServer) error
	mustEmbedUnimplementedWatcherServer()
}

// UnimplementedWatcherServer must be embedded to have forward compatible implementations.
type UnimplementedWatcherServer struct {
}

func (UnimplementedWatcherServer) GetInfo(context.Context, *TickerRequest) (*TickerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInfo not implemented")
}
func (UnimplementedWatcherServer) MoreInfo(context.Context, *TickerRequest) (*CompanyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MoreInfo not implemented")
}
func (UnimplementedWatcherServer) SubscribeTicker(Watcher_SubscribeTickerServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeTicker not implemented")
}
func (UnimplementedWatcherServer) mustEmbedUnimplementedWatcherServer() {}

// UnsafeWatcherServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WatcherServer will
// result in compilation errors.
type UnsafeWatcherServer interface {
	mustEmbedUnimplementedWatcherServer()
}

func RegisterWatcherServer(s grpc.ServiceRegistrar, srv WatcherServer) {
	s.RegisterService(&Watcher_ServiceDesc, srv)
}

func _Watcher_GetInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TickerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WatcherServer).GetInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Watcher_GetInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WatcherServer).GetInfo(ctx, req.(*TickerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Watcher_MoreInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TickerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WatcherServer).MoreInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Watcher_MoreInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WatcherServer).MoreInfo(ctx, req.(*TickerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Watcher_SubscribeTicker_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(WatcherServer).SubscribeTicker(&watcherSubscribeTickerServer{stream})
}

type Watcher_SubscribeTickerServer interface {
	Send(*TickerResponse) error
	Recv() (*TickerRequest, error)
	grpc.ServerStream
}

type watcherSubscribeTickerServer struct {
	grpc.ServerStream
}

func (x *watcherSubscribeTickerServer) Send(m *TickerResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *watcherSubscribeTickerServer) Recv() (*TickerRequest, error) {
	m := new(TickerRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Watcher_ServiceDesc is the grpc.ServiceDesc for Watcher service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Watcher_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Watcher",
	HandlerType: (*WatcherServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetInfo",
			Handler:    _Watcher_GetInfo_Handler,
		},
		{
			MethodName: "MoreInfo",
			Handler:    _Watcher_MoreInfo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeTicker",
			Handler:       _Watcher_SubscribeTicker_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "protos/watcher.proto",
}
