// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// HumorAgentClient is the client API for HumorAgent service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HumorAgentClient interface {
	Tts(ctx context.Context, in *TtsRequest, opts ...grpc.CallOption) (*TtsReply, error)
}

type humorAgentClient struct {
	cc grpc.ClientConnInterface
}

func NewHumorAgentClient(cc grpc.ClientConnInterface) HumorAgentClient {
	return &humorAgentClient{cc}
}

func (c *humorAgentClient) Tts(ctx context.Context, in *TtsRequest, opts ...grpc.CallOption) (*TtsReply, error) {
	out := new(TtsReply)
	err := c.cc.Invoke(ctx, "/api.HumorAgent/Tts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HumorAgentServer is the server API for HumorAgent service.
// All implementations must embed UnimplementedHumorAgentServer
// for forward compatibility
type HumorAgentServer interface {
	Tts(context.Context, *TtsRequest) (*TtsReply, error)
	mustEmbedUnimplementedHumorAgentServer()
}

// UnimplementedHumorAgentServer must be embedded to have forward compatible implementations.
type UnimplementedHumorAgentServer struct {
}

func (UnimplementedHumorAgentServer) Tts(context.Context, *TtsRequest) (*TtsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Tts not implemented")
}
func (UnimplementedHumorAgentServer) mustEmbedUnimplementedHumorAgentServer() {}

// UnsafeHumorAgentServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HumorAgentServer will
// result in compilation errors.
type UnsafeHumorAgentServer interface {
	mustEmbedUnimplementedHumorAgentServer()
}

func RegisterHumorAgentServer(s grpc.ServiceRegistrar, srv HumorAgentServer) {
	s.RegisterService(&_HumorAgent_serviceDesc, srv)
}

func _HumorAgent_Tts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TtsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HumorAgentServer).Tts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.HumorAgent/Tts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HumorAgentServer).Tts(ctx, req.(*TtsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _HumorAgent_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.HumorAgent",
	HandlerType: (*HumorAgentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Tts",
			Handler:    _HumorAgent_Tts_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/humoragent.proto",
}