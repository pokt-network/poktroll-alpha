// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: poktroll/poktroll/query.proto

package types

import (
	context "context"
	fmt "fmt"
	query "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
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

// QueryParamsRequest is request type for the Query/Params RPC method.
type QueryParamsRequest struct {
}

func (m *QueryParamsRequest) Reset()         { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryParamsRequest) ProtoMessage()    {}
func (*QueryParamsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_71a7e047b63a0f83, []int{0}
}
func (m *QueryParamsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsRequest.Merge(m, src)
}
func (m *QueryParamsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsRequest proto.InternalMessageInfo

// QueryParamsResponse is response type for the Query/Params RPC method.
type QueryParamsResponse struct {
	// params holds all the parameters of this module.
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

func (m *QueryParamsResponse) Reset()         { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryParamsResponse) ProtoMessage()    {}
func (*QueryParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_71a7e047b63a0f83, []int{1}
}
func (m *QueryParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsResponse.Merge(m, src)
}
func (m *QueryParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsResponse proto.InternalMessageInfo

func (m *QueryParamsResponse) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

type QueryGetSessionRequest struct {
}

func (m *QueryGetSessionRequest) Reset()         { *m = QueryGetSessionRequest{} }
func (m *QueryGetSessionRequest) String() string { return proto.CompactTextString(m) }
func (*QueryGetSessionRequest) ProtoMessage()    {}
func (*QueryGetSessionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_71a7e047b63a0f83, []int{2}
}
func (m *QueryGetSessionRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetSessionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetSessionRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetSessionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetSessionRequest.Merge(m, src)
}
func (m *QueryGetSessionRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetSessionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetSessionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetSessionRequest proto.InternalMessageInfo

type QueryGetSessionResponse struct {
	PlaceholderText string `protobuf:"bytes,1,opt,name=placeholder_text,json=placeholderText,proto3" json:"placeholder_text,omitempty"`
}

func (m *QueryGetSessionResponse) Reset()         { *m = QueryGetSessionResponse{} }
func (m *QueryGetSessionResponse) String() string { return proto.CompactTextString(m) }
func (*QueryGetSessionResponse) ProtoMessage()    {}
func (*QueryGetSessionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_71a7e047b63a0f83, []int{3}
}
func (m *QueryGetSessionResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetSessionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetSessionResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetSessionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetSessionResponse.Merge(m, src)
}
func (m *QueryGetSessionResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetSessionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetSessionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetSessionResponse proto.InternalMessageInfo

func (m *QueryGetSessionResponse) GetPlaceholderText() string {
	if m != nil {
		return m.PlaceholderText
	}
	return ""
}

type QueryServicersRequest struct {
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryServicersRequest) Reset()         { *m = QueryServicersRequest{} }
func (m *QueryServicersRequest) String() string { return proto.CompactTextString(m) }
func (*QueryServicersRequest) ProtoMessage()    {}
func (*QueryServicersRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_71a7e047b63a0f83, []int{4}
}
func (m *QueryServicersRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryServicersRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryServicersRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryServicersRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryServicersRequest.Merge(m, src)
}
func (m *QueryServicersRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryServicersRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryServicersRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryServicersRequest proto.InternalMessageInfo

func (m *QueryServicersRequest) GetPagination() *query.PageRequest {
	if m != nil {
		return m.Pagination
	}
	return nil
}

type QueryServicersResponse struct {
	Servicers  []*Servicer         `protobuf:"bytes,1,rep,name=servicers,proto3" json:"servicers,omitempty"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryServicersResponse) Reset()         { *m = QueryServicersResponse{} }
func (m *QueryServicersResponse) String() string { return proto.CompactTextString(m) }
func (*QueryServicersResponse) ProtoMessage()    {}
func (*QueryServicersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_71a7e047b63a0f83, []int{5}
}
func (m *QueryServicersResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryServicersResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryServicersResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryServicersResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryServicersResponse.Merge(m, src)
}
func (m *QueryServicersResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryServicersResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryServicersResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryServicersResponse proto.InternalMessageInfo

func (m *QueryServicersResponse) GetServicers() []*Servicer {
	if m != nil {
		return m.Servicers
	}
	return nil
}

func (m *QueryServicersResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryParamsRequest)(nil), "poktroll.poktroll.QueryParamsRequest")
	proto.RegisterType((*QueryParamsResponse)(nil), "poktroll.poktroll.QueryParamsResponse")
	proto.RegisterType((*QueryGetSessionRequest)(nil), "poktroll.poktroll.QueryGetSessionRequest")
	proto.RegisterType((*QueryGetSessionResponse)(nil), "poktroll.poktroll.QueryGetSessionResponse")
	proto.RegisterType((*QueryServicersRequest)(nil), "poktroll.poktroll.QueryServicersRequest")
	proto.RegisterType((*QueryServicersResponse)(nil), "poktroll.poktroll.QueryServicersResponse")
}

func init() { proto.RegisterFile("poktroll/poktroll/query.proto", fileDescriptor_71a7e047b63a0f83) }

var fileDescriptor_71a7e047b63a0f83 = []byte{
	// 499 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x93, 0x31, 0x6f, 0x13, 0x31,
	0x14, 0xc7, 0xe3, 0x16, 0x22, 0xc5, 0x1d, 0x00, 0x53, 0x20, 0xb9, 0x16, 0x37, 0x9c, 0x20, 0x24,
	0x1d, 0xce, 0x6a, 0x3a, 0x20, 0xd6, 0x0a, 0xd1, 0x0d, 0x95, 0x2b, 0x13, 0x4b, 0xe4, 0x84, 0xa7,
	0x23, 0xe2, 0x7a, 0xbe, 0xda, 0x6e, 0x95, 0x32, 0x32, 0xb0, 0xb0, 0x20, 0xb1, 0x21, 0x3e, 0x50,
	0xc7, 0x4a, 0x2c, 0x4c, 0x08, 0x25, 0x7c, 0x10, 0x14, 0xdb, 0x97, 0xbb, 0x72, 0x17, 0xb5, 0x9b,
	0xe5, 0xf7, 0x7f, 0xef, 0xff, 0xf3, 0x7b, 0xcf, 0xf8, 0x61, 0x2a, 0x3e, 0x68, 0x29, 0xe2, 0x98,
	0x2d, 0x0e, 0xc7, 0x27, 0x20, 0xcf, 0x82, 0x54, 0x0a, 0x2d, 0xc8, 0x9d, 0xec, 0x36, 0xc8, 0x0e,
	0xde, 0x7a, 0x24, 0x22, 0x61, 0xa2, 0x6c, 0x7e, 0xb2, 0x42, 0x6f, 0x33, 0x12, 0x22, 0x8a, 0x81,
	0xf1, 0x74, 0xcc, 0x78, 0x92, 0x08, 0xcd, 0xf5, 0x58, 0x24, 0xca, 0x45, 0xb7, 0x47, 0x42, 0x1d,
	0x09, 0xc5, 0x86, 0x5c, 0x81, 0xad, 0xcf, 0x4e, 0x77, 0x86, 0xa0, 0xf9, 0x0e, 0x4b, 0x79, 0x34,
	0x4e, 0x8c, 0xd8, 0x69, 0x69, 0x99, 0x28, 0xe5, 0x92, 0x1f, 0x65, 0xb5, 0xb6, 0xca, 0x71, 0x05,
	0x4a, 0xe5, 0x05, 0x2a, 0x9e, 0xc4, 0x47, 0x5a, 0x48, 0x1b, 0xf6, 0xd7, 0x31, 0x79, 0x3d, 0x27,
	0x38, 0x30, 0x45, 0x43, 0x38, 0x3e, 0x01, 0xa5, 0xfd, 0x57, 0xf8, 0xee, 0xa5, 0x5b, 0x95, 0x8a,
	0x44, 0x01, 0x79, 0x86, 0xeb, 0xd6, 0xbc, 0x89, 0xda, 0xa8, 0xbb, 0xd6, 0x6f, 0x05, 0xa5, 0x86,
	0x04, 0x36, 0x65, 0xef, 0xc6, 0xf9, 0xef, 0xad, 0x5a, 0xe8, 0xe4, 0x7e, 0x13, 0xdf, 0x37, 0xf5,
	0xf6, 0x41, 0x1f, 0x5a, 0xba, 0xcc, 0xe9, 0x05, 0x7e, 0x50, 0x8a, 0x38, 0xb7, 0x1e, 0xbe, 0x9d,
	0xc6, 0x7c, 0x04, 0xef, 0x45, 0xfc, 0x0e, 0xe4, 0x40, 0xc3, 0x44, 0x1b, 0xdf, 0x46, 0x78, 0xab,
	0x70, 0xff, 0x06, 0x26, 0xda, 0x1f, 0xe0, 0x7b, 0xa6, 0xca, 0x21, 0xc8, 0xd3, 0xf1, 0x08, 0x64,
	0xf6, 0x10, 0xf2, 0x12, 0xe3, 0xbc, 0xa5, 0x8e, 0xba, 0x13, 0xd8, 0xfe, 0x07, 0xf3, 0xfe, 0x07,
	0x76, 0xbe, 0xae, 0xff, 0xc1, 0x01, 0x8f, 0xc0, 0xe5, 0x86, 0x85, 0x4c, 0xff, 0x07, 0x72, 0x2f,
	0x28, 0x38, 0x38, 0xcc, 0xe7, 0xb8, 0xa1, 0xb2, 0xcb, 0x26, 0x6a, 0xaf, 0x76, 0xd7, 0xfa, 0x1b,
	0x15, 0x7d, 0xc9, 0x12, 0xc3, 0x5c, 0x4d, 0xf6, 0x2f, 0xd1, 0xad, 0x18, 0xba, 0xa7, 0x57, 0xd2,
	0x59, 0xdf, 0x22, 0x5e, 0xff, 0xfb, 0x2a, 0xbe, 0x69, 0xf0, 0xc8, 0x47, 0x5c, 0xb7, 0x13, 0x20,
	0x4f, 0x2a, 0x20, 0xca, 0xa3, 0xf6, 0x3a, 0x57, 0xc9, 0xac, 0x9d, 0xff, 0xe8, 0xd3, 0xcf, 0xbf,
	0xdf, 0x56, 0x36, 0x48, 0x8b, 0x2d, 0xdb, 0x48, 0xf2, 0x05, 0x61, 0x9c, 0xcf, 0x91, 0xf4, 0x96,
	0x55, 0x2e, 0x6d, 0x81, 0xb7, 0x7d, 0x1d, 0xa9, 0x03, 0xe9, 0x18, 0x90, 0x36, 0xa1, 0x15, 0x20,
	0x11, 0xe8, 0x81, 0x5b, 0x7f, 0xf2, 0x19, 0xe1, 0xc6, 0x62, 0x5a, 0xa4, 0xbb, 0xcc, 0xe1, 0xff,
	0x95, 0xf1, 0x7a, 0xd7, 0x50, 0x3a, 0x94, 0xc7, 0x06, 0x85, 0x92, 0x4d, 0x56, 0xf5, 0x0b, 0x9d,
	0x7a, 0x6f, 0xf7, 0x7c, 0x4a, 0xd1, 0xc5, 0x94, 0xa2, 0x3f, 0x53, 0x8a, 0xbe, 0xce, 0x68, 0xed,
	0x62, 0x46, 0x6b, 0xbf, 0x66, 0xb4, 0xf6, 0xb6, 0xb5, 0x50, 0x4f, 0xf2, 0x44, 0x7d, 0x96, 0x82,
	0x1a, 0xd6, 0xcd, 0xf7, 0xdc, 0xfd, 0x17, 0x00, 0x00, 0xff, 0xff, 0x73, 0xd2, 0x81, 0x48, 0x92,
	0x04, 0x00, 0x00,
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
	// Parameters queries the parameters of the module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	// Queries a list of GetSession items.
	GetSession(ctx context.Context, in *QueryGetSessionRequest, opts ...grpc.CallOption) (*QueryGetSessionResponse, error)
	// Queries for all staked Servicers
	Servicers(ctx context.Context, in *QueryServicersRequest, opts ...grpc.CallOption) (*QueryServicersResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, "/poktroll.poktroll.Query/Params", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetSession(ctx context.Context, in *QueryGetSessionRequest, opts ...grpc.CallOption) (*QueryGetSessionResponse, error) {
	out := new(QueryGetSessionResponse)
	err := c.cc.Invoke(ctx, "/poktroll.poktroll.Query/GetSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Servicers(ctx context.Context, in *QueryServicersRequest, opts ...grpc.CallOption) (*QueryServicersResponse, error) {
	out := new(QueryServicersResponse)
	err := c.cc.Invoke(ctx, "/poktroll.poktroll.Query/Servicers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Parameters queries the parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// Queries a list of GetSession items.
	GetSession(context.Context, *QueryGetSessionRequest) (*QueryGetSessionResponse, error)
	// Queries for all staked Servicers
	Servicers(context.Context, *QueryServicersRequest) (*QueryServicersResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Params(ctx context.Context, req *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (*UnimplementedQueryServer) GetSession(ctx context.Context, req *QueryGetSessionRequest) (*QueryGetSessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSession not implemented")
}
func (*UnimplementedQueryServer) Servicers(ctx context.Context, req *QueryServicersRequest) (*QueryServicersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Servicers not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poktroll.poktroll.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poktroll.poktroll.Query/GetSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetSession(ctx, req.(*QueryGetSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Servicers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryServicersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Servicers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poktroll.poktroll.Query/Servicers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Servicers(ctx, req.(*QueryServicersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "poktroll.poktroll.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "GetSession",
			Handler:    _Query_GetSession_Handler,
		},
		{
			MethodName: "Servicers",
			Handler:    _Query_Servicers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "poktroll/poktroll/query.proto",
}

func (m *QueryParamsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QueryGetSessionRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetSessionRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetSessionRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryGetSessionResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetSessionResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetSessionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.PlaceholderText) > 0 {
		i -= len(m.PlaceholderText)
		copy(dAtA[i:], m.PlaceholderText)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.PlaceholderText)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryServicersRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryServicersRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryServicersRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryServicersResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryServicersResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryServicersResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Servicers) > 0 {
		for iNdEx := len(m.Servicers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Servicers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryParamsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryGetSessionRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryGetSessionResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.PlaceholderText)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryServicersRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryServicersResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Servicers) > 0 {
		for _, e := range m.Servicers {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryParamsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryParamsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryGetSessionRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryGetSessionRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetSessionRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryGetSessionResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryGetSessionResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetSessionResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PlaceholderText", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PlaceholderText = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryServicersRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryServicersRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryServicersRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageRequest{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryServicersResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryServicersResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryServicersResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Servicers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Servicers = append(m.Servicers, &Servicer{})
			if err := m.Servicers[len(m.Servicers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageResponse{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
