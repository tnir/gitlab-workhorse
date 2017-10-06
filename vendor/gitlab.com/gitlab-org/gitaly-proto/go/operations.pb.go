// Code generated by protoc-gen-go. DO NOT EDIT.
// source: operations.proto

package gitaly

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type UserCreateBranchRequest struct {
	Repository *Repository `protobuf:"bytes,1,opt,name=repository" json:"repository,omitempty"`
	BranchName []byte      `protobuf:"bytes,2,opt,name=branch_name,json=branchName,proto3" json:"branch_name,omitempty"`
	User       *User       `protobuf:"bytes,3,opt,name=user" json:"user,omitempty"`
	StartPoint []byte      `protobuf:"bytes,4,opt,name=start_point,json=startPoint,proto3" json:"start_point,omitempty"`
}

func (m *UserCreateBranchRequest) Reset()                    { *m = UserCreateBranchRequest{} }
func (m *UserCreateBranchRequest) String() string            { return proto.CompactTextString(m) }
func (*UserCreateBranchRequest) ProtoMessage()               {}
func (*UserCreateBranchRequest) Descriptor() ([]byte, []int) { return fileDescriptor6, []int{0} }

func (m *UserCreateBranchRequest) GetRepository() *Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

func (m *UserCreateBranchRequest) GetBranchName() []byte {
	if m != nil {
		return m.BranchName
	}
	return nil
}

func (m *UserCreateBranchRequest) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *UserCreateBranchRequest) GetStartPoint() []byte {
	if m != nil {
		return m.StartPoint
	}
	return nil
}

type UserCreateBranchResponse struct {
	Branch *Branch `protobuf:"bytes,1,opt,name=branch" json:"branch,omitempty"`
	// Error returned by the pre-receive hook. If no error was thrown,
	// it's the empty string ("")
	PreReceiveError string `protobuf:"bytes,2,opt,name=pre_receive_error,json=preReceiveError" json:"pre_receive_error,omitempty"`
}

func (m *UserCreateBranchResponse) Reset()                    { *m = UserCreateBranchResponse{} }
func (m *UserCreateBranchResponse) String() string            { return proto.CompactTextString(m) }
func (*UserCreateBranchResponse) ProtoMessage()               {}
func (*UserCreateBranchResponse) Descriptor() ([]byte, []int) { return fileDescriptor6, []int{1} }

func (m *UserCreateBranchResponse) GetBranch() *Branch {
	if m != nil {
		return m.Branch
	}
	return nil
}

func (m *UserCreateBranchResponse) GetPreReceiveError() string {
	if m != nil {
		return m.PreReceiveError
	}
	return ""
}

type UserDeleteBranchRequest struct {
	Repository *Repository `protobuf:"bytes,1,opt,name=repository" json:"repository,omitempty"`
	BranchName []byte      `protobuf:"bytes,2,opt,name=branch_name,json=branchName,proto3" json:"branch_name,omitempty"`
	User       *User       `protobuf:"bytes,3,opt,name=user" json:"user,omitempty"`
}

func (m *UserDeleteBranchRequest) Reset()                    { *m = UserDeleteBranchRequest{} }
func (m *UserDeleteBranchRequest) String() string            { return proto.CompactTextString(m) }
func (*UserDeleteBranchRequest) ProtoMessage()               {}
func (*UserDeleteBranchRequest) Descriptor() ([]byte, []int) { return fileDescriptor6, []int{2} }

func (m *UserDeleteBranchRequest) GetRepository() *Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

func (m *UserDeleteBranchRequest) GetBranchName() []byte {
	if m != nil {
		return m.BranchName
	}
	return nil
}

func (m *UserDeleteBranchRequest) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

type UserDeleteBranchResponse struct {
	PreReceiveError string `protobuf:"bytes,1,opt,name=pre_receive_error,json=preReceiveError" json:"pre_receive_error,omitempty"`
}

func (m *UserDeleteBranchResponse) Reset()                    { *m = UserDeleteBranchResponse{} }
func (m *UserDeleteBranchResponse) String() string            { return proto.CompactTextString(m) }
func (*UserDeleteBranchResponse) ProtoMessage()               {}
func (*UserDeleteBranchResponse) Descriptor() ([]byte, []int) { return fileDescriptor6, []int{3} }

func (m *UserDeleteBranchResponse) GetPreReceiveError() string {
	if m != nil {
		return m.PreReceiveError
	}
	return ""
}

type UserDeleteTagRequest struct {
	Repository *Repository `protobuf:"bytes,1,opt,name=repository" json:"repository,omitempty"`
	TagName    []byte      `protobuf:"bytes,2,opt,name=tag_name,json=tagName,proto3" json:"tag_name,omitempty"`
	User       *User       `protobuf:"bytes,3,opt,name=user" json:"user,omitempty"`
}

func (m *UserDeleteTagRequest) Reset()                    { *m = UserDeleteTagRequest{} }
func (m *UserDeleteTagRequest) String() string            { return proto.CompactTextString(m) }
func (*UserDeleteTagRequest) ProtoMessage()               {}
func (*UserDeleteTagRequest) Descriptor() ([]byte, []int) { return fileDescriptor6, []int{4} }

func (m *UserDeleteTagRequest) GetRepository() *Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

func (m *UserDeleteTagRequest) GetTagName() []byte {
	if m != nil {
		return m.TagName
	}
	return nil
}

func (m *UserDeleteTagRequest) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

type UserDeleteTagResponse struct {
	PreReceiveError string `protobuf:"bytes,1,opt,name=pre_receive_error,json=preReceiveError" json:"pre_receive_error,omitempty"`
}

func (m *UserDeleteTagResponse) Reset()                    { *m = UserDeleteTagResponse{} }
func (m *UserDeleteTagResponse) String() string            { return proto.CompactTextString(m) }
func (*UserDeleteTagResponse) ProtoMessage()               {}
func (*UserDeleteTagResponse) Descriptor() ([]byte, []int) { return fileDescriptor6, []int{5} }

func (m *UserDeleteTagResponse) GetPreReceiveError() string {
	if m != nil {
		return m.PreReceiveError
	}
	return ""
}

type UserCreateTagRequest struct {
	Repository     *Repository `protobuf:"bytes,1,opt,name=repository" json:"repository,omitempty"`
	TagName        []byte      `protobuf:"bytes,2,opt,name=tag_name,json=tagName,proto3" json:"tag_name,omitempty"`
	User           *User       `protobuf:"bytes,3,opt,name=user" json:"user,omitempty"`
	TargetRevision []byte      `protobuf:"bytes,4,opt,name=target_revision,json=targetRevision,proto3" json:"target_revision,omitempty"`
	Message        []byte      `protobuf:"bytes,5,opt,name=message,proto3" json:"message,omitempty"`
}

func (m *UserCreateTagRequest) Reset()                    { *m = UserCreateTagRequest{} }
func (m *UserCreateTagRequest) String() string            { return proto.CompactTextString(m) }
func (*UserCreateTagRequest) ProtoMessage()               {}
func (*UserCreateTagRequest) Descriptor() ([]byte, []int) { return fileDescriptor6, []int{6} }

func (m *UserCreateTagRequest) GetRepository() *Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

func (m *UserCreateTagRequest) GetTagName() []byte {
	if m != nil {
		return m.TagName
	}
	return nil
}

func (m *UserCreateTagRequest) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *UserCreateTagRequest) GetTargetRevision() []byte {
	if m != nil {
		return m.TargetRevision
	}
	return nil
}

func (m *UserCreateTagRequest) GetMessage() []byte {
	if m != nil {
		return m.Message
	}
	return nil
}

type UserCreateTagResponse struct {
	Tag             *Tag   `protobuf:"bytes,1,opt,name=tag" json:"tag,omitempty"`
	Exists          bool   `protobuf:"varint,2,opt,name=exists" json:"exists,omitempty"`
	PreReceiveError string `protobuf:"bytes,3,opt,name=pre_receive_error,json=preReceiveError" json:"pre_receive_error,omitempty"`
}

func (m *UserCreateTagResponse) Reset()                    { *m = UserCreateTagResponse{} }
func (m *UserCreateTagResponse) String() string            { return proto.CompactTextString(m) }
func (*UserCreateTagResponse) ProtoMessage()               {}
func (*UserCreateTagResponse) Descriptor() ([]byte, []int) { return fileDescriptor6, []int{7} }

func (m *UserCreateTagResponse) GetTag() *Tag {
	if m != nil {
		return m.Tag
	}
	return nil
}

func (m *UserCreateTagResponse) GetExists() bool {
	if m != nil {
		return m.Exists
	}
	return false
}

func (m *UserCreateTagResponse) GetPreReceiveError() string {
	if m != nil {
		return m.PreReceiveError
	}
	return ""
}

func init() {
	proto.RegisterType((*UserCreateBranchRequest)(nil), "gitaly.UserCreateBranchRequest")
	proto.RegisterType((*UserCreateBranchResponse)(nil), "gitaly.UserCreateBranchResponse")
	proto.RegisterType((*UserDeleteBranchRequest)(nil), "gitaly.UserDeleteBranchRequest")
	proto.RegisterType((*UserDeleteBranchResponse)(nil), "gitaly.UserDeleteBranchResponse")
	proto.RegisterType((*UserDeleteTagRequest)(nil), "gitaly.UserDeleteTagRequest")
	proto.RegisterType((*UserDeleteTagResponse)(nil), "gitaly.UserDeleteTagResponse")
	proto.RegisterType((*UserCreateTagRequest)(nil), "gitaly.UserCreateTagRequest")
	proto.RegisterType((*UserCreateTagResponse)(nil), "gitaly.UserCreateTagResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for OperationService service

type OperationServiceClient interface {
	UserCreateBranch(ctx context.Context, in *UserCreateBranchRequest, opts ...grpc.CallOption) (*UserCreateBranchResponse, error)
	UserDeleteBranch(ctx context.Context, in *UserDeleteBranchRequest, opts ...grpc.CallOption) (*UserDeleteBranchResponse, error)
	UserCreateTag(ctx context.Context, in *UserCreateTagRequest, opts ...grpc.CallOption) (*UserCreateTagResponse, error)
	UserDeleteTag(ctx context.Context, in *UserDeleteTagRequest, opts ...grpc.CallOption) (*UserDeleteTagResponse, error)
}

type operationServiceClient struct {
	cc *grpc.ClientConn
}

func NewOperationServiceClient(cc *grpc.ClientConn) OperationServiceClient {
	return &operationServiceClient{cc}
}

func (c *operationServiceClient) UserCreateBranch(ctx context.Context, in *UserCreateBranchRequest, opts ...grpc.CallOption) (*UserCreateBranchResponse, error) {
	out := new(UserCreateBranchResponse)
	err := grpc.Invoke(ctx, "/gitaly.OperationService/UserCreateBranch", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *operationServiceClient) UserDeleteBranch(ctx context.Context, in *UserDeleteBranchRequest, opts ...grpc.CallOption) (*UserDeleteBranchResponse, error) {
	out := new(UserDeleteBranchResponse)
	err := grpc.Invoke(ctx, "/gitaly.OperationService/UserDeleteBranch", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *operationServiceClient) UserCreateTag(ctx context.Context, in *UserCreateTagRequest, opts ...grpc.CallOption) (*UserCreateTagResponse, error) {
	out := new(UserCreateTagResponse)
	err := grpc.Invoke(ctx, "/gitaly.OperationService/UserCreateTag", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *operationServiceClient) UserDeleteTag(ctx context.Context, in *UserDeleteTagRequest, opts ...grpc.CallOption) (*UserDeleteTagResponse, error) {
	out := new(UserDeleteTagResponse)
	err := grpc.Invoke(ctx, "/gitaly.OperationService/UserDeleteTag", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for OperationService service

type OperationServiceServer interface {
	UserCreateBranch(context.Context, *UserCreateBranchRequest) (*UserCreateBranchResponse, error)
	UserDeleteBranch(context.Context, *UserDeleteBranchRequest) (*UserDeleteBranchResponse, error)
	UserCreateTag(context.Context, *UserCreateTagRequest) (*UserCreateTagResponse, error)
	UserDeleteTag(context.Context, *UserDeleteTagRequest) (*UserDeleteTagResponse, error)
}

func RegisterOperationServiceServer(s *grpc.Server, srv OperationServiceServer) {
	s.RegisterService(&_OperationService_serviceDesc, srv)
}

func _OperationService_UserCreateBranch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserCreateBranchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperationServiceServer).UserCreateBranch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gitaly.OperationService/UserCreateBranch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperationServiceServer).UserCreateBranch(ctx, req.(*UserCreateBranchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OperationService_UserDeleteBranch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserDeleteBranchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperationServiceServer).UserDeleteBranch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gitaly.OperationService/UserDeleteBranch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperationServiceServer).UserDeleteBranch(ctx, req.(*UserDeleteBranchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OperationService_UserCreateTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserCreateTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperationServiceServer).UserCreateTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gitaly.OperationService/UserCreateTag",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperationServiceServer).UserCreateTag(ctx, req.(*UserCreateTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OperationService_UserDeleteTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserDeleteTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperationServiceServer).UserDeleteTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gitaly.OperationService/UserDeleteTag",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperationServiceServer).UserDeleteTag(ctx, req.(*UserDeleteTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _OperationService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gitaly.OperationService",
	HandlerType: (*OperationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UserCreateBranch",
			Handler:    _OperationService_UserCreateBranch_Handler,
		},
		{
			MethodName: "UserDeleteBranch",
			Handler:    _OperationService_UserDeleteBranch_Handler,
		},
		{
			MethodName: "UserCreateTag",
			Handler:    _OperationService_UserCreateTag_Handler,
		},
		{
			MethodName: "UserDeleteTag",
			Handler:    _OperationService_UserDeleteTag_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "operations.proto",
}

func init() { proto.RegisterFile("operations.proto", fileDescriptor6) }

var fileDescriptor6 = []byte{
	// 472 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x54, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0x8e, 0x9b, 0xe2, 0x96, 0x49, 0x68, 0xc3, 0x8a, 0x1f, 0x13, 0x51, 0xd5, 0xda, 0x03, 0x54,
	0x1c, 0x72, 0x08, 0x6f, 0x40, 0x81, 0x63, 0x41, 0x0b, 0x88, 0xa3, 0xb5, 0x0d, 0x23, 0x77, 0xa5,
	0xc6, 0x6b, 0x66, 0x37, 0x11, 0xe1, 0x05, 0xb8, 0xf2, 0x2a, 0xbc, 0x06, 0xcf, 0xc1, 0x83, 0x20,
	0x7b, 0xd7, 0x89, 0xed, 0x38, 0x12, 0x82, 0x03, 0x3d, 0xee, 0x37, 0xa3, 0x6f, 0xbe, 0xef, 0xd3,
	0xcc, 0xc2, 0x48, 0xe7, 0x48, 0xd2, 0x2a, 0x9d, 0x99, 0x49, 0x4e, 0xda, 0x6a, 0x16, 0xa6, 0xca,
	0xca, 0xeb, 0xd5, 0x78, 0x68, 0xae, 0x24, 0xe1, 0x27, 0x87, 0xf2, 0x1f, 0x01, 0x3c, 0xfc, 0x60,
	0x90, 0xce, 0x09, 0xa5, 0xc5, 0x17, 0x24, 0xb3, 0xd9, 0x95, 0xc0, 0xcf, 0x0b, 0x34, 0x96, 0x4d,
	0x01, 0x08, 0x73, 0x6d, 0x94, 0xd5, 0xb4, 0x8a, 0x82, 0x38, 0x38, 0x1b, 0x4c, 0xd9, 0xc4, 0xd1,
	0x4c, 0xc4, 0xba, 0x22, 0x6a, 0x5d, 0xec, 0x14, 0x06, 0x97, 0x25, 0x49, 0x92, 0xc9, 0x39, 0x46,
	0x7b, 0x71, 0x70, 0x36, 0x14, 0xe0, 0xa0, 0x0b, 0x39, 0x47, 0x16, 0xc3, 0xfe, 0xc2, 0x20, 0x45,
	0xfd, 0x92, 0x6e, 0x58, 0xd1, 0x15, 0x1a, 0x44, 0x59, 0x29, 0x28, 0x8c, 0x95, 0x64, 0x93, 0x5c,
	0xab, 0xcc, 0x46, 0xfb, 0x8e, 0xa2, 0x84, 0xde, 0x16, 0x08, 0xcf, 0x20, 0xda, 0x96, 0x6c, 0x72,
	0x9d, 0x19, 0x64, 0x4f, 0x20, 0x74, 0xc3, 0xbc, 0xde, 0xa3, 0x6a, 0x80, 0xef, 0xf3, 0x55, 0xf6,
	0x0c, 0xee, 0xe6, 0x84, 0x09, 0xe1, 0x0c, 0xd5, 0x12, 0x13, 0x24, 0xd2, 0x54, 0xaa, 0xbd, 0x2d,
	0x8e, 0x73, 0x42, 0xe1, 0xf0, 0x57, 0x05, 0xcc, 0xbf, 0xfb, 0x8c, 0x5e, 0xe2, 0x35, 0xde, 0x8c,
	0x8c, 0xf8, 0x6b, 0x17, 0x41, 0x53, 0x91, 0x8f, 0xa0, 0xd3, 0x5a, 0xd0, 0x6d, 0xed, 0x5b, 0x00,
	0xf7, 0x36, 0x44, 0xef, 0x65, 0xfa, 0x2f, 0xbe, 0x1e, 0xc1, 0xa1, 0x95, 0x69, 0xdd, 0xd4, 0x81,
	0x95, 0xe9, 0x1f, 0x3a, 0x3a, 0x87, 0xfb, 0x2d, 0x21, 0x7f, 0x61, 0xe7, 0xa7, 0xb7, 0xe3, 0x56,
	0xe3, 0x3f, 0xda, 0x61, 0x4f, 0xe1, 0xd8, 0x4a, 0x4a, 0xd1, 0x26, 0x84, 0x4b, 0x65, 0x94, 0xce,
	0xfc, 0x22, 0x1f, 0x39, 0x58, 0x78, 0x94, 0x45, 0x70, 0x30, 0x47, 0x63, 0x64, 0x8a, 0xd1, 0x2d,
	0x37, 0xc4, 0x3f, 0xf9, 0x57, 0x97, 0x48, 0xcd, 0x8b, 0x4f, 0xe4, 0x04, 0xfa, 0x56, 0xa6, 0xde,
	0xc5, 0xa0, 0x1a, 0x5e, 0x74, 0x14, 0x38, 0x7b, 0x00, 0x21, 0x7e, 0x51, 0xc6, 0x9a, 0x52, 0xf5,
	0xa1, 0xf0, 0xaf, 0xee, 0x20, 0xfb, 0x9d, 0x41, 0x4e, 0x7f, 0xed, 0xc1, 0xe8, 0x4d, 0xf5, 0x83,
	0xbc, 0x43, 0x5a, 0xaa, 0x19, 0xb2, 0x8f, 0x30, 0x6a, 0xdf, 0x1d, 0x3b, 0xad, 0x7b, 0xef, 0xf8,
	0x44, 0xc6, 0xf1, 0xee, 0x06, 0x67, 0x87, 0xf7, 0x2a, 0xe2, 0xfa, 0x36, 0x37, 0x89, 0x3b, 0x2e,
	0xaf, 0x49, 0xdc, 0x75, 0x08, 0xbc, 0xc7, 0x2e, 0xe0, 0x4e, 0x23, 0x42, 0xf6, 0x78, 0x5b, 0xcd,
	0x66, 0x4b, 0xc6, 0x27, 0x3b, 0xaa, 0x6d, 0xbe, 0xf5, 0x92, 0x36, 0xf9, 0xda, 0x47, 0xd4, 0xe4,
	0xdb, 0xda, 0x6c, 0xde, 0xbb, 0x0c, 0xcb, 0x4f, 0xf8, 0xf9, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x4d, 0x0d, 0x02, 0xa3, 0xae, 0x05, 0x00, 0x00,
}