// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: errcode/v1/errors.proto

package errcode

import (
	"reflect"
	"sync"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 错误
type Error int32

const (
	// 未知类型
	// @code=UNKNOWN
	Error_ERR_UNKNOWN Error = 0
	// 参数错误
	// @code=INVALID_ARGUMENT
	Error_ERR_INVALID_ARGUMENT Error = 1
	// 找不到资源
	// @code=NOT_FOUND
	Error_ERR_NOT_FOUND Error = 2
	// db错误
	// @code=INTERNAL
	Error_ERR_DB_ERROR Error = 3
)

// Enum value maps for Error.
var (
	Error_name = map[int32]string{
		0: "ERR_UNKNOWN",
		1: "ERR_INVALID_ARGUMENT",
		2: "ERR_NOT_FOUND",
		3: "ERR_DB_ERROR",
	}
	Error_value = map[string]int32{
		"ERR_UNKNOWN":          0,
		"ERR_INVALID_ARGUMENT": 1,
		"ERR_NOT_FOUND":        2,
		"ERR_DB_ERROR":         3,
	}
)

func (x Error) Enum() *Error {
	p := new(Error)
	*p = x
	return p
}

func (x Error) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Error) Descriptor() protoreflect.EnumDescriptor {
	return file_errcode_v1_errors_proto_enumTypes[0].Descriptor()
}

func (Error) Type() protoreflect.EnumType {
	return &file_errcode_v1_errors_proto_enumTypes[0]
}

func (x Error) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Error.Descriptor instead.
func (Error) EnumDescriptor() ([]byte, []int) {
	return file_errcode_v1_errors_proto_rawDescGZIP(), []int{0}
}

var File_errcode_v1_errors_proto protoreflect.FileDescriptor

var file_errcode_v1_errors_proto_rawDesc = []byte{
	0x0a, 0x17, 0x65, 0x72, 0x72, 0x63, 0x6f, 0x64, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x65, 0x72, 0x72, 0x63, 0x6f,
	0x64, 0x65, 0x2e, 0x76, 0x31, 0x2a, 0x57, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x0f,
	0x0a, 0x0b, 0x45, 0x52, 0x52, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12,
	0x18, 0x0a, 0x14, 0x45, 0x52, 0x52, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x41,
	0x52, 0x47, 0x55, 0x4d, 0x45, 0x4e, 0x54, 0x10, 0x01, 0x12, 0x11, 0x0a, 0x0d, 0x45, 0x52, 0x52,
	0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x02, 0x12, 0x10, 0x0a, 0x0c,
	0x45, 0x52, 0x52, 0x5f, 0x44, 0x42, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x03, 0x42, 0x16,
	0x5a, 0x14, 0x65, 0x72, 0x72, 0x63, 0x6f, 0x64, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x65, 0x72, 0x72,
	0x63, 0x6f, 0x64, 0x65, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_errcode_v1_errors_proto_rawDescOnce sync.Once
	file_errcode_v1_errors_proto_rawDescData = file_errcode_v1_errors_proto_rawDesc
)

func file_errcode_v1_errors_proto_rawDescGZIP() []byte {
	file_errcode_v1_errors_proto_rawDescOnce.Do(func() {
		file_errcode_v1_errors_proto_rawDescData = protoimpl.X.CompressGZIP(file_errcode_v1_errors_proto_rawDescData)
	})
	return file_errcode_v1_errors_proto_rawDescData
}

var file_errcode_v1_errors_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_errcode_v1_errors_proto_goTypes = []interface{}{
	(Error)(0), // 0: errcode.v1.Error
}
var file_errcode_v1_errors_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_errcode_v1_errors_proto_init() }
func file_errcode_v1_errors_proto_init() {
	if File_errcode_v1_errors_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_errcode_v1_errors_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_errcode_v1_errors_proto_goTypes,
		DependencyIndexes: file_errcode_v1_errors_proto_depIdxs,
		EnumInfos:         file_errcode_v1_errors_proto_enumTypes,
	}.Build()
	File_errcode_v1_errors_proto = out.File
	file_errcode_v1_errors_proto_rawDesc = nil
	file_errcode_v1_errors_proto_goTypes = nil
	file_errcode_v1_errors_proto_depIdxs = nil
}
