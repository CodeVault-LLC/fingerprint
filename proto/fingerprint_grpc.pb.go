// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.2
// source: proto/fingerprint.proto

package fingerprint

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	FingerprintService_AddFingerprint_FullMethodName   = "/fingerprint.FingerprintService/AddFingerprint"
	FingerprintService_GetFingerprint_FullMethodName   = "/fingerprint.FingerprintService/GetFingerprint"
	FingerprintService_MatchFingerprint_FullMethodName = "/fingerprint.FingerprintService/MatchFingerprint"
)

// FingerprintServiceClient is the client API for FingerprintService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FingerprintServiceClient interface {
	AddFingerprint(ctx context.Context, in *AddFingerprintRequest, opts ...grpc.CallOption) (*AddFingerprintResponse, error)
	GetFingerprint(ctx context.Context, in *GetFingerprintRequest, opts ...grpc.CallOption) (*GetFingerprintResponse, error)
	MatchFingerprint(ctx context.Context, in *MatchFingerprintRequest, opts ...grpc.CallOption) (*MatchFingerprintResponse, error)
}

type fingerprintServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFingerprintServiceClient(cc grpc.ClientConnInterface) FingerprintServiceClient {
	return &fingerprintServiceClient{cc}
}

func (c *fingerprintServiceClient) AddFingerprint(ctx context.Context, in *AddFingerprintRequest, opts ...grpc.CallOption) (*AddFingerprintResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddFingerprintResponse)
	err := c.cc.Invoke(ctx, FingerprintService_AddFingerprint_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fingerprintServiceClient) GetFingerprint(ctx context.Context, in *GetFingerprintRequest, opts ...grpc.CallOption) (*GetFingerprintResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetFingerprintResponse)
	err := c.cc.Invoke(ctx, FingerprintService_GetFingerprint_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fingerprintServiceClient) MatchFingerprint(ctx context.Context, in *MatchFingerprintRequest, opts ...grpc.CallOption) (*MatchFingerprintResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MatchFingerprintResponse)
	err := c.cc.Invoke(ctx, FingerprintService_MatchFingerprint_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FingerprintServiceServer is the server API for FingerprintService service.
// All implementations must embed UnimplementedFingerprintServiceServer
// for forward compatibility.
type FingerprintServiceServer interface {
	AddFingerprint(context.Context, *AddFingerprintRequest) (*AddFingerprintResponse, error)
	GetFingerprint(context.Context, *GetFingerprintRequest) (*GetFingerprintResponse, error)
	MatchFingerprint(context.Context, *MatchFingerprintRequest) (*MatchFingerprintResponse, error)
	mustEmbedUnimplementedFingerprintServiceServer()
}

// UnimplementedFingerprintServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedFingerprintServiceServer struct{}

func (UnimplementedFingerprintServiceServer) AddFingerprint(context.Context, *AddFingerprintRequest) (*AddFingerprintResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddFingerprint not implemented")
}
func (UnimplementedFingerprintServiceServer) GetFingerprint(context.Context, *GetFingerprintRequest) (*GetFingerprintResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFingerprint not implemented")
}
func (UnimplementedFingerprintServiceServer) MatchFingerprint(context.Context, *MatchFingerprintRequest) (*MatchFingerprintResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MatchFingerprint not implemented")
}
func (UnimplementedFingerprintServiceServer) mustEmbedUnimplementedFingerprintServiceServer() {}
func (UnimplementedFingerprintServiceServer) testEmbeddedByValue()                            {}

// UnsafeFingerprintServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FingerprintServiceServer will
// result in compilation errors.
type UnsafeFingerprintServiceServer interface {
	mustEmbedUnimplementedFingerprintServiceServer()
}

func RegisterFingerprintServiceServer(s grpc.ServiceRegistrar, srv FingerprintServiceServer) {
	// If the following call pancis, it indicates UnimplementedFingerprintServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&FingerprintService_ServiceDesc, srv)
}

func _FingerprintService_AddFingerprint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddFingerprintRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FingerprintServiceServer).AddFingerprint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FingerprintService_AddFingerprint_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FingerprintServiceServer).AddFingerprint(ctx, req.(*AddFingerprintRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FingerprintService_GetFingerprint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFingerprintRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FingerprintServiceServer).GetFingerprint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FingerprintService_GetFingerprint_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FingerprintServiceServer).GetFingerprint(ctx, req.(*GetFingerprintRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FingerprintService_MatchFingerprint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MatchFingerprintRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FingerprintServiceServer).MatchFingerprint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FingerprintService_MatchFingerprint_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FingerprintServiceServer).MatchFingerprint(ctx, req.(*MatchFingerprintRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FingerprintService_ServiceDesc is the grpc.ServiceDesc for FingerprintService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FingerprintService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fingerprint.FingerprintService",
	HandlerType: (*FingerprintServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddFingerprint",
			Handler:    _FingerprintService_AddFingerprint_Handler,
		},
		{
			MethodName: "GetFingerprint",
			Handler:    _FingerprintService_GetFingerprint_Handler,
		},
		{
			MethodName: "MatchFingerprint",
			Handler:    _FingerprintService_MatchFingerprint_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/fingerprint.proto",
}
