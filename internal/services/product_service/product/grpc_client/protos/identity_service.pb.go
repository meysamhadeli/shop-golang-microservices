// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.26.1
// source: internal/services/product_service/product/grpc_client/protos/identity_service.proto

package identity_service

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
	Name   string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *User) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type GetUserByIdReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
}

func (x *GetUserByIdReq) Reset() {
	*x = GetUserByIdReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserByIdReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserByIdReq) ProtoMessage() {}

func (x *GetUserByIdReq) ProtoReflect() protoreflect.Message {
	mi := &file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserByIdReq.ProtoReflect.Descriptor instead.
func (*GetUserByIdReq) Descriptor() ([]byte, []int) {
	return file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetUserByIdReq) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type GetUserByIdRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *User `protobuf:"bytes,1,opt,name=User,proto3" json:"User,omitempty"`
}

func (x *GetUserByIdRes) Reset() {
	*x = GetUserByIdRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserByIdRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserByIdRes) ProtoMessage() {}

func (x *GetUserByIdRes) ProtoReflect() protoreflect.Message {
	mi := &file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserByIdRes.ProtoReflect.Descriptor instead.
func (*GetUserByIdRes) Descriptor() ([]byte, []int) {
	return file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_rawDescGZIP(), []int{2}
}

func (x *GetUserByIdRes) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

var File_internal_services_product_service_product_grpc_client_protos_identity_service_proto protoreflect.FileDescriptor

var file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_rawDesc = []byte{
	0x0a, 0x53, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2f, 0x67, 0x72, 0x70, 0x63,
	0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x69,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0x32, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12,
	0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x28, 0x0a, 0x0e, 0x47,
	0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x12, 0x16, 0x0a,
	0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x3c, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72,
	0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x73, 0x12, 0x2a, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x55,
	0x73, 0x65, 0x72, 0x32, 0x64, 0x0a, 0x0f, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x51, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x42, 0x79, 0x49, 0x64, 0x12, 0x20, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72,
	0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x1a, 0x20, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x73, 0x42, 0x15, 0x5a, 0x13, 0x2e, 0x2f, 0x3b,
	0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_rawDescOnce sync.Once
	file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_rawDescData = file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_rawDesc
)

func file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_rawDescGZIP() []byte {
	file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_rawDescOnce.Do(func() {
		file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_rawDescData)
	})
	return file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_rawDescData
}

var file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_goTypes = []interface{}{
	(*User)(nil),           // 0: identity_service.User
	(*GetUserByIdReq)(nil), // 1: identity_service.GetUserByIdReq
	(*GetUserByIdRes)(nil), // 2: identity_service.GetUserByIdRes
}
var file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_depIdxs = []int32{
	0, // 0: identity_service.GetUserByIdRes.User:type_name -> identity_service.User
	1, // 1: identity_service.IdentityService.GetUserById:input_type -> identity_service.GetUserByIdReq
	2, // 2: identity_service.IdentityService.GetUserById:output_type -> identity_service.GetUserByIdRes
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() {
	file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_init()
}
func file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_init() {
	if File_internal_services_product_service_product_grpc_client_protos_identity_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserByIdReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserByIdRes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_goTypes,
		DependencyIndexes: file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_depIdxs,
		MessageInfos:      file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_msgTypes,
	}.Build()
	File_internal_services_product_service_product_grpc_client_protos_identity_service_proto = out.File
	file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_rawDesc = nil
	file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_goTypes = nil
	file_internal_services_product_service_product_grpc_client_protos_identity_service_proto_depIdxs = nil
}
