// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.12.4
// source: home.proto

package service

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type EmptyRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EmptyRequest) Reset() {
	*x = EmptyRequest{}
	mi := &file_home_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EmptyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmptyRequest) ProtoMessage() {}

func (x *EmptyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_home_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmptyRequest.ProtoReflect.Descriptor instead.
func (*EmptyRequest) Descriptor() ([]byte, []int) {
	return file_home_proto_rawDescGZIP(), []int{0}
}

type HomeDataResponse struct {
	state            protoimpl.MessageState `protogen:"open.v1"`
	Categories       []*Category            `protobuf:"bytes,1,rep,name=categories,proto3" json:"categories,omitempty"`
	FeaturedProducts []*Product             `protobuf:"bytes,2,rep,name=featured_products,json=featuredProducts,proto3" json:"featured_products,omitempty"`
	Brands           []*Brand               `protobuf:"bytes,3,rep,name=brands,proto3" json:"brands,omitempty"`
	Error            *ErrorResponse         `protobuf:"bytes,4,opt,name=error,proto3" json:"error,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *HomeDataResponse) Reset() {
	*x = HomeDataResponse{}
	mi := &file_home_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HomeDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HomeDataResponse) ProtoMessage() {}

func (x *HomeDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_home_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HomeDataResponse.ProtoReflect.Descriptor instead.
func (*HomeDataResponse) Descriptor() ([]byte, []int) {
	return file_home_proto_rawDescGZIP(), []int{1}
}

func (x *HomeDataResponse) GetCategories() []*Category {
	if x != nil {
		return x.Categories
	}
	return nil
}

func (x *HomeDataResponse) GetFeaturedProducts() []*Product {
	if x != nil {
		return x.FeaturedProducts
	}
	return nil
}

func (x *HomeDataResponse) GetBrands() []*Brand {
	if x != nil {
		return x.Brands
	}
	return nil
}

func (x *HomeDataResponse) GetError() *ErrorResponse {
	if x != nil {
		return x.Error
	}
	return nil
}

var File_home_proto protoreflect.FileDescriptor

var file_home_proto_rawDesc = string([]byte{
	0x0a, 0x0a, 0x68, 0x6f, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x0e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x0e, 0x0a, 0x0c, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x22, 0xd2, 0x01, 0x0a, 0x10, 0x48, 0x6f, 0x6d, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f, 0x0a, 0x0a, 0x63, 0x61, 0x74, 0x65,
	0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x0a, 0x63,
	0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x12, 0x3b, 0x0a, 0x11, 0x66, 0x65, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x64, 0x5f, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x52, 0x10, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x64, 0x50, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x12, 0x24, 0x0a, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42,
	0x72, 0x61, 0x6e, 0x64, 0x52, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x73, 0x12, 0x2a, 0x0a, 0x05,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x32, 0x4a, 0x0a, 0x0b, 0x48, 0x6f, 0x6d, 0x65,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3b, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x48, 0x6f,
	0x6d, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x48, 0x6f, 0x6d, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0b, 0x5a, 0x09, 0x2e, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_home_proto_rawDescOnce sync.Once
	file_home_proto_rawDescData []byte
)

func file_home_proto_rawDescGZIP() []byte {
	file_home_proto_rawDescOnce.Do(func() {
		file_home_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_home_proto_rawDesc), len(file_home_proto_rawDesc)))
	})
	return file_home_proto_rawDescData
}

var file_home_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_home_proto_goTypes = []any{
	(*EmptyRequest)(nil),     // 0: proto.EmptyRequest
	(*HomeDataResponse)(nil), // 1: proto.HomeDataResponse
	(*Category)(nil),         // 2: proto.Category
	(*Product)(nil),          // 3: proto.Product
	(*Brand)(nil),            // 4: proto.Brand
	(*ErrorResponse)(nil),    // 5: proto.ErrorResponse
}
var file_home_proto_depIdxs = []int32{
	2, // 0: proto.HomeDataResponse.categories:type_name -> proto.Category
	3, // 1: proto.HomeDataResponse.featured_products:type_name -> proto.Product
	4, // 2: proto.HomeDataResponse.brands:type_name -> proto.Brand
	5, // 3: proto.HomeDataResponse.error:type_name -> proto.ErrorResponse
	0, // 4: proto.HomeService.GetHomeData:input_type -> proto.EmptyRequest
	1, // 5: proto.HomeService.GetHomeData:output_type -> proto.HomeDataResponse
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_home_proto_init() }
func file_home_proto_init() {
	if File_home_proto != nil {
		return
	}
	file_entities_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_home_proto_rawDesc), len(file_home_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_home_proto_goTypes,
		DependencyIndexes: file_home_proto_depIdxs,
		MessageInfos:      file_home_proto_msgTypes,
	}.Build()
	File_home_proto = out.File
	file_home_proto_goTypes = nil
	file_home_proto_depIdxs = nil
}
