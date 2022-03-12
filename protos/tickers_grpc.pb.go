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

// TickerServiceClient is the client API for TickerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TickerServiceClient interface {
	GetTickers(ctx context.Context, in *GetTickersRequest, opts ...grpc.CallOption) (*GetTickersResponse, error)
}

type tickerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTickerServiceClient(cc grpc.ClientConnInterface) TickerServiceClient {
	return &tickerServiceClient{cc}
}

func (c *tickerServiceClient) GetTickers(ctx context.Context, in *GetTickersRequest, opts ...grpc.CallOption) (*GetTickersResponse, error) {
	out := new(GetTickersResponse)
	err := c.cc.Invoke(ctx, "/protos.TickerService/GetTickers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TickerServiceServer is the server API for TickerService service.
// All implementations must embed UnimplementedTickerServiceServer
// for forward compatibility
type TickerServiceServer interface {
	GetTickers(context.Context, *GetTickersRequest) (*GetTickersResponse, error)
	mustEmbedUnimplementedTickerServiceServer()
}

// UnimplementedTickerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTickerServiceServer struct {
}

func (UnimplementedTickerServiceServer) GetTickers(context.Context, *GetTickersRequest) (*GetTickersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTickers not implemented")
}
func (UnimplementedTickerServiceServer) mustEmbedUnimplementedTickerServiceServer() {}

// UnsafeTickerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TickerServiceServer will
// result in compilation errors.
type UnsafeTickerServiceServer interface {
	mustEmbedUnimplementedTickerServiceServer()
}

func RegisterTickerServiceServer(s grpc.ServiceRegistrar, srv TickerServiceServer) {
	s.RegisterService(&TickerService_ServiceDesc, srv)
}

func _TickerService_GetTickers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTickersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TickerServiceServer).GetTickers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.TickerService/GetTickers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TickerServiceServer).GetTickers(ctx, req.(*GetTickersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TickerService_ServiceDesc is the grpc.ServiceDesc for TickerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TickerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protos.TickerService",
	HandlerType: (*TickerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTickers",
			Handler:    _TickerService_GetTickers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/tickers.proto",
}
