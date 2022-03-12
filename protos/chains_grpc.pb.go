// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

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

// ChainServiceClient is the client API for ChainService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChainServiceClient interface {
	GetChainInfo(ctx context.Context, in *GetChainInfoRequest, opts ...grpc.CallOption) (*GetChainInfoResponse, error)
	GetChain(ctx context.Context, in *GetChainRequest, opts ...grpc.CallOption) (*GetChainResponse, error)
}

type chainServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChainServiceClient(cc grpc.ClientConnInterface) ChainServiceClient {
	return &chainServiceClient{cc}
}

func (c *chainServiceClient) GetChainInfo(ctx context.Context, in *GetChainInfoRequest, opts ...grpc.CallOption) (*GetChainInfoResponse, error) {
	out := new(GetChainInfoResponse)
	err := c.cc.Invoke(ctx, "/protos.ChainService/GetChainInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chainServiceClient) GetChain(ctx context.Context, in *GetChainRequest, opts ...grpc.CallOption) (*GetChainResponse, error) {
	out := new(GetChainResponse)
	err := c.cc.Invoke(ctx, "/protos.ChainService/GetChain", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChainServiceServer is the server API for ChainService service.
// All implementations must embed UnimplementedChainServiceServer
// for forward compatibility
type ChainServiceServer interface {
	GetChainInfo(context.Context, *GetChainInfoRequest) (*GetChainInfoResponse, error)
	GetChain(context.Context, *GetChainRequest) (*GetChainResponse, error)
	mustEmbedUnimplementedChainServiceServer()
}

// UnimplementedChainServiceServer must be embedded to have forward compatible implementations.
type UnimplementedChainServiceServer struct {
}

func (UnimplementedChainServiceServer) GetChainInfo(context.Context, *GetChainInfoRequest) (*GetChainInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChainInfo not implemented")
}
func (UnimplementedChainServiceServer) GetChain(context.Context, *GetChainRequest) (*GetChainResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChain not implemented")
}
func (UnimplementedChainServiceServer) mustEmbedUnimplementedChainServiceServer() {}

// UnsafeChainServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChainServiceServer will
// result in compilation errors.
type UnsafeChainServiceServer interface {
	mustEmbedUnimplementedChainServiceServer()
}

func RegisterChainServiceServer(s grpc.ServiceRegistrar, srv ChainServiceServer) {
	s.RegisterService(&ChainService_ServiceDesc, srv)
}

func _ChainService_GetChainInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChainInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChainServiceServer).GetChainInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.ChainService/GetChainInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChainServiceServer).GetChainInfo(ctx, req.(*GetChainInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChainService_GetChain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChainRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChainServiceServer).GetChain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.ChainService/GetChain",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChainServiceServer).GetChain(ctx, req.(*GetChainRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ChainService_ServiceDesc is the grpc.ServiceDesc for ChainService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChainService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protos.ChainService",
	HandlerType: (*ChainServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetChainInfo",
			Handler:    _ChainService_GetChainInfo_Handler,
		},
		{
			MethodName: "GetChain",
			Handler:    _ChainService_GetChain_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/chains.proto",
}
