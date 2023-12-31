// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: proto/service.proto

package proto

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

// BroadcastClient is the client API for Broadcast service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BroadcastClient interface {
	CreateStream(ctx context.Context, in *Connect, opts ...grpc.CallOption) (Broadcast_CreateStreamClient, error)
	BroadcastMessage(ctx context.Context, opts ...grpc.CallOption) (Broadcast_BroadcastMessageClient, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	GetInfo(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*InfoResponse, error)
	Follow(ctx context.Context, opts ...grpc.CallOption) (Broadcast_FollowClient, error)
	UserExist(ctx context.Context, in *UserExistRequest, opts ...grpc.CallOption) (*UserExistResponse, error)
}

type broadcastClient struct {
	cc grpc.ClientConnInterface
}

func NewBroadcastClient(cc grpc.ClientConnInterface) BroadcastClient {
	return &broadcastClient{cc}
}

func (c *broadcastClient) CreateStream(ctx context.Context, in *Connect, opts ...grpc.CallOption) (Broadcast_CreateStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Broadcast_ServiceDesc.Streams[0], "/proto.Broadcast/CreateStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &broadcastCreateStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Broadcast_CreateStreamClient interface {
	Recv() (*ServerResponse, error)
	grpc.ClientStream
}

type broadcastCreateStreamClient struct {
	grpc.ClientStream
}

func (x *broadcastCreateStreamClient) Recv() (*ServerResponse, error) {
	m := new(ServerResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *broadcastClient) BroadcastMessage(ctx context.Context, opts ...grpc.CallOption) (Broadcast_BroadcastMessageClient, error) {
	stream, err := c.cc.NewStream(ctx, &Broadcast_ServiceDesc.Streams[1], "/proto.Broadcast/BroadcastMessage", opts...)
	if err != nil {
		return nil, err
	}
	x := &broadcastBroadcastMessageClient{stream}
	return x, nil
}

type Broadcast_BroadcastMessageClient interface {
	Send(*MessageRequest) error
	CloseAndRecv() (*Close, error)
	grpc.ClientStream
}

type broadcastBroadcastMessageClient struct {
	grpc.ClientStream
}

func (x *broadcastBroadcastMessageClient) Send(m *MessageRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *broadcastBroadcastMessageClient) CloseAndRecv() (*Close, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Close)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *broadcastClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/proto.Broadcast/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *broadcastClient) GetInfo(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*InfoResponse, error) {
	out := new(InfoResponse)
	err := c.cc.Invoke(ctx, "/proto.Broadcast/GetInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *broadcastClient) Follow(ctx context.Context, opts ...grpc.CallOption) (Broadcast_FollowClient, error) {
	stream, err := c.cc.NewStream(ctx, &Broadcast_ServiceDesc.Streams[2], "/proto.Broadcast/Follow", opts...)
	if err != nil {
		return nil, err
	}
	x := &broadcastFollowClient{stream}
	return x, nil
}

type Broadcast_FollowClient interface {
	Send(*FollowerRequest) error
	Recv() (*LeaderResponse, error)
	grpc.ClientStream
}

type broadcastFollowClient struct {
	grpc.ClientStream
}

func (x *broadcastFollowClient) Send(m *FollowerRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *broadcastFollowClient) Recv() (*LeaderResponse, error) {
	m := new(LeaderResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *broadcastClient) UserExist(ctx context.Context, in *UserExistRequest, opts ...grpc.CallOption) (*UserExistResponse, error) {
	out := new(UserExistResponse)
	err := c.cc.Invoke(ctx, "/proto.Broadcast/UserExist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BroadcastServer is the server API for Broadcast service.
// All implementations should embed UnimplementedBroadcastServer
// for forward compatibility
type BroadcastServer interface {
	CreateStream(*Connect, Broadcast_CreateStreamServer) error
	BroadcastMessage(Broadcast_BroadcastMessageServer) error
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	GetInfo(context.Context, *EmptyRequest) (*InfoResponse, error)
	Follow(Broadcast_FollowServer) error
	UserExist(context.Context, *UserExistRequest) (*UserExistResponse, error)
}

// UnimplementedBroadcastServer should be embedded to have forward compatible implementations.
type UnimplementedBroadcastServer struct {
}

func (UnimplementedBroadcastServer) CreateStream(*Connect, Broadcast_CreateStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method CreateStream not implemented")
}
func (UnimplementedBroadcastServer) BroadcastMessage(Broadcast_BroadcastMessageServer) error {
	return status.Errorf(codes.Unimplemented, "method BroadcastMessage not implemented")
}
func (UnimplementedBroadcastServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedBroadcastServer) GetInfo(context.Context, *EmptyRequest) (*InfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInfo not implemented")
}
func (UnimplementedBroadcastServer) Follow(Broadcast_FollowServer) error {
	return status.Errorf(codes.Unimplemented, "method Follow not implemented")
}
func (UnimplementedBroadcastServer) UserExist(context.Context, *UserExistRequest) (*UserExistResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserExist not implemented")
}

// UnsafeBroadcastServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BroadcastServer will
// result in compilation errors.
type UnsafeBroadcastServer interface {
	mustEmbedUnimplementedBroadcastServer()
}

func RegisterBroadcastServer(s grpc.ServiceRegistrar, srv BroadcastServer) {
	s.RegisterService(&Broadcast_ServiceDesc, srv)
}

func _Broadcast_CreateStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Connect)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BroadcastServer).CreateStream(m, &broadcastCreateStreamServer{stream})
}

type Broadcast_CreateStreamServer interface {
	Send(*ServerResponse) error
	grpc.ServerStream
}

type broadcastCreateStreamServer struct {
	grpc.ServerStream
}

func (x *broadcastCreateStreamServer) Send(m *ServerResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Broadcast_BroadcastMessage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BroadcastServer).BroadcastMessage(&broadcastBroadcastMessageServer{stream})
}

type Broadcast_BroadcastMessageServer interface {
	SendAndClose(*Close) error
	Recv() (*MessageRequest, error)
	grpc.ServerStream
}

type broadcastBroadcastMessageServer struct {
	grpc.ServerStream
}

func (x *broadcastBroadcastMessageServer) SendAndClose(m *Close) error {
	return x.ServerStream.SendMsg(m)
}

func (x *broadcastBroadcastMessageServer) Recv() (*MessageRequest, error) {
	m := new(MessageRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Broadcast_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BroadcastServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Broadcast/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BroadcastServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broadcast_GetInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BroadcastServer).GetInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Broadcast/GetInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BroadcastServer).GetInfo(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broadcast_Follow_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BroadcastServer).Follow(&broadcastFollowServer{stream})
}

type Broadcast_FollowServer interface {
	Send(*LeaderResponse) error
	Recv() (*FollowerRequest, error)
	grpc.ServerStream
}

type broadcastFollowServer struct {
	grpc.ServerStream
}

func (x *broadcastFollowServer) Send(m *LeaderResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *broadcastFollowServer) Recv() (*FollowerRequest, error) {
	m := new(FollowerRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Broadcast_UserExist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserExistRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BroadcastServer).UserExist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Broadcast/UserExist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BroadcastServer).UserExist(ctx, req.(*UserExistRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Broadcast_ServiceDesc is the grpc.ServiceDesc for Broadcast service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Broadcast_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Broadcast",
	HandlerType: (*BroadcastServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Broadcast_Ping_Handler,
		},
		{
			MethodName: "GetInfo",
			Handler:    _Broadcast_GetInfo_Handler,
		},
		{
			MethodName: "UserExist",
			Handler:    _Broadcast_UserExist_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "CreateStream",
			Handler:       _Broadcast_CreateStream_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "BroadcastMessage",
			Handler:       _Broadcast_BroadcastMessage_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Follow",
			Handler:       _Broadcast_Follow_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/service.proto",
}
