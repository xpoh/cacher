// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.2
// source: api/cacher.proto

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

// CacherClient is the client API for Cacher service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CacherClient interface {
	GetRandomDataStream(ctx context.Context, in *Request, opts ...grpc.CallOption) (Cacher_GetRandomDataStreamClient, error)
}

type cacherClient struct {
	cc grpc.ClientConnInterface
}

func NewCacherClient(cc grpc.ClientConnInterface) CacherClient {
	return &cacherClient{cc}
}

func (c *cacherClient) GetRandomDataStream(ctx context.Context, in *Request, opts ...grpc.CallOption) (Cacher_GetRandomDataStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Cacher_ServiceDesc.Streams[0], "/cacher.Cacher/GetRandomDataStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &cacherGetRandomDataStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Cacher_GetRandomDataStreamClient interface {
	Recv() (*Response, error)
	grpc.ClientStream
}

type cacherGetRandomDataStreamClient struct {
	grpc.ClientStream
}

func (x *cacherGetRandomDataStreamClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CacherServer is the server API for Cacher service.
// All implementations must embed UnimplementedCacherServer
// for forward compatibility
type CacherServer interface {
	GetRandomDataStream(*Request, Cacher_GetRandomDataStreamServer) error
	mustEmbedUnimplementedCacherServer()
}

// UnimplementedCacherServer must be embedded to have forward compatible implementations.
type UnimplementedCacherServer struct {
}

func (UnimplementedCacherServer) GetRandomDataStream(*Request, Cacher_GetRandomDataStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method GetRandomDataStream not implemented")
}
func (UnimplementedCacherServer) mustEmbedUnimplementedCacherServer() {}

// UnsafeCacherServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CacherServer will
// result in compilation errors.
type UnsafeCacherServer interface {
	mustEmbedUnimplementedCacherServer()
}

func RegisterCacherServer(s grpc.ServiceRegistrar, srv CacherServer) {
	s.RegisterService(&Cacher_ServiceDesc, srv)
}

func _Cacher_GetRandomDataStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CacherServer).GetRandomDataStream(m, &cacherGetRandomDataStreamServer{stream})
}

type Cacher_GetRandomDataStreamServer interface {
	Send(*Response) error
	grpc.ServerStream
}

type cacherGetRandomDataStreamServer struct {
	grpc.ServerStream
}

func (x *cacherGetRandomDataStreamServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

// Cacher_ServiceDesc is the grpc.ServiceDesc for Cacher service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Cacher_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cacher.Cacher",
	HandlerType: (*CacherServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetRandomDataStream",
			Handler:       _Cacher_GetRandomDataStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/cacher.proto",
}