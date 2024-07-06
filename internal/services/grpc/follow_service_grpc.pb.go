// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: internal/services/proto/follow_service.proto

package socialnetwork_pb

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

// FollowServiceClient is the client API for FollowService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FollowServiceClient interface {
	// Follow a user
	FollowUser(ctx context.Context, in *FollowUserRequest, opts ...grpc.CallOption) (*FollowUserResponse, error)
	// Unfollow a user
	UnfollowUser(ctx context.Context, in *UnfollowUserRequest, opts ...grpc.CallOption) (*UnfollowUserResponse, error)
	// Get the list of users followed by a specific user
	GetFollowing(ctx context.Context, in *GetFollowingRequest, opts ...grpc.CallOption) (*GetFollowingResponse, error)
}

type followServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFollowServiceClient(cc grpc.ClientConnInterface) FollowServiceClient {
	return &followServiceClient{cc}
}

func (c *followServiceClient) FollowUser(ctx context.Context, in *FollowUserRequest, opts ...grpc.CallOption) (*FollowUserResponse, error) {
	out := new(FollowUserResponse)
	err := c.cc.Invoke(ctx, "/socialnetwork.FollowService/FollowUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *followServiceClient) UnfollowUser(ctx context.Context, in *UnfollowUserRequest, opts ...grpc.CallOption) (*UnfollowUserResponse, error) {
	out := new(UnfollowUserResponse)
	err := c.cc.Invoke(ctx, "/socialnetwork.FollowService/UnfollowUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *followServiceClient) GetFollowing(ctx context.Context, in *GetFollowingRequest, opts ...grpc.CallOption) (*GetFollowingResponse, error) {
	out := new(GetFollowingResponse)
	err := c.cc.Invoke(ctx, "/socialnetwork.FollowService/GetFollowing", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FollowServiceServer is the server API for FollowService service.
// All implementations must embed UnimplementedFollowServiceServer
// for forward compatibility
type FollowServiceServer interface {
	// Follow a user
	FollowUser(context.Context, *FollowUserRequest) (*FollowUserResponse, error)
	// Unfollow a user
	UnfollowUser(context.Context, *UnfollowUserRequest) (*UnfollowUserResponse, error)
	// Get the list of users followed by a specific user
	GetFollowing(context.Context, *GetFollowingRequest) (*GetFollowingResponse, error)
	mustEmbedUnimplementedFollowServiceServer()
}

// UnimplementedFollowServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFollowServiceServer struct {
}

func (UnimplementedFollowServiceServer) FollowUser(context.Context, *FollowUserRequest) (*FollowUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FollowUser not implemented")
}
func (UnimplementedFollowServiceServer) UnfollowUser(context.Context, *UnfollowUserRequest) (*UnfollowUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnfollowUser not implemented")
}
func (UnimplementedFollowServiceServer) GetFollowing(context.Context, *GetFollowingRequest) (*GetFollowingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFollowing not implemented")
}
func (UnimplementedFollowServiceServer) mustEmbedUnimplementedFollowServiceServer() {}

// UnsafeFollowServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FollowServiceServer will
// result in compilation errors.
type UnsafeFollowServiceServer interface {
	mustEmbedUnimplementedFollowServiceServer()
}

func RegisterFollowServiceServer(s grpc.ServiceRegistrar, srv FollowServiceServer) {
	s.RegisterService(&FollowService_ServiceDesc, srv)
}

func _FollowService_FollowUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FollowUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FollowServiceServer).FollowUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/socialnetwork.FollowService/FollowUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FollowServiceServer).FollowUser(ctx, req.(*FollowUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FollowService_UnfollowUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnfollowUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FollowServiceServer).UnfollowUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/socialnetwork.FollowService/UnfollowUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FollowServiceServer).UnfollowUser(ctx, req.(*UnfollowUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FollowService_GetFollowing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFollowingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FollowServiceServer).GetFollowing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/socialnetwork.FollowService/GetFollowing",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FollowServiceServer).GetFollowing(ctx, req.(*GetFollowingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FollowService_ServiceDesc is the grpc.ServiceDesc for FollowService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FollowService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "socialnetwork.FollowService",
	HandlerType: (*FollowServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FollowUser",
			Handler:    _FollowService_FollowUser_Handler,
		},
		{
			MethodName: "UnfollowUser",
			Handler:    _FollowService_UnfollowUser_Handler,
		},
		{
			MethodName: "GetFollowing",
			Handler:    _FollowService_GetFollowing_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/services/proto/follow_service.proto",
}
