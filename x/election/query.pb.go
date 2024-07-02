// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: samricotta/election/v1/query.proto

package election

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

func init() {
	proto.RegisterFile("samricotta/election/v1/query.proto", fileDescriptor_01cf0d11eb1370f3)
}

var fileDescriptor_01cf0d11eb1370f3 = []byte{
	// 192 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2a, 0x4e, 0xcc, 0x2d,
	0xca, 0x4c, 0xce, 0x2f, 0x29, 0x49, 0xd4, 0x4f, 0xcd, 0x49, 0x4d, 0x2e, 0xc9, 0xcc, 0xcf, 0xd3,
	0x2f, 0x33, 0xd4, 0x2f, 0x2c, 0x4d, 0x2d, 0xaa, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12,
	0x43, 0xa8, 0xd1, 0x83, 0xa9, 0xd1, 0x2b, 0x33, 0x94, 0x92, 0x49, 0xcf, 0xcf, 0x4f, 0xcf, 0x49,
	0xd5, 0x4f, 0x2c, 0xc8, 0xd4, 0x4f, 0xcc, 0xcb, 0xcb, 0x2f, 0x49, 0x04, 0xc9, 0x14, 0x43, 0x74,
	0x49, 0x49, 0x27, 0xe7, 0x17, 0xe7, 0xe6, 0x17, 0x43, 0x4c, 0x42, 0x33, 0x52, 0x4a, 0x30, 0x31,
	0x37, 0x33, 0x2f, 0x5f, 0x1f, 0x4c, 0x42, 0x85, 0x44, 0xd2, 0xf3, 0xd3, 0xf3, 0xc1, 0x4c, 0x7d,
	0x10, 0x0b, 0x22, 0x6a, 0xc4, 0xce, 0xc5, 0x1a, 0x08, 0xd2, 0xe7, 0x64, 0x7f, 0xe2, 0x91, 0x1c,
	0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f, 0x1e, 0xc9, 0x31, 0x4e, 0x78, 0x2c, 0xc7, 0x70, 0xe1, 0xb1,
	0x1c, 0xc3, 0x8d, 0xc7, 0x72, 0x0c, 0x51, 0xaa, 0xe9, 0x99, 0x25, 0x19, 0xa5, 0x49, 0x7a, 0xc9,
	0xf9, 0xb9, 0xfa, 0x48, 0xbe, 0x29, 0xcb, 0x2f, 0x49, 0xd5, 0xaf, 0x80, 0x7b, 0x2a, 0x89, 0x0d,
	0x6c, 0xa0, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xe6, 0x1d, 0x2a, 0x66, 0xf2, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

// QueryServer is the server API for Query service.
type QueryServer interface {
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "samricotta.election.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "samricotta/election/v1/query.proto",
}