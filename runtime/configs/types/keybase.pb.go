// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.4
// source: keybase.proto

package types

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

type KeybaseType int32

const (
	KeybaseType_FILE KeybaseType = 0
	// Hashicorp Vault
	KeybaseType_VAULT KeybaseType = 1
)

// Enum value maps for KeybaseType.
var (
	KeybaseType_name = map[int32]string{
		0: "FILE",
		1: "VAULT",
	}
	KeybaseType_value = map[string]int32{
		"FILE":  0,
		"VAULT": 1,
	}
)

func (x KeybaseType) Enum() *KeybaseType {
	p := new(KeybaseType)
	*p = x
	return p
}

func (x KeybaseType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (KeybaseType) Descriptor() protoreflect.EnumDescriptor {
	return file_keybase_proto_enumTypes[0].Descriptor()
}

func (KeybaseType) Type() protoreflect.EnumType {
	return &file_keybase_proto_enumTypes[0]
}

func (x KeybaseType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use KeybaseType.Descriptor instead.
func (KeybaseType) EnumDescriptor() ([]byte, []int) {
	return file_keybase_proto_rawDescGZIP(), []int{0}
}

var File_keybase_proto protoreflect.FileDescriptor

var file_keybase_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6b, 0x65, 0x79, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x6b, 0x65, 0x79, 0x62, 0x61, 0x73, 0x65, 0x2a, 0x22, 0x0a, 0x0b, 0x4b, 0x65, 0x79, 0x62,
	0x61, 0x73, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x46, 0x49, 0x4c, 0x45, 0x10,
	0x00, 0x12, 0x09, 0x0a, 0x05, 0x56, 0x41, 0x55, 0x4c, 0x54, 0x10, 0x01, 0x42, 0x36, 0x5a, 0x34,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x6f, 0x6b, 0x74, 0x2d,
	0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x70, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x2f, 0x72,
	0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x2f, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_keybase_proto_rawDescOnce sync.Once
	file_keybase_proto_rawDescData = file_keybase_proto_rawDesc
)

func file_keybase_proto_rawDescGZIP() []byte {
	file_keybase_proto_rawDescOnce.Do(func() {
		file_keybase_proto_rawDescData = protoimpl.X.CompressGZIP(file_keybase_proto_rawDescData)
	})
	return file_keybase_proto_rawDescData
}

var file_keybase_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_keybase_proto_goTypes = []interface{}{
	(KeybaseType)(0), // 0: keybase.KeybaseType
}
var file_keybase_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_keybase_proto_init() }
func file_keybase_proto_init() {
	if File_keybase_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_keybase_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_keybase_proto_goTypes,
		DependencyIndexes: file_keybase_proto_depIdxs,
		EnumInfos:         file_keybase_proto_enumTypes,
	}.Build()
	File_keybase_proto = out.File
	file_keybase_proto_rawDesc = nil
	file_keybase_proto_goTypes = nil
	file_keybase_proto_depIdxs = nil
}
