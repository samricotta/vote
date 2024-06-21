// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: samricotta/vote/vote/v1/tx.proto

package votev1

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
	Msg_NewVote_FullMethodName     = "/samricotta.vote.vote.v1.Msg/NewVote"
	Msg_ResolveVote_FullMethodName = "/samricotta.vote.vote.v1.Msg/ResolveVote"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	NewVote(ctx context.Context, in *MsgNewVote, opts ...grpc.CallOption) (*MsgNewVoteResponse, error)
	ResolveVote(ctx context.Context, in *MsgResolveVote, opts ...grpc.CallOption) (*ResolveVoteResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) NewVote(ctx context.Context, in *MsgNewVote, opts ...grpc.CallOption) (*MsgNewVoteResponse, error) {
	out := new(MsgNewVoteResponse)
	err := c.cc.Invoke(ctx, Msg_NewVote_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) ResolveVote(ctx context.Context, in *MsgResolveVote, opts ...grpc.CallOption) (*ResolveVoteResponse, error) {
	out := new(ResolveVoteResponse)
	err := c.cc.Invoke(ctx, Msg_ResolveVote_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	NewVote(context.Context, *MsgNewVote) (*MsgNewVoteResponse, error)
	ResolveVote(context.Context, *MsgResolveVote) (*ResolveVoteResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) NewVote(context.Context, *MsgNewVote) (*MsgNewVoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewVote not implemented")
}
func (UnimplementedMsgServer) ResolveVote(context.Context, *MsgResolveVote) (*ResolveVoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResolveVote not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}

// UnsafeMsgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServer will
// result in compilation errors.
type UnsafeMsgServer interface {
	mustEmbedUnimplementedMsgServer()
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_NewVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgNewVote)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).NewVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_NewVote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).NewVote(ctx, req.(*MsgNewVote))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ResolveVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgResolveVote)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ResolveVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_ResolveVote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ResolveVote(ctx, req.(*MsgResolveVote))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "samricotta.vote.vote.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewVote",
			Handler:    _Msg_NewVote_Handler,
		},
		{
			MethodName: "ResolveVote",
			Handler:    _Msg_ResolveVote_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "samricotta/vote/vote/v1/tx.proto",
}