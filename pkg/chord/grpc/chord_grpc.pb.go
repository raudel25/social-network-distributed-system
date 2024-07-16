// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.12.4
// source: pkg/chord/grpc/chord.proto

package chord_pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	Chord_FindSuccessor_FullMethodName         = "/chord.Chord/FindSuccessor"
	Chord_GetPredecessor_FullMethodName        = "/chord.Chord/GetPredecessor"
	Chord_GetSuccessorAndNotify_FullMethodName = "/chord.Chord/GetSuccessorAndNotify"
	Chord_Notify_FullMethodName                = "/chord.Chord/Notify"
	Chord_Ping_FullMethodName                  = "/chord.Chord/Ping"
	Chord_PingLeader_FullMethodName            = "/chord.Chord/PingLeader"
	Chord_Election_FullMethodName              = "/chord.Chord/Election"
	Chord_Get_FullMethodName                   = "/chord.Chord/Get"
	Chord_Set_FullMethodName                   = "/chord.Chord/Set"
	Chord_SetPartition_FullMethodName          = "/chord.Chord/SetPartition"
	Chord_ResolveData_FullMethodName           = "/chord.Chord/ResolveData"
	Chord_Remove_FullMethodName                = "/chord.Chord/Remove"
)

// ChordClient is the client API for Chord service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChordClient interface {
	FindSuccessor(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*NodeResponse, error)
	GetPredecessor(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*NodeResponse, error)
	GetSuccessorAndNotify(ctx context.Context, in *NodeIndexRequest, opts ...grpc.CallOption) (*NodeResponse, error)
	Notify(ctx context.Context, in *NodeRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	Ping(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	PingLeader(ctx context.Context, in *TimeRequest, opts ...grpc.CallOption) (*TimeResponse, error)
	Election(ctx context.Context, in *ElectionRequest, opts ...grpc.CallOption) (*NodeResponse, error)
	Get(ctx context.Context, in *KeyRequest, opts ...grpc.CallOption) (*StatusValueResponse, error)
	Set(ctx context.Context, in *KeyValueRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	SetPartition(ctx context.Context, in *PartitionRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	ResolveData(ctx context.Context, in *PartitionRequest, opts ...grpc.CallOption) (*ResolveDataResponse, error)
	Remove(ctx context.Context, in *KeyRequest, opts ...grpc.CallOption) (*StatusResponse, error)
}

type chordClient struct {
	cc grpc.ClientConnInterface
}

func NewChordClient(cc grpc.ClientConnInterface) ChordClient {
	return &chordClient{cc}
}

func (c *chordClient) FindSuccessor(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*NodeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NodeResponse)
	err := c.cc.Invoke(ctx, Chord_FindSuccessor_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) GetPredecessor(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*NodeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NodeResponse)
	err := c.cc.Invoke(ctx, Chord_GetPredecessor_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) GetSuccessorAndNotify(ctx context.Context, in *NodeIndexRequest, opts ...grpc.CallOption) (*NodeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NodeResponse)
	err := c.cc.Invoke(ctx, Chord_GetSuccessorAndNotify_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) Notify(ctx context.Context, in *NodeRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, Chord_Notify_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) Ping(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, Chord_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) PingLeader(ctx context.Context, in *TimeRequest, opts ...grpc.CallOption) (*TimeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TimeResponse)
	err := c.cc.Invoke(ctx, Chord_PingLeader_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) Election(ctx context.Context, in *ElectionRequest, opts ...grpc.CallOption) (*NodeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NodeResponse)
	err := c.cc.Invoke(ctx, Chord_Election_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) Get(ctx context.Context, in *KeyRequest, opts ...grpc.CallOption) (*StatusValueResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatusValueResponse)
	err := c.cc.Invoke(ctx, Chord_Get_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) Set(ctx context.Context, in *KeyValueRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, Chord_Set_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) SetPartition(ctx context.Context, in *PartitionRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, Chord_SetPartition_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) ResolveData(ctx context.Context, in *PartitionRequest, opts ...grpc.CallOption) (*ResolveDataResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResolveDataResponse)
	err := c.cc.Invoke(ctx, Chord_ResolveData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chordClient) Remove(ctx context.Context, in *KeyRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, Chord_Remove_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChordServer is the server API for Chord service.
// All implementations must embed UnimplementedChordServer
// for forward compatibility
type ChordServer interface {
	FindSuccessor(context.Context, *IdRequest) (*NodeResponse, error)
	GetPredecessor(context.Context, *EmptyRequest) (*NodeResponse, error)
	GetSuccessorAndNotify(context.Context, *NodeIndexRequest) (*NodeResponse, error)
	Notify(context.Context, *NodeRequest) (*StatusResponse, error)
	Ping(context.Context, *EmptyRequest) (*StatusResponse, error)
	PingLeader(context.Context, *TimeRequest) (*TimeResponse, error)
	Election(context.Context, *ElectionRequest) (*NodeResponse, error)
	Get(context.Context, *KeyRequest) (*StatusValueResponse, error)
	Set(context.Context, *KeyValueRequest) (*StatusResponse, error)
	SetPartition(context.Context, *PartitionRequest) (*StatusResponse, error)
	ResolveData(context.Context, *PartitionRequest) (*ResolveDataResponse, error)
	Remove(context.Context, *KeyRequest) (*StatusResponse, error)
	mustEmbedUnimplementedChordServer()
}

// UnimplementedChordServer must be embedded to have forward compatible implementations.
type UnimplementedChordServer struct {
}

func (UnimplementedChordServer) FindSuccessor(context.Context, *IdRequest) (*NodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindSuccessor not implemented")
}
func (UnimplementedChordServer) GetPredecessor(context.Context, *EmptyRequest) (*NodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPredecessor not implemented")
}
func (UnimplementedChordServer) GetSuccessorAndNotify(context.Context, *NodeIndexRequest) (*NodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSuccessorAndNotify not implemented")
}
func (UnimplementedChordServer) Notify(context.Context, *NodeRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Notify not implemented")
}
func (UnimplementedChordServer) Ping(context.Context, *EmptyRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedChordServer) PingLeader(context.Context, *TimeRequest) (*TimeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PingLeader not implemented")
}
func (UnimplementedChordServer) Election(context.Context, *ElectionRequest) (*NodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Election not implemented")
}
func (UnimplementedChordServer) Get(context.Context, *KeyRequest) (*StatusValueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedChordServer) Set(context.Context, *KeyValueRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedChordServer) SetPartition(context.Context, *PartitionRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetPartition not implemented")
}
func (UnimplementedChordServer) ResolveData(context.Context, *PartitionRequest) (*ResolveDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResolveData not implemented")
}
func (UnimplementedChordServer) Remove(context.Context, *KeyRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}
func (UnimplementedChordServer) mustEmbedUnimplementedChordServer() {}

// UnsafeChordServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChordServer will
// result in compilation errors.
type UnsafeChordServer interface {
	mustEmbedUnimplementedChordServer()
}

func RegisterChordServer(s grpc.ServiceRegistrar, srv ChordServer) {
	s.RegisterService(&Chord_ServiceDesc, srv)
}

func _Chord_FindSuccessor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).FindSuccessor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chord_FindSuccessor_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).FindSuccessor(ctx, req.(*IdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_GetPredecessor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).GetPredecessor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chord_GetPredecessor_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).GetPredecessor(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_GetSuccessorAndNotify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NodeIndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).GetSuccessorAndNotify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chord_GetSuccessorAndNotify_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).GetSuccessorAndNotify(ctx, req.(*NodeIndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_Notify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).Notify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chord_Notify_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).Notify(ctx, req.(*NodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chord_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).Ping(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_PingLeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TimeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).PingLeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chord_PingLeader_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).PingLeader(ctx, req.(*TimeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_Election_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ElectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).Election(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chord_Election_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).Election(ctx, req.(*ElectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chord_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).Get(ctx, req.(*KeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyValueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chord_Set_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).Set(ctx, req.(*KeyValueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_SetPartition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PartitionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).SetPartition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chord_SetPartition_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).SetPartition(ctx, req.(*PartitionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_ResolveData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PartitionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).ResolveData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chord_ResolveData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).ResolveData(ctx, req.(*PartitionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chord_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChordServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chord_Remove_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChordServer).Remove(ctx, req.(*KeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Chord_ServiceDesc is the grpc.ServiceDesc for Chord service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Chord_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chord.Chord",
	HandlerType: (*ChordServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindSuccessor",
			Handler:    _Chord_FindSuccessor_Handler,
		},
		{
			MethodName: "GetPredecessor",
			Handler:    _Chord_GetPredecessor_Handler,
		},
		{
			MethodName: "GetSuccessorAndNotify",
			Handler:    _Chord_GetSuccessorAndNotify_Handler,
		},
		{
			MethodName: "Notify",
			Handler:    _Chord_Notify_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _Chord_Ping_Handler,
		},
		{
			MethodName: "PingLeader",
			Handler:    _Chord_PingLeader_Handler,
		},
		{
			MethodName: "Election",
			Handler:    _Chord_Election_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _Chord_Get_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _Chord_Set_Handler,
		},
		{
			MethodName: "SetPartition",
			Handler:    _Chord_SetPartition_Handler,
		},
		{
			MethodName: "ResolveData",
			Handler:    _Chord_ResolveData_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _Chord_Remove_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/chord/grpc/chord.proto",
}
