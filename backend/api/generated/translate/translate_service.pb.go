// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v3.12.4
// source: translate_service.proto

package translate

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

type TranslateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Text          string                 `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TranslateRequest) Reset() {
	*x = TranslateRequest{}
	mi := &file_translate_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TranslateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TranslateRequest) ProtoMessage() {}

func (x *TranslateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_translate_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TranslateRequest.ProtoReflect.Descriptor instead.
func (*TranslateRequest) Descriptor() ([]byte, []int) {
	return file_translate_service_proto_rawDescGZIP(), []int{0}
}

func (x *TranslateRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

type TranslateResult struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Lines         string                 `protobuf:"bytes,1,opt,name=lines,proto3" json:"lines,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TranslateResult) Reset() {
	*x = TranslateResult{}
	mi := &file_translate_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TranslateResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TranslateResult) ProtoMessage() {}

func (x *TranslateResult) ProtoReflect() protoreflect.Message {
	mi := &file_translate_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TranslateResult.ProtoReflect.Descriptor instead.
func (*TranslateResult) Descriptor() ([]byte, []int) {
	return file_translate_service_proto_rawDescGZIP(), []int{1}
}

func (x *TranslateResult) GetLines() string {
	if x != nil {
		return x.Lines
	}
	return ""
}

var File_translate_service_proto protoreflect.FileDescriptor

var file_translate_service_proto_rawDesc = []byte{
	0x0a, 0x17, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x6c, 0x61, 0x74, 0x65, 0x22, 0x26, 0x0a, 0x10, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x22, 0x27, 0x0a, 0x0f,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6e, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x6c, 0x69, 0x6e, 0x65, 0x73, 0x32, 0x61, 0x0a, 0x10, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61,
	0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4d, 0x0a, 0x12, 0x50, 0x72, 0x6f,
	0x63, 0x65, 0x73, 0x73, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x1b, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x2e, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x6c, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x42, 0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x4f, 0x53, 0x6f, 0x6d, 0x6e, 0x75, 0x73, 0x2f,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6c, 0x61, 0x74, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61,
	0x74, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_translate_service_proto_rawDescOnce sync.Once
	file_translate_service_proto_rawDescData = file_translate_service_proto_rawDesc
)

func file_translate_service_proto_rawDescGZIP() []byte {
	file_translate_service_proto_rawDescOnce.Do(func() {
		file_translate_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_translate_service_proto_rawDescData)
	})
	return file_translate_service_proto_rawDescData
}

var file_translate_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_translate_service_proto_goTypes = []any{
	(*TranslateRequest)(nil), // 0: translate.TranslateRequest
	(*TranslateResult)(nil),  // 1: translate.TranslateResult
}
var file_translate_service_proto_depIdxs = []int32{
	0, // 0: translate.TranslateService.ProcessTranslation:input_type -> translate.TranslateRequest
	1, // 1: translate.TranslateService.ProcessTranslation:output_type -> translate.TranslateResult
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_translate_service_proto_init() }
func file_translate_service_proto_init() {
	if File_translate_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_translate_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_translate_service_proto_goTypes,
		DependencyIndexes: file_translate_service_proto_depIdxs,
		MessageInfos:      file_translate_service_proto_msgTypes,
	}.Build()
	File_translate_service_proto = out.File
	file_translate_service_proto_rawDesc = nil
	file_translate_service_proto_goTypes = nil
	file_translate_service_proto_depIdxs = nil
}
