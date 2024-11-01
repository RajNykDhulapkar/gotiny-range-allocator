// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v4.25.1
// source: allocator_svc.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	RangeAllocator_AllocateRange_FullMethodName     = "/rangeallocator.v1.RangeAllocator/AllocateRange"
	RangeAllocator_GetRange_FullMethodName          = "/rangeallocator.v1.RangeAllocator/GetRange"
	RangeAllocator_UpdateRangeStatus_FullMethodName = "/rangeallocator.v1.RangeAllocator/UpdateRangeStatus"
	RangeAllocator_GetHealth_FullMethodName         = "/rangeallocator.v1.RangeAllocator/GetHealth"
)

// RangeAllocatorClient is the client API for RangeAllocator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RangeAllocatorClient interface {
	AllocateRange(ctx context.Context, in *AllocateRangeRequest, opts ...grpc.CallOption) (*AllocateRangeResponse, error)
	GetRange(ctx context.Context, in *GetRangeRequest, opts ...grpc.CallOption) (*Range, error)
	UpdateRangeStatus(ctx context.Context, in *UpdateRangeStatusRequest, opts ...grpc.CallOption) (*Range, error)
	GetHealth(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*HealthResponse, error)
}

type rangeAllocatorClient struct {
	cc grpc.ClientConnInterface
}

func NewRangeAllocatorClient(cc grpc.ClientConnInterface) RangeAllocatorClient {
	return &rangeAllocatorClient{cc}
}

func (c *rangeAllocatorClient) AllocateRange(ctx context.Context, in *AllocateRangeRequest, opts ...grpc.CallOption) (*AllocateRangeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AllocateRangeResponse)
	err := c.cc.Invoke(ctx, RangeAllocator_AllocateRange_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rangeAllocatorClient) GetRange(ctx context.Context, in *GetRangeRequest, opts ...grpc.CallOption) (*Range, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Range)
	err := c.cc.Invoke(ctx, RangeAllocator_GetRange_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rangeAllocatorClient) UpdateRangeStatus(ctx context.Context, in *UpdateRangeStatusRequest, opts ...grpc.CallOption) (*Range, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Range)
	err := c.cc.Invoke(ctx, RangeAllocator_UpdateRangeStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rangeAllocatorClient) GetHealth(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*HealthResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HealthResponse)
	err := c.cc.Invoke(ctx, RangeAllocator_GetHealth_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RangeAllocatorServer is the server API for RangeAllocator service.
// All implementations must embed UnimplementedRangeAllocatorServer
// for forward compatibility.
type RangeAllocatorServer interface {
	AllocateRange(context.Context, *AllocateRangeRequest) (*AllocateRangeResponse, error)
	GetRange(context.Context, *GetRangeRequest) (*Range, error)
	UpdateRangeStatus(context.Context, *UpdateRangeStatusRequest) (*Range, error)
	GetHealth(context.Context, *emptypb.Empty) (*HealthResponse, error)
	mustEmbedUnimplementedRangeAllocatorServer()
}

// UnimplementedRangeAllocatorServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRangeAllocatorServer struct{}

func (UnimplementedRangeAllocatorServer) AllocateRange(context.Context, *AllocateRangeRequest) (*AllocateRangeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllocateRange not implemented")
}
func (UnimplementedRangeAllocatorServer) GetRange(context.Context, *GetRangeRequest) (*Range, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRange not implemented")
}
func (UnimplementedRangeAllocatorServer) UpdateRangeStatus(context.Context, *UpdateRangeStatusRequest) (*Range, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRangeStatus not implemented")
}
func (UnimplementedRangeAllocatorServer) GetHealth(context.Context, *emptypb.Empty) (*HealthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHealth not implemented")
}
func (UnimplementedRangeAllocatorServer) mustEmbedUnimplementedRangeAllocatorServer() {}
func (UnimplementedRangeAllocatorServer) testEmbeddedByValue()                        {}

// UnsafeRangeAllocatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RangeAllocatorServer will
// result in compilation errors.
type UnsafeRangeAllocatorServer interface {
	mustEmbedUnimplementedRangeAllocatorServer()
}

func RegisterRangeAllocatorServer(s grpc.ServiceRegistrar, srv RangeAllocatorServer) {
	// If the following call pancis, it indicates UnimplementedRangeAllocatorServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RangeAllocator_ServiceDesc, srv)
}

func _RangeAllocator_AllocateRange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AllocateRangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RangeAllocatorServer).AllocateRange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RangeAllocator_AllocateRange_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RangeAllocatorServer).AllocateRange(ctx, req.(*AllocateRangeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RangeAllocator_GetRange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RangeAllocatorServer).GetRange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RangeAllocator_GetRange_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RangeAllocatorServer).GetRange(ctx, req.(*GetRangeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RangeAllocator_UpdateRangeStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRangeStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RangeAllocatorServer).UpdateRangeStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RangeAllocator_UpdateRangeStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RangeAllocatorServer).UpdateRangeStatus(ctx, req.(*UpdateRangeStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RangeAllocator_GetHealth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RangeAllocatorServer).GetHealth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RangeAllocator_GetHealth_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RangeAllocatorServer).GetHealth(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// RangeAllocator_ServiceDesc is the grpc.ServiceDesc for RangeAllocator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RangeAllocator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rangeallocator.v1.RangeAllocator",
	HandlerType: (*RangeAllocatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AllocateRange",
			Handler:    _RangeAllocator_AllocateRange_Handler,
		},
		{
			MethodName: "GetRange",
			Handler:    _RangeAllocator_GetRange_Handler,
		},
		{
			MethodName: "UpdateRangeStatus",
			Handler:    _RangeAllocator_UpdateRangeStatus_Handler,
		},
		{
			MethodName: "GetHealth",
			Handler:    _RangeAllocator_GetHealth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "allocator_svc.proto",
}
